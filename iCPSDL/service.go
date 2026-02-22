package iCPSDL

import "fmt"
import "iCPSDL/util"

type ServiceMessage struct {
    Jmap util.JsonMap
    Request string
    Command string
    State string
    Process string
    ConfArgs string
    Ack string
}

/*********************
  iCPSDL execute a command
 *********************/

func Execute(args []string) {
    iCPSDL := NewiCPSDL()
    iCPSDL.Execute(args)
}

 /*********************
   iCPSDL service
  *********************/

func Service(eventChannel chan *ServiceMessage, engineChannel chan *ServiceMessage) {
    fmt.Println("[iCPSDL Service] Starting.")
    iCPSDL := NewiCPSDL()

    for iCPSDL.quit == false {
        select {
            case input := <- eventChannel:
                _, output := iCPSDL.handleRequest(input.Command)
                eventChannel <- &ServiceMessage{Jmap: nil, Request: "", Command: "", Process: "", State: "", ConfArgs: "", Ack: output}
            case input := <- engineChannel:
                switch input.Request {
                    case "initalise":
                        fmt.Println("[iCPSDL Service] Initialising.")
                        _, output := iCPSDL.handleRequest(input.Command)
                        engineChannel <- &ServiceMessage{Jmap: nil, Request: "", Command: "", Process: "", State: "", ConfArgs: "", Ack: output}
                    case "configure":
                        command := fmt.Sprintf("seGraph := translate %v\n", input.Process)
                        command += fmt.Sprintf("trees := traverse %v seGraph", input.State)
                        expr, _ := iCPSDL.handleRequest(command)
                        f := expr.(*forest)
                        exprs := make([]expression, len(f.trees))
                        min_len := 0
                        min_index := 0
                        var min_expr *local = nil
                        output := ""
                        expr_output := ""
                        for i := range f.trees {
                            command = fmt.Sprintf("lc%v := configure trees[%v] %v\n", i, i, input.ConfArgs)
                            exprs[i], output = iCPSDL.handleRequest(command)
                            lc := exprs[i].(*local)
                            if min_expr == nil || len(lc.configuration) < min_len {
                                min_index = i
                                min_expr = lc
                                expr_output = output
                                min_len = len(lc.configuration)
                            }
                        }
                        command = fmt.Sprintf("mermaid trees[%v]", min_index)
                        iCPSDL.handleRequest(command)
                        fmt.Println("[iCPSDL Service] Select control loop configuration. Create a corresponding mermaid diagram.")
                        jmap := min_expr.json()
                        //fmt.Println(jmap, "\n")
                        engineChannel <- &ServiceMessage{Jmap: jmap, Request: "", Command: "", Process: "", State: "", ConfArgs: "", Ack: expr_output}
                }
        }
    }
}
