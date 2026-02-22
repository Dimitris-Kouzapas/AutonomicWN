package ast

import (
    "fmt"
    "strings"
    "sessions/util"
)

type program struct {
    file string
    modules map[string]*module
    mainModule *module
    suffix string
}

func NewProgram(file string) *program {
    return &program{
        file:    file,
        modules: make(map[string]*module),
        suffix:  ".sess",
    }
}

func (p *program) Parse() (*module, util.ErrorLog) {
    log := newCustomErrorListener()
    if !strings.HasSuffix(p.file, p.suffix) {
        msg := fmt.Sprintf("file %q should have suffix %q.", p.file, p.suffix)
        return nil, log.Add(util.NewSystemError(msg))
    }
    moduleName := strings.TrimSuffix(p.file, p.suffix)
    mod := p.parseModule(moduleName, log)
    if mod == nil {
        return nil, log
    }

    p.mainModule = newModule(nil, mod.line)
    for _, m := range p.modules {
        for _, decl := range m.decls {
            p.mainModule.addDeclaration(decl)
        }
    }
    p.mainModule.filename = mod.file()
    return p.mainModule, log
}

func (p *program) parseModule(moduleName string, log *customErrorListener) *module {
    if m, ok := p.modules[moduleName]; ok {
        return m
    }
    filename := moduleName + p.suffix
    m := parse(filename, log)
    if m == nil {
        return nil
    }
    m.setFilename(filename)
    p.modules[moduleName] = m
    for _, imp := range m.getImports() {
        imp = strings.ReplaceAll(imp, ".", "\\")
        p.parseModule(imp, log)
    }
    return m
}

