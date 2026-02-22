package iCPSDL

import "fmt"
import "os"
import "bufio"

import "iCPSDL/util"

type iCPSDL struct {
    expressions map[string]expression
    suffix string
    infoFile string
    quit bool
}

func NewiCPSDL() (this *iCPSDL) {
    this = new(iCPSDL)
    this.expressions = make(map[string]expression)
    this.suffix = ".dss"
    this.infoFile = "./res/help.res"
    this.quit = false
    return
}

func (this *iCPSDL) executeCommand(input string) {
    log := newCustomErrorListener("terminal")
    expressions := parseString(input, log)
    for _, expr := range expressions {
        expr.setFilename("terminal")
    }

    if log.HasErrors() == true {
        log.PrintErrors()
    } else {
        _, log := this.executeStack(expressions)
        if log.HasErrors() {
            log.PrintErrors()
        }
    }
}

func (this *iCPSDL) terminal() {
    reader := bufio.NewReader(os.Stdin)
    for this.quit == false {
        fmt.Print(">: ")
        input, _ := reader.ReadString('\n')

        this.executeCommand(input)
        /*log := newCustomErrorListener("terminal")
        expressions := parseString(input, log)

        for _, expr := range expressions {
            expr.setFilename("terminal")
        }

        if log.HasErrors() == true {
            log.PrintErrors()
        } else {
            _, log := this.executeStack(expressions)
            if log.HasErrors() {
                log.PrintErrors()
            }
        }*/
    }
}

func (this *iCPSDL) handleRequest(input string) (expression, string) {
    log := newCustomErrorListener("service")
    expressions := parseString(input, log)
    for _, expr := range expressions {
        expr.setFilename("service")
    }

    if log.HasErrors() == true {
        return nil, log.ToString()
    } else {
        expr, log := this.executeStack(expressions)
        if log.HasErrors() {
            return nil, log.ToString()
        } else {
            stream := util.NewStream()
            expr.prettyPrint(stream)

            return expr, stream.ToString()
        }
    }
    return nil, ""
}

func (this *iCPSDL) Execute(args []string) {
    term := false
    for i := 0; i < len(args); i++ {
        if args[i] == "-t" {
            term = true
        } else if args[i] == "-load" {
            i += 1
            this.executeCommand("load " + args[i])
        } else {
            fmt.Println("  invalid command line argument: ", args[i])
        }
    }

    if term == true {
        this.terminal()
    }
}

func (this *iCPSDL) executeStack(expressions []expression) (expression, util.ErrorLog) {
    var expr expression = nil
    for _, stmt := range expressions {
        var log util.ErrorLog
        expr, log = stmt.execute(this)
        if log.HasErrors() == true {
            return nil, log
        }
    }
    return expr, util.NewErrorLog()
}

func (this *iCPSDL) Request(input string) string {
    _, output := this.handleRequest(input)
    return output
}

func (this *iCPSDL) String() string {
    return "iCPSDL service"
}