package iCPSDL

import "strconv"
import "iCPSDL/parser"

type ontologyVisitor struct {
    parser.BaseontologyVisitor
}

func (this *ontologyVisitor) VisitProgram(ctx *parser.ProgramContext) []expression {
    expressions := make([]expression, 0)
    for _, s := range ctx.AllExpression() {
        expr := this.VisitExpression(s.(*parser.ExpressionContext))
        expressions = append(expressions, expr)
    }
    if ctx.Exit() != nil {
        exit := this.VisitExit(ctx.Exit().(*parser.ExitContext))
        expressions = append(expressions, exit)
    }
    return expressions
}

func (this *ontologyVisitor) VisitExpression(ctx *parser.ExpressionContext) expression {
    if ctx.Domain() != nil {
        return this.VisitDomain(ctx.Domain().(*parser.DomainContext))
    } else if ctx.Process() != nil {
        return this.VisitProcess(ctx.Process().(*parser.ProcessContext))
    } else if ctx.Local_configuration() != nil {
        return this.VisitLocal_configuration(ctx.Local_configuration().(*parser.Local_configurationContext))
    } else if ctx.Repository() != nil {
        return this.VisitRepository(ctx.Repository().(*parser.RepositoryContext))
    } else if ctx.Global_configuration() != nil {
        return this.VisitGlobal_configuration(ctx.Global_configuration().(*parser.Global_configurationContext))
    } else if ctx.Load() != nil {
        return this.VisitLoad(ctx.Load().(*parser.LoadContext))
    } else if ctx.Mermaid() != nil {
        return this.VisitMermaid(ctx.Mermaid().(*parser.MermaidContext))
    } else if ctx.Assignment() != nil {
        return this.VisitAssignment(ctx.Assignment().(*parser.AssignmentContext))
    } else if ctx.Remove() != nil {
        return this.VisitRemove(ctx.Remove().(*parser.RemoveContext))
    } else if ctx.Show() != nil {
        return this.VisitShow(ctx.Show().(*parser.ShowContext))
    } else if ctx.Translate() != nil {
        return this.VisitTranslate(ctx.Translate().(*parser.TranslateContext))
    } else if ctx.Traverse() != nil {
        return this.VisitTraverse(ctx.Traverse().(*parser.TraverseContext))
    } else if ctx.Configure() != nil {
        return this.VisitConfigure(ctx.Configure().(*parser.ConfigureContext))
    } else if ctx.Compose() != nil {
        return this.VisitCompose(ctx.Compose().(*parser.ComposeContext))
    } else if ctx.Project() != nil {
        return this.VisitProject(ctx.Project().(*parser.ProjectContext))
    } else if ctx.Info() != nil {
        return this.VisitInfo(ctx.Info().(*parser.InfoContext))
    } else if ctx.LSQ() != nil {
        line := ctx.INT().GetSymbol().GetLine()
        expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
        index, _ := strconv.Atoi(ctx.INT().GetText())
        return newAccess(expr, index, line)
    } else if ctx.ID() != nil {
        line := ctx.ID().GetSymbol().GetLine()
        id := ctx.ID().GetText()
        return newIdExpression(id, line)
    } else if ctx.Expression() != nil {
        return this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    }
    return nil
}

func (this *ontologyVisitor) VisitExit(ctx *parser.ExitContext) *exitExpression {
  line := 0
  if ctx.EXIT() != nil {
    line = ctx.EXIT().GetSymbol().GetLine()
    } else if ctx.QUIT() != nil {
      line = ctx.QUIT().GetSymbol().GetLine()
    }
    return newExit(line)
  }

func (this *ontologyVisitor) VisitLoad(ctx *parser.LoadContext) *loadExpression {
  line := ctx.LOAD().GetSymbol().GetLine()
  name := ""
  if ctx.ID() != nil {
      name = ctx.ID().GetText()
  } else if ctx.PATH() != nil {
      name = ctx.PATH().GetText()
  }
  return newLoad(name, line)
}

func (this *ontologyVisitor) VisitMermaid(ctx *parser.MermaidContext) *mermaidExpression {
  line := ctx.MERMAID().GetSymbol().GetLine()
  expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
  return newMermaid(expr, line)
}

func (this *ontologyVisitor) VisitAssignment(ctx *parser.AssignmentContext) *assignmentExpression {
    line := ctx.ASGN().GetSymbol().GetLine()
    id := ctx.ID().GetText()
    expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    return newAssignment(id, expr, line)
}

func (this *ontologyVisitor) VisitRemove(ctx *parser.RemoveContext) *removeExpression {
    line := ctx.REMOVE().GetSymbol().GetLine()
    id := ctx.ID().GetText()
    return newRemove(id, line)
}

func (this *ontologyVisitor) VisitShow(ctx *parser.ShowContext) *showExpression {
    line := 0
    if ctx.SHOW() != nil {
        line = ctx.SHOW().GetSymbol().GetLine()
    } else if ctx.LS() != nil {
        line = ctx.LS().GetSymbol().GetLine()
    }
    var expr expression = nil
    if ctx.Expression() != nil {
        expr = this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    }
    return newShow(expr, line)
}

func (this *ontologyVisitor) VisitTranslate(ctx *parser.TranslateContext) *translateExpression {
    line := ctx.TRANSLATE().GetSymbol().GetLine()
    expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    return newTranslate(expr, line)
}

func (this *ontologyVisitor) VisitTraverse(ctx *parser.TraverseContext) *traverseExpression {
    line := ctx.TRAVERSE().GetSymbol().GetLine()
    node := ctx.ID(0).GetText() + "." + ctx.ID(1).GetText()
    expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    return newTraverse(node, expr, line)
}

func (this *ontologyVisitor) VisitConfigure(ctx *parser.ConfigureContext) *configureExpression {
    line := ctx.CONFIGURE().GetSymbol().GetLine()
    rep := this.VisitExpression(ctx.Expression(0).(*parser.ExpressionContext))
    trav := this.VisitExpression(ctx.Expression(1).(*parser.ExpressionContext))
    controller := ctx.ID(0).GetText()
    actuator := ctx.ID(1).GetText()
    return newConfigure(rep, trav, controller, actuator, line)
}

func (this *ontologyVisitor) VisitCompose(ctx *parser.ComposeContext) *composeExpression {
    line := ctx.COMPOSE().GetSymbol().GetLine()
    expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    return newCompose(expr, line)
}

func (this *ontologyVisitor) VisitProject(ctx *parser.ProjectContext) *projectExpression {
    line := ctx.PROJECT().GetSymbol().GetLine()
    expr := this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    return newProject(expr, line)
}

func (this *ontologyVisitor) VisitInfo(ctx *parser.InfoContext) *infoExpression {
    line := ctx.INFO().GetSymbol().GetLine()
    var expr expression = nil
    if ctx.Expression() != nil {
        expr = this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    }
    return newInfo(expr, line)
}

/*****************
  domain
 ****************/

func (this *ontologyVisitor) VisitDomain(ctx *parser.DomainContext) *domain {
    line := ctx.DOMAIN().GetSymbol().GetLine()
    dmn := newDomain(line)
    for _, decl := range ctx.AllDomain_declaration() {
        this.VisitDomain_declaration(decl.(*parser.Domain_declarationContext), dmn)
    }
    return dmn
}

func (this *ontologyVisitor) VisitDomain_declaration(ctx *parser.Domain_declarationContext, dmn *domain) *domain {
    if ctx.Property() != nil {
        this.VisitProperty(ctx.Property().(*parser.PropertyContext), dmn)
    } else if ctx.Model() != nil {
        this.VisitModel(ctx.Model().(*parser.ModelContext), dmn)
    } else if ctx.Class() != nil {
        this.VisitClass(ctx.Class().(*parser.ClassContext), dmn)
    } else if ctx.Translation() != nil {
        this.VisitTranslation(ctx.Translation().(*parser.TranslationContext), dmn)
    }
	  return dmn
}

func (this *ontologyVisitor) VisitProperty(ctx *parser.PropertyContext, dmn *domain) *domain {
    for _, t := range ctx.AllType() {
        prop := this.VisitType(t.(*parser.TypeContext))
        dmn.add(prop)
    }

  return dmn
}

func (this *ontologyVisitor) VisitType(ctx *parser.TypeContext) imodel {
    if ctx.LBRA() != nil {
        line := ctx.LBRA().GetSymbol().GetLine()
        name := ctx.ID(0).GetText()
        enum := newEnumerationProperty(name, line)
        for i, id := range ctx.AllID() {
            if i == 0 {
                continue
            }
            enum.addLabel(id.GetText())
        }
        return enum
    } else if len(ctx.AllID()) == 1 {
        line := ctx.ID(0).GetSymbol().GetLine()
        return newProperty(ctx.ID(0).GetText(), line)
    }

	  return nil
}

func (this *ontologyVisitor) VisitModel(ctx *parser.ModelContext, dmn *domain) *domain {
    for _, id := range ctx.AllID() {
        line := id.GetSymbol().GetLine()
        mdl := newModel(id.GetText(), line)
        dmn.add(mdl)
    }
    return dmn
}

