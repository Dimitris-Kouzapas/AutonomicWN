package main
import (
    "fmt"
    "os"

    "sessions"
    "sessions/util"
    "sessions/api"
    "iCPSDL"
)

var programList []string = []string {
    "main.sess",
}

func parseArgs(args []string) *util.Config {
    defs := []util.FlagDef {
        {"verbose", []string{"v"}, false, "Enable verbose mode"},
        {"analysis", []string{"a"}, false, "Enable analysis mode"},
        {"no-exec", []string{"ne"}, false, "Disable execution mode"},
        {"prettyPrint", []string{"p"}, false, "Enable prettyPrint mode"},
        {"usage", []string{"u"}, false, "Print usage"},
    }

    cfg := util.ParseFlags(defs, args)
    if ok1, ok2 := cfg.Bool("usage"); ok1 && ok2 {
        fmt.Fprintln(os.Stdout, cfg.Usage())
    } 
    files := cfg.Args()
    if len(files) == 0 {
        cfg.SetArgs(programList)
    }
    return cfg
}


func runFile(file string, cfg *util.Config) {
    fmt.Println("#############################################################")
    fmt.Println("Compiling file: ", file)
    program := ast.NewProgram(file)
    m, log := program.Parse()
    if log.HasErrors() {
        log.PrintErrors()
        return
    }

    elog, rlog := m.Validate(cfg)
    if elog.HasErrors() {
        elog.PrintErrors()
        return
    }

    if rlog.HasReports() {
        rlog.PrintReports()
    }

    if ok1, ok2 := cfg.Bool("prettyPrint"); ok1 && ok2 {
        stream := util.NewStream()
        m.PrettyPrint(stream)
        fmt.Fprintln(os.Stdout, stream)
    }

    if ok1, ok2 := cfg.Bool("no-exec"); !ok1 && ok2 {
        fmt.Println("Starting execution")
        m.Execute()
        fmt.Println("Finished execution")
    }
    fmt.Println("#############################################################")
}

func newiCPSDL() api.API {
    return iCPSDL.NewiCPSDL()
}

func configure() {
    api.Register("iCPSDL", newiCPSDL)
}

/*********************
 * main
 *********************/

func main() {
    cfg := parseArgs(os.Args[1:])
    files := cfg.Args()

    configure()

    for _, file := range files {
        runFile(file, cfg)    
    }


    // iCPSDL := iCPSDL.NewiCPSDL()
    // iCPSDL.Execute([]string{"-t"})
}




