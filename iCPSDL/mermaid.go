package iCPSDL

import "fmt"
import "iCPSDL/util"
import (
    "net/http"
)

/*****************************************************************************
 * mermaid server
 *****************************************************************************/

func createMermaidHTML(mermaid string) string {
    html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Mermaid Diagram</title>
            <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
            <script>mermaid.initialize({startOnLoad:true});</script>
        </head>
        <body>
            <div class="mermaid">`
    html += mermaid
    html += `
            </div>
        </body>
    </html>
    `
    return html
}

func createIndexHTML() string {
    return `
        <!DOCTYPE html>
            <html lang="en">
                <head>
                    <meta charset="UTF-8">
                    <meta name="viewport" content="width=device-width, initial-scale=1.0">
                    <title>Mermaid Diagram Request</title>
                </head>
                <body>
                    <h1>Please render a Mermaid diagram</h1>
                </body>
            </html>
    `
}

var mmd string = ""
func ListenAndServeLocalHost() {
    http.HandleFunc("/mermaid", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        html := createIndexHTML()
        if mmd != "" { html = createMermaidHTML(mmd) }
        w.Write([]byte(html))
    })

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
    //server := &http.Server{
    //    Addr:    ":8080", // Set the bind address
    //    Handler: nil,    // Use the default mux (router)
    //}

    // Start the server
    //if err := server.ListenAndServe(); err != nil {
    //    fmt.Println("Error starting server:", err)
    //}
}

/*****************************************************************************
 * mermaid expession
 *****************************************************************************/

type mermaidPrinter interface {
   mermaid(*util.Stream)
}

type mermaidExpression struct {
    expressionImpl
    mmd string
}

func newMermaid(expr expression, line int) (this *mermaidExpression) {
    this = new(mermaidExpression)
    this.newExpression(expr, line)
    this.mmd = ""
    return
}

func (this *mermaidExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expr.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    var mmd mermaidPrinter = nil
    switch exec := executed.(type) {
        case *process:
            mmd = exec
        case *arGraph:
            mmd = exec
        case *arTree:
            mmd = exec
        //case *forest:
            //for _, tree := range exec.trees {
            //    this.createDiagram(tree, log)
            //}
        default:
            error := fmt.Sprintf("%v %v is neither a process nor a state estimation graph.", util.CapitaliseFirst(executed.getType()), this.getModule())
            this.reportCommandLineError(error, log)
            return nil, log
    }
    stream := util.NewStream()
    mmd.mermaid(stream)
    this.mmd = stream.ToString()
    this.createDiagram(log)
    return this, log
}

func (this *mermaidExpression) createDiagram(log util.ErrorLog) {
    html := createMermaidHTML(this.mmd)
    file := this.getModule()
    //filename, err := util.CreateTmpFile(html, "html")
    //if err != nil {
    //    error := fmt.Sprintf("Failed to create mermaid file: %s", err)
    //    log.Add(util.NewSystemError(error))
    //    return
    //}
    filename := file + ".html"
    err := util.CreateFile(filename, html)
    if err != nil {
        error := fmt.Sprintf("Failed to create mermaid file: %s", err)
        log.Add(util.NewSystemError(error))
        return
    }

    err = util.OpenBrowser(filename)
    if err != nil {
        log.Add(util.NewSystemError(err.Error()))
    }
    return
}

func (this *mermaidExpression) render() {
  mmd = this.mmd
}

func (this *mermaidExpression) prettyPrint(stream *util.Stream) {
    stream.Print(this.mmd)
    //stream.Print("mermaid ")
    //this.expressionImpl.prettyPrint(stream)
}

func (this *mermaidExpression) getType() string {
    return "mermaid"
}

func (this *mermaidExpression) info() string {
    return "Type: mermaid, module type: " + this.expressionImpl.getType()
}