func (this *ontologyVisitor) VisitClass(ctx *parser.ClassContext, dmn *domain) *domain {
    name := ctx.ID(0).GetText()
    nodes := make([]string, 0)
    for i := 1; i < len(ctx.AllID()); i++ {
        nodes = append(nodes, ctx.ID(i).GetText())
    }
    internalEdges := make([]*internalEdge, 0)
    for _, ie := range ctx.AllInternal_edge() {
        c := this.VisitInternal_edge(ie.(*parser.Internal_edgeContext))
        internalEdges = append(internalEdges, c)
    }
    var decl domainDecl = nil
    if ctx.PHYSICAL() != nil {
        line := ctx.PHYSICAL().GetSymbol().GetLine()
        decl = newPhysical(name, nodes, internalEdges, line)
    } else if ctx.ACTUATOR() != nil {
        line := ctx.ACTUATOR().GetSymbol().GetLine()
        decl = newActuator(name, nodes, internalEdges, line)
    }
    dmn.add(decl)
    return dmn
}

func (this *ontologyVisitor) VisitInternal_edge(ctx *parser.Internal_edgeContext) *internalEdge {
    line := ctx.ARROW().GetSymbol().GetLine()
    node1 := ctx.ID(0).GetText()
    node2 := ctx.ID(1).GetText()
    return newInternalEdge(node1, node2, line)
}

func (this *ontologyVisitor) VisitTranslation(ctx *parser.TranslationContext, dmn *domain) *domain {
    line := ctx.TRANSLATION().GetSymbol().GetLine()
    node1 := ctx.ID(0).GetText()
    node2 := ctx.ID(1).GetText()
    targets := make([]*externalEdge, 0)
    for _, aconn := range ctx.AllArg_connection() {
        ac := this.VisitArg_connection(aconn.(*parser.Arg_connectionContext))
        targets = append(targets, ac)
    }
    trans := newTranslation(node1, node2, targets, line)
    dmn.add(trans)
    return dmn
}

func (this *ontologyVisitor) VisitArg_connection(ctx *parser.Arg_connectionContext) *externalEdge {
    line := ctx.ARROW().GetSymbol().GetLine()
    c1 := ctx.ID(0).GetText()
    m1 := ctx.ID(1).GetText()
    c2 := ctx.ID(2).GetText()
    m2 := ctx.ID(3).GetText()
    return newExternalEdge(c1, m1, c2, m2, line)
}

func (this *ontologyVisitor) VisitConnection(ctx *parser.ConnectionContext) *connection {
    line := ctx.ARROW().GetSymbol().GetLine()
    node1 := ctx.ID(0).GetText()
    node2 := ctx.ID(1).GetText()
    return newConnection(node1, node2, line)
}

/*****************
  process
 ****************/

func (this *ontologyVisitor) VisitProcess(ctx *parser.ProcessContext) *process {
    line := ctx.PROCESS().GetSymbol().GetLine()
    prdgm := ctx.ID().GetText()
  	proc := newProcess(prdgm, line)
    for _, decl := range ctx.AllProcess_declaration() {
        this.VisitProcess_declaration(decl.(*parser.Process_declarationContext), proc)
    }
    return proc
}

func (this *ontologyVisitor) VisitProcess_declaration(ctx *parser.Process_declarationContext, proc *process) *process {
    if ctx.Device() != nil {
        this.VisitDevice(ctx.Device().(*parser.DeviceContext), proc)
    } else if ctx.Component() != nil {
        this.VisitComponent(ctx.Component().(*parser.ComponentContext), proc)
    } else if ctx.Connection_decl() != nil {
        this.VisitConnection_decl(ctx.Connection_decl().(*parser.Connection_declContext), proc)
    }
    return proc
}

func (this *ontologyVisitor) VisitDevice(ctx *parser.DeviceContext, proc *process) *process {
    line := ctx.DEVICE().GetSymbol().GetLine()
    for _, id := range ctx.AllID() {
        h := newHardware(id.GetText(), line)
        proc.add(h)
    }
    return proc
}

