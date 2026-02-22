package ast

import (
    "sessions/util"
)

type module struct {
    baseNode
    imports []string
    decls []declaration
    nameDecl []declaration
    abstractionMap map[string]*abstractionDeclaration
    typeDeclMap map[string]*typeDeclaration
    //brokerRegistry map[string]*brokerDeclaration
    valueMap map[string]*valueDeclaration
    sessionAsgmMap map[string]*sessionAssignment
}

func newModule(imports []string, line int) *module {
    return &module {
        baseNode: baseNode {
            line: line,
        },
        imports: imports,
        decls: make([]declaration, 0),
        nameDecl: make([]declaration, 0),
        abstractionMap: make(map[string]*abstractionDeclaration),
        typeDeclMap: make(map[string]*typeDeclaration),
        valueMap: make(map[string]*valueDeclaration),
        sessionAsgmMap: make(map[string]*sessionAssignment),
    }
}

func (m *module) addDeclaration(decl declaration) {
    m.decls = append(m.decls, decl)
    switch d := decl.(type) {
        case *abstractionDeclaration:
            m.nameDecl = append(m.nameDecl, decl)
            m.abstractionMap[d.getName()] = d
        case *typeDeclaration:
            m.nameDecl = append(m.nameDecl, decl)
            m.typeDeclMap[d.getName()] = d
        case *valueDeclaration:
            m.nameDecl = append(m.nameDecl, decl)
            m.valueMap[d.getName()] = d
        case *sessionAssignment:
            m.sessionAsgmMap[d.getName()] = d
    }
}

func (m *module) setFilename(filename string) {
    m.baseNode.setFilename(filename)
    for _, decl := range m.decls {
        decl.setFilename(filename)
    }
}

func (m *module) getImports() []string { return m.imports }

func (m *module) getDecl(name string) declaration {
    if decl, ok := m.abstractionMap[name]; ok {
        return decl
    }
    return nil
}

func (m *module) getType(name string) typedef {
    if tdecl, ok := m.typeDeclMap[name]; ok {
        return tdecl.tdef
    }

    return nil
}

func (m *module) getSession(name string) typedef {
    if decl, ok := m.sessionAsgmMap[name]; ok {
        return decl.session
    }
    return nil
}

func (m *module) getAbstraction(name string) *abstraction {
    if abstDecl, ok := m.abstractionMap[name]; ok {
        return abstDecl.abstr
    }
    return nil
}

func (m *module) getValue(name string) expression {
    if decl, ok := m.valueMap[name]; ok {
        return decl.value
    }
    return nil
}

func (m *module) typeCheck(log util.ErrorLog) {
    for i := range m.nameDecl {
        for j := i + 1; j < len(m.nameDecl); j++ {
            if m.nameDecl[i].getName() == m.nameDecl[j].getName() {
                m.nameDecl[j].reportErrorf(
                    log,
                    "duplicate definition of name: %q; originally defined in %s:%d",
                    m.nameDecl[j].getName(), m.nameDecl[i].file(), m.nameDecl[i].lineno(),
                )
            }
        }
    }
    ctx := newTypeCheckContext(m)
    for _, decl := range m.decls {
        decl.typeCheck(ctx, log)
    }
}

func (m *module) expressionCheck(log util.ErrorLog) {
    ctx := newExpressionCheckContext(m)
    for _, decl := range m.decls {
        decl.expressionCheck(ctx, log)
    }
}

func (m *module) projectionCheck(config *util.Config, elog util.ErrorLog, rlog util.ReportLog) {
    ctx := newProjectionCheckContext(config)
    for _, decl := range m.decls {
        decl.projectionCheck(ctx, elog, rlog)
    }
}

func (m *module) sessionCheck(config *util.Config, elog util.ErrorLog, rlog util.ReportLog) {
    ctx := newSessionCheckContext(config)
    for _, decl := range m.decls {
        decl.sessionCheck(ctx, elog, rlog)
    }
}

func (m *module) Validate(config *util.Config) (elog util.ErrorLog, rlog util.ReportLog) {
    elog = util.NewErrorLog()
    rlog = util.NewReportLog()
    m.typeCheck(elog)
    if elog.HasErrors() {
        return
    }
    m.projectionCheck(config, elog, rlog)
    if elog.HasErrors() {
        return
    }
    m.expressionCheck(elog)
    if elog.HasErrors() {
        return
    }
    m.sessionCheck(config, elog, rlog)
    return
}

func (m *module) Execute() {
    ctx := newEvaluationContext(m)
    if main := m.getDecl("main"); main != nil {
        main.execute(ctx)
    }
    ctx.wait()
}

func (m *module) PrettyPrint(iw util.IndentedWriter) {
    for _, imp := range m.imports {
        iw.Println("import", imp)
    }
    iw.Println()
    for _, d := range m.decls {
        d.prettyPrint(iw)
        iw.Println()
        iw.Println()
    }
}

func (m *module) String() string {
    s := ""
    for _, imp := range m.imports {
        s += "import " + imp + "\n"
    }
    for _, d := range m.decls {
        s += d.String() + "\n"
    }
    return s
}