func (this *ontologyVisitor) VisitComponent(ctx *parser.ComponentContext, proc *process) *process {
    if ctx.PHYSICAL() != nil {
        line := ctx.PHYSICAL().GetSymbol().GetLine()
        length := len(ctx.AllID())
        class := ctx.ID(length - 1).GetText()
        for i := 0; i < length - 1; i++ {
            comp := newPhysicalComp(ctx.ID(i).GetText(), class, line)
            proc.add(comp)
        }
    } else if ctx.ACTUATOR() != nil {
        line := ctx.ACTUATOR().GetSymbol().GetLine()
        length := len(ctx.AllID())
        act := ctx.ID(length - 1).GetText()
        for i := 0; i < length - 1; i++ {
            name := ctx.ID(i).GetText()
            i += 1
            hard := ctx.ID(i).GetText()
            comp := newActuatorComp(name, hard, act, line)
            proc.add(comp)
      }
    } else if ctx.SENSOR() != nil {
        line := ctx.SENSOR().GetSymbol().GetLine()
        length := len(ctx.AllID())
        act := ctx.ID(length - 1).GetText()
        for i := 0; i < length - 1; i++ {
            name := ctx.ID(i).GetText()
            i += 1
            hard := ctx.ID(i).GetText()
            comp := newSensorComp(name, hard, act, line)
            proc.add(comp)
      }
    }
    return proc
}

func (this *ontologyVisitor) VisitConnection_decl(ctx *parser.Connection_declContext, proc *process) *process {
    for _, conn := range ctx.AllConnection() {
        c := this.VisitConnection(conn.(*parser.ConnectionContext))
        proc.add(c)
    }
    return proc
}

/*****************
  local configuration
 ****************/

func (this *ontologyVisitor) VisitLocal_configuration(ctx *parser.Local_configurationContext) *local {
    line := ctx.LOCAL().GetSymbol().GetLine()
    lc := newLocal(line)
    for i, l := range ctx.AllLocal() {
        sess := this.VisitLocal(l.(*parser.LocalContext))
        id := ctx.ID(i).GetText()
        lc.add(id, sess)
    }
    return lc
}

func (this *ontologyVisitor) VisitLocal(ctx *parser.LocalContext) session {
    if ctx.END() != nil {
        line := ctx.END().GetSymbol().GetLine()
        return newEnd(line)
    } else if ctx.Send() != nil {
        return this.VisitSend(ctx.Send().(*parser.SendContext))
    } else if ctx.Receive() != nil {
        return this.VisitReceive(ctx.Receive().(*parser.ReceiveContext))
    } else if ctx.Loop() != nil {
        return this.VisitLoop(ctx.Loop().(*parser.LoopContext))
    } else if ctx.Label() != nil {
        return this.VisitLabel(ctx.Label().(*parser.LabelContext))
    } else if ctx.Select() != nil {
        return this.VisitSelect(ctx.Select().(*parser.SelectContext))
    } else if ctx.Branch() != nil {
        return this.VisitBranch(ctx.Branch().(*parser.BranchContext))
    }
	return nil
}

func (this *ontologyVisitor) VisitSend(ctx *parser.SendContext) *send {
    line := ctx.BANG().GetSymbol().GetLine()
    to := ctx.ID().GetText()
    tdef := this.VisitType(ctx.Type().(*parser.TypeContext))
    var cont session = nil
    //if ctx.Local() != nil {
    cont = this.VisitLocal(ctx.Local().(*parser.LocalContext))
    //}
    s := newSend(to, tdef, cont, line)
    return s
}

func (this *ontologyVisitor) VisitReceive(ctx *parser.ReceiveContext) *receive {
    line := ctx.QUESTION().GetSymbol().GetLine()
    from := ctx.ID().GetText()
    tdef := this.VisitType(ctx.Type().(*parser.TypeContext))
    var cont session = nil
    //if ctx.Local() != nil {
    cont = this.VisitLocal(ctx.Local().(*parser.LocalContext))
    //}
    r := newReceive(from, tdef, cont, line)
    return r
}

func (this *ontologyVisitor) VisitLoop(ctx *parser.LoopContext) *loop {
    label := ctx.ID().GetText()
    line := ctx.DOT().GetSymbol().GetLine()
    s := this.VisitLocal(ctx.Local().(*parser.LocalContext))
    return newLoop(label, s, line)
}

func (this *ontologyVisitor) VisitLabel(ctx *parser.LabelContext) *labelVar {
  label := ctx.ID().GetText()
  line := ctx.ID().GetSymbol().GetLine()
  return newLabel(label, line)
}

func (this *ontologyVisitor) VisitSelect(ctx *parser.SelectContext) *sel {
    line := ctx.BANG().GetSymbol().GetLine()
    adv := ctx.ID(0).GetText()
    enum := ctx.ID(1).GetText()
    sl := newSelect(adv, enum, line)
    for i, s := range ctx.AllLocal() {
        ss := this.VisitLocal(s.(*parser.LocalContext))
        e := ctx.ID(i + 2).GetText()
        sl.addSession(e, ss)
    }
    return sl
}

func (this *ontologyVisitor) VisitBranch(ctx *parser.BranchContext) *branch {
    line := ctx.QUESTION().GetSymbol().GetLine()
    adv := ctx.ID(0).GetText()
    enum := ctx.ID(1).GetText()
    bra := newBranch(adv, enum, line)
    for i, s := range ctx.AllLocal() {
        ss := this.VisitLocal(s.(*parser.LocalContext))
        e := ctx.ID(i + 2).GetText()
        bra.addSession(e, ss)
    }
    return bra
}

/*****************
  global configuration
 ****************/

func (this *ontologyVisitor) VisitGlobal_configuration(ctx *parser.Global_configurationContext) *global {
   line := ctx.GLOBAL().GetSymbol().GetLine()
   session := this.VisitGlobal(ctx.Global().(*parser.GlobalContext))
   return newGlobal(session, line)
}

func (this *ontologyVisitor) VisitGlobal(ctx *parser.GlobalContext) session {
    if ctx.END() != nil {
        line := ctx.END().GetSymbol().GetLine()
        return newEnd(line)
    } else if ctx.Pass() != nil {
        return this.VisitPass(ctx.Pass().(*parser.PassContext))
    } else if ctx.Global_loop() != nil {
        return this.VisitGlobal_loop(ctx.Global_loop().(*parser.Global_loopContext))
    } else if ctx.Label() != nil {
        return this.VisitLabel(ctx.Label().(*parser.LabelContext))
    } else if ctx.Choice() != nil {
        return this.VisitChoice(ctx.Choice().(*parser.ChoiceContext))
    }
	return nil
}

func (this *ontologyVisitor) VisitPass(ctx *parser.PassContext) *pass {
    line := ctx.ARROW().GetSymbol().GetLine()
    from := ctx.ID(0).GetText()
    to := ctx.ID(1).GetText()
    tdef := this.VisitType(ctx.Type().(*parser.TypeContext))
    var cont session = nil
    if ctx.Global() != nil {
        cont = this.VisitGlobal(ctx.Global().(*parser.GlobalContext))
    }
    return newPass(from, to, tdef, cont, line)
}

func (this *ontologyVisitor) VisitGlobal_loop(ctx *parser.Global_loopContext) *loop {
    label := ctx.ID().GetText()
    line := ctx.DOT().GetSymbol().GetLine()
    s := this.VisitGlobal(ctx.Global().(*parser.GlobalContext))
    return newLoop(label, s, line)
}

func (this *ontologyVisitor) VisitChoice(ctx *parser.ChoiceContext) *choice {
    line := ctx.ARROW().GetSymbol().GetLine()
    from := ctx.ID(0).GetText()
    to := ctx.ID(1).GetText()
    enum := ctx.ID(2).GetText()
    ch := newChoice(from, to, enum, line)
    for i, s := range ctx.AllGlobal() {
        ss := this.VisitGlobal(s.(*parser.GlobalContext))
        e := ctx.ID(i + 3).GetText()
        ch.addSession(e, ss)
    }
    return ch
}

/*****************
  repository
 ****************/

func (this *ontologyVisitor) VisitRepository(ctx *parser.RepositoryContext) *repository {
    line := ctx.REPOSITORY().GetSymbol().GetLine()
    prdgm := ctx.ID().GetText()
  	rep := newRepository(prdgm, line)
    for _, decl := range ctx.AllRepository_decl() {
        elem := this.VisitRepository_decl(decl.(*parser.Repository_declContext))
        rep.add(elem)
    }
    return rep
}

func (this *ontologyVisitor) VisitRepository_decl(ctx *parser.Repository_declContext) irole {
    class := ctx.ID(0).GetText()
    name := ctx.ID(1).GetText()
    sess := this.VisitLocal(ctx.Local().(*parser.LocalContext))
    if ctx.ESTIMATE() != nil {
        line := ctx.ESTIMATE().GetSymbol().GetLine()
        return newEstimate(class, name, sess, line)
    } else if ctx.SENSE() != nil {
        line := ctx.SENSE().GetSymbol().GetLine()
        return newSense(class, name, sess, line)
    } else if ctx.CONTROL() != nil {
        line := ctx.CONTROL().GetSymbol().GetLine()
        return newControl(class, name, sess, line)
    } else if ctx.ACTUATE() != nil {
        line := ctx.ACTUATE().GetSymbol().GetLine()
        return newActuate(class, name, sess, line)
    }
	return nil
}
