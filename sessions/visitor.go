package ast

import (
    "sessions/parser"
)

type sessionsVisitor struct {
    parser.BasesessionsVisitor
}

func (this *sessionsVisitor) VisitModule(ctx *parser.ModuleContext) *module {
    imp := this.VisitImports(ctx.Imports().(*parser.ImportsContext))
    mod := newModule(imp, 0)

    for _, decl := range ctx.AllDeclaration() {
        d := this.VisitDeclaration(decl.(*parser.DeclarationContext))
        mod.addDeclaration(d)
    }
    return mod
}

func (this *sessionsVisitor) VisitImports(ctx *parser.ImportsContext) []string {
    // imports := make([]string, len(ctx.AllID()))
    // for i, imp := range ctx.AllID() {
    //     imports[i] = imp.GetText()
    // }
    imports := make([]string, len(ctx.AllPath()))
    for i, path := range ctx.AllPath() {
        imports[i] = this.VisitPath(path.(*parser.PathContext))
    }
    return imports
}

func (v *sessionsVisitor) VisitPath(ctx *parser.PathContext) string {
    path := ""
    for i, id := range ctx.AllID() {
        if i != 0 {
            path += "."
        }
        path += id.GetText()
    }
    return path
}

func (this *sessionsVisitor) VisitDeclaration(ctx *parser.DeclarationContext) declaration {
    if ctx.DeclarationSugar() != nil {
        return this.VisitDeclarationSugar(ctx.DeclarationSugar().(*parser.DeclarationSugarContext))
    }

    id, _ := this.VisitName(ctx.Name().(*parser.NameContext))
    if ctx.VAL() != nil {
        line := ctx.VAL().GetSymbol().GetLine()
        if ctx.Abstraction() != nil {
            abstr := this.VisitAbstraction(ctx.Abstraction().(*parser.AbstractionContext))
            return newAbstractionDeclaration(id, abstr, line)
        } else if ctx.PrimaryTerm() != nil {
            term := this.VisitPrimaryTerm(ctx.PrimaryTerm().(*parser.PrimaryTermContext))
            return newValueDeclaration(id, term, line)
        }
    } else if ctx.TYPE() != nil {
        line := ctx.TYPE().GetSymbol().GetLine()
        var tdef typedef
        if ctx.Session() != nil {
            tdef = this.VisitSession(ctx.Session().(*parser.SessionContext))
        } else if ctx.Configuration() != nil {
            tdef = this.VisitConfiguration(ctx.Configuration().(*parser.ConfigurationContext))
        } else if ctx.RecordType() != nil {
            tdef = this.VisitRecordType(ctx.RecordType().(*parser.RecordTypeContext))
        } else if ctx.BrokerType() != nil {
            tdef = this.VisitBrokerType(ctx.BrokerType().(*parser.BrokerTypeContext))
        }
        return newTypeDeclaration(id, tdef, line)
    } else if ctx.COLON() != nil {
        line := ctx.COLON().GetSymbol().GetLine()
        tdef := this.VisitSessionDef(ctx.SessionDef().(*parser.SessionDefContext))
        return newSessionAssignment(id, tdef, line)
    }
    return nil
}

func (this *sessionsVisitor) VisitDeclarationSugar(ctx *parser.DeclarationSugarContext) declaration {
    id, _ := this.VisitName(ctx.Name().(*parser.NameContext))
    if ctx.PROC() != nil {
        pars := make([]*participantExpr, len(ctx.AllParticipant()))
        for i, p := range ctx.AllParticipant() {
            pars[i] = this.VisitParticipant(p.(*parser.ParticipantContext))
        }
        proc := this.VisitProcess(ctx.Process().(*parser.ProcessContext))
        line := ctx.PROC().GetSymbol().GetLine()
        abstr := newAbstraction(proc, pars, line)
        return newAbstractionDeclaration(id, abstr, line)
    } else if ctx.RECORD() != nil {
        line := ctx.RECORD().GetSymbol().GetLine()
        rec := this.VisitPrimaryRecord(ctx.PrimaryRecord().(*parser.PrimaryRecordContext))
        return newValueDeclaration(id, rec, line)
    } else if ctx.BROKER() != nil {
        line := ctx.BROKER().GetSymbol().GetLine()
        broker := this.VisitBroker(ctx.Broker().(*parser.BrokerContext))
        return newValueDeclaration(id, broker, line)
    }
    return nil
}

func (this *sessionsVisitor) VisitName(ctx *parser.NameContext) (string, int) {
    return ctx.ID().GetText(), ctx.ID().GetSymbol().GetLine()
}

func (this *sessionsVisitor) VisitParticipant(ctx *parser.ParticipantContext) *participantExpr {
    return newParticipantExpr(ctx.ID().GetText(), ctx.ID().GetSymbol().GetLine())
}

func (this *sessionsVisitor) VisitVariableDef(ctx *parser.VariableDefContext) *variableExpr {
    id, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
    tdef := this.VisitType(ctx.Type_().(*parser.TypeContext))
    return newVariableExpr(id, tdef, line)
}

//TODO redundant <- same as VisitName
func (this *sessionsVisitor) VisitVariable(ctx *parser.VariableContext) (string, int) {
    return ctx.ID().GetText(), ctx.ID().GetSymbol().GetLine()
}

func (this *sessionsVisitor) VisitAbstraction(ctx *parser.AbstractionContext) *abstraction {
    pars := make([]*participantExpr, len(ctx.AllParticipant()))
    for i, p := range ctx.AllParticipant() {
        pars[i] = this.VisitParticipant(p.(*parser.ParticipantContext))
    }

    proc := this.VisitProcess(ctx.Process().(*parser.ProcessContext))
    line := ctx.DOT().GetSymbol().GetLine()
    return newAbstraction(proc, pars, line)
}

func (this *sessionsVisitor) VisitApplication(ctx *parser.ApplicationContext) *application {
    args := make([]*participantExpr, len(ctx.AllParticipant()))
    for i, a := range ctx.AllParticipant() {
        args[i] = this.VisitParticipant(a.(*parser.ParticipantContext))
    }

    adef := this.VisitTerm(ctx.Term().(*parser.TermContext))
    line := ctx.Term().GetStart().GetLine()

    return newApplication(adef, args, line)
}

func (this *sessionsVisitor) VisitConcurrent(ctx *parser.ConcurrentContext) (*participantExpr, []*participantExpr, []*application, int) {
    var participant *participantExpr = nil
    if ctx.Name() != nil {
        id, line := this.VisitName(ctx.Name().(*parser.NameContext))
        participant = newParticipantExpr(id, line)
    }
    participants := make([]*participantExpr, len(ctx.AllApplication()))
    applications := make([]*application, len(ctx.AllApplication()))
    for i, r := range ctx.AllApplication() {
        participants[i] = this.VisitParticipant(ctx.Participant(i).(*parser.ParticipantContext))
        applications[i] = this.VisitApplication(r.(*parser.ApplicationContext))
    }
    line := ctx.CONC().GetSymbol().GetLine()
    return participant, participants, applications, line
}

func (this *sessionsVisitor) VisitRequest(ctx *parser.RequestContext) (*variableExpr, *participantExpr, []*participantExpr, []*application, int) {
    vrb, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
    channel := newVariableExpr(vrb, nil, line)
    id, line := this.VisitName(ctx.Name().(*parser.NameContext))
    participant := newParticipantExpr(id, line)

    participants := make([]*participantExpr, len(ctx.AllApplication()))
    applications := make([]*application, len(ctx.AllApplication()))
    for i, r := range ctx.AllApplication() {
        participants[i] = this.VisitParticipant(ctx.Participant(i).(*parser.ParticipantContext))
        applications[i] = this.VisitApplication(r.(*parser.ApplicationContext))
    }
    line = ctx.REQ().GetSymbol().GetLine()
    return channel, participant, participants, applications, line
}

func (this *sessionsVisitor) VisitAcceptProc(ctx *parser.AcceptProcContext) *acceptProc {
    vrb, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
    channel := newVariableExpr(vrb, nil, line)
    participants := make([]*participantExpr, len(ctx.AllParticipant()))
    applications := make([]*application, len(ctx.AllApplication()))
    for i, r := range ctx.AllApplication() {
        participants[i] = this.VisitParticipant(ctx.Participant(i).(*parser.ParticipantContext))
        applications[i] = this.VisitApplication(r.(*parser.ApplicationContext))
    }
    line = ctx.ACC().GetSymbol().GetLine()
    return newAcceptProc(channel, participants, applications, line)
}

func (this *sessionsVisitor) VisitRecurse(ctx *parser.RecurseContext) *recurseProc {
    vrb, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
    channel := newVariableExpr(vrb, nil, line)
    participants := make([]*participantExpr, len(ctx.AllParticipant()) - 1)
    applications := make([]*application, len(ctx.AllApplication()) - 1)
    var rec *participantExpr
    var app *application
    for i, r := range ctx.AllApplication() {
        p := this.VisitParticipant(ctx.Participant(i).(*parser.ParticipantContext))
        a := this.VisitApplication(r.(*parser.ApplicationContext))
        if i == 0 {
            rec = p
            app = a

        } else {
            participants[i - 1] = p
            applications[i - 1] = a
        }
    }
    line = ctx.REC().GetSymbol().GetLine()
    return newRecurseProc(channel, rec, app, participants, applications, line)
}

func (this *sessionsVisitor) VisitProcess(ctx *parser.ProcessContext) process {
    if ctx.BlockProc() != nil {
        return this.VisitBlockProc(ctx.BlockProc().(*parser.BlockProcContext))
    } else if ctx.SequenceProc() != nil {
        return this.VisitSequenceProc(ctx.SequenceProc().(*parser.SequenceProcContext))
    } else if ctx.IfThenElseProc() != nil {
        return this.VisitIfThenElseProc(ctx.IfThenElseProc().(*parser.IfThenElseProcContext))
    } else if ctx.SelectProc() != nil {
        return this.VisitSelectProc(ctx.SelectProc().(*parser.SelectProcContext))
    } else if ctx.BranchProc() != nil {
        return this.VisitBranchProc(ctx.BranchProc().(*parser.BranchProcContext))
    } else if ctx.AcceptProc() != nil {
        return this.VisitAcceptProc(ctx.AcceptProc().(*parser.AcceptProcContext))
    } else if ctx.Recurse() != nil {
        return this.VisitRecurse(ctx.Recurse().(*parser.RecurseContext))
    } else if ctx.TerminateProc() != nil {
        return this.VisitTerminateProc(ctx.TerminateProc().(*parser.TerminateProcContext))
    }
    return nil
}

func (this *sessionsVisitor) VisitBlockProc(ctx *parser.BlockProcContext) process {
    return this.VisitProcess(ctx.Process().(*parser.ProcessContext))
}

func (this *sessionsVisitor) VisitSequenceProc(ctx *parser.SequenceProcContext) process {
    var proc process = nil
    if ctx.Process() != nil {
        proc = this.VisitProcess(ctx.Process().(*parser.ProcessContext))
    }

    if ctx.Prefix() != nil {
        expr, line := this.VisitPrefix(ctx.Prefix().(*parser.PrefixContext))
        if proc == nil {
            proc = newTerminate(line)
        }
        return newSequentialProc(expr, proc, line)
    } else if ctx.Concurrent() != nil {
        participant, participants, applications, line := this.VisitConcurrent(ctx.Concurrent().(*parser.ConcurrentContext))
        if proc == nil {
            proc = newTerminate(line)
        }
        return newIntroProc(participant, participants, applications, proc, line)
    } else if ctx.Request() != nil {
        channel, participant, participants, applications, line := this.VisitRequest(ctx.Request().(*parser.RequestContext))
        if proc == nil {
            proc = newTerminate(line)
        }
        return newRequestProc(channel, participant, participants, applications, proc, line)
    } else if ctx.Call() != nil {
        exprs, adef, variables, line:= this.VisitCall(ctx.Call().(*parser.CallContext))
        if proc == nil {
            proc = newTerminate(line)
        }
        return newCall(exprs, adef, variables, proc, line)
    }

    return nil
}

func (this *sessionsVisitor) VisitPrefix(ctx *parser.PrefixContext) (seqExpr, int) {
    line := 0
    if ctx.Receive() != nil {
        line = ctx.Receive().GetStart().GetLine()
        return this.VisitReceive(ctx.Receive().(*parser.ReceiveContext)), line
    } else if ctx.Let() != nil {
        line = ctx.Let().GetStart().GetLine()
        return this.VisitLet(ctx.Let().(*parser.LetContext)), line
    } else if ctx.Out() != nil {
        line = ctx.Out().GetStart().GetLine()
        return this.VisitOut(ctx.Out().(*parser.OutContext)), line
    } else if ctx.Inp() != nil {
        line = ctx.Inp().GetStart().GetLine()
        return this.VisitInp(ctx.Inp().(*parser.InpContext)), line
    } else if ctx.CloseBroker() != nil {
        line := ctx.CloseBroker().GetStart().GetLine()
        return this.VisitCloseBroker(ctx.CloseBroker().(*parser.CloseBrokerContext)), line
    }
    return nil, line
}

func (this *sessionsVisitor) VisitLet(ctx *parser.LetContext) *letExpr {
    id, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
    var tdef typedef = nil
    if ctx.Type_() != nil {
        tdef = this.VisitType(ctx.Type_().(*parser.TypeContext))
    }
    variable := newVariableExpr(id, tdef, line)
    expr := this.VisitLogicalOr(ctx.LogicalOr().(*parser.LogicalOrContext))
    line = ctx.LET().GetSymbol().GetLine()
    return newLetExpr(variable, expr, line)
}

func (this *sessionsVisitor) VisitCloseBroker(ctx *parser.CloseBrokerContext) *closeExpr {
    id, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
    variable := newVariableExpr(id, nil, line)
    line = ctx.ARROW().GetSymbol().GetLine()
    return newCloseExpr(variable, line)
}

func (this *sessionsVisitor) VisitIfThenElseProc(ctx *parser.IfThenElseProcContext) process {
    cond := this.VisitLogicalOr(ctx.LogicalOr().(*parser.LogicalOrContext))
    thenProc := this.VisitProcess(ctx.Process(0).(*parser.ProcessContext))
    elseProc := this.VisitProcess(ctx.Process(1).(*parser.ProcessContext))
    line := ctx.IF().GetSymbol().GetLine()
    return newIfThenElseProc(cond, thenProc, elseProc, line)
}

func (this *sessionsVisitor) VisitSelectProc(ctx *parser.SelectProcContext) process {
    line := 0
    participant := this.VisitParticipant(ctx.Participant().(*parser.ParticipantContext))
    conds := make([]expression, len(ctx.AllLogicalOr()))
    labels := make([]*labelExpr, len(ctx.AllLabelExpression()))
    processes := make([]process, len(ctx.AllProcess()))
    if ctx.SELECT() != nil {
        line = ctx.SELECT().GetSymbol().GetLine()
        for i := range ctx.AllLogicalOr() {
            conds[i] = this.VisitLogicalOr(ctx.LogicalOr(i).(*parser.LogicalOrContext))
            labels[i] = this.VisitLabelExpression(ctx.LabelExpression(i).(*parser.LabelExpressionContext))
            processes[i] = this.VisitProcess(ctx.Process(i).(*parser.ProcessContext))
        }
        psize := len(ctx.AllProcess()) - 1
        labels[psize] = this.VisitLabelExpression(ctx.LabelExpression(psize).(*parser.LabelExpressionContext))
        processes[psize] = this.VisitProcess(ctx.Process(psize).(*parser.ProcessContext))
    } else /*if ctx.ARROW() != nil*/ {
        line = ctx.ARROW(0).GetSymbol().GetLine()
        labels[0] = this.VisitLabelExpression(ctx.LabelExpression(0).(*parser.LabelExpressionContext))
        if len(ctx.AllProcess()) != 0 {
            processes[0] = this.VisitProcess(ctx.Process(0).(*parser.ProcessContext))
        } else {
            processes = append(processes, newTerminate(line))
        }
    }

    return newSelectProc(participant, conds, labels, processes, line)
}

func (this *sessionsVisitor) VisitBranchProc(ctx *parser.BranchProcContext) process {
    participant := this.VisitParticipant(ctx.Participant().(*parser.ParticipantContext))

    labels := make([]*labelExpr, len(ctx.AllLabelExpression()))
    processes := make([]process, len(ctx.AllProcess()))
    for i, p := range ctx.AllProcess() {
        labels[i] = this.VisitLabelExpression(ctx.LabelExpression(i).(*parser.LabelExpressionContext))
        processes[i] = this.VisitProcess(p.(*parser.ProcessContext))
    }
    line := ctx.BRANCH().GetSymbol().GetLine()
    return newBranchProc(participant, labels, processes, line)
}

func (this *sessionsVisitor) VisitTerminateProc(ctx *parser.TerminateProcContext) *terminate {
    return newTerminate(ctx.TERMINATE().GetSymbol().GetLine())
}

func (this *sessionsVisitor) VisitReceive(ctx *parser.ReceiveContext) seqExpr {
    var expr expression
    if ctx.Send() != nil {
        expr = this.VisitSend(ctx.Send().(*parser.SendContext))
    } else if ctx.Participant() != nil {
        expr = this.VisitParticipant(ctx.Participant().(*parser.ParticipantContext))
    }
    if ctx.ARROW() == nil {
        return expr.(seqExpr)
    }
    line  := ctx.ARROW().GetSymbol().GetLine()
    variable := this.VisitVariableDef(ctx.VariableDef().(*parser.VariableDefContext))
    return newReceiveExpr(variable, expr, line)
}

func (this *sessionsVisitor) VisitSend(ctx *parser.SendContext) seqExpr {
    var lexpr seqExpr = nil
    if ctx.Receive() != nil {
        lexpr = this.VisitReceive(ctx.Receive().(*parser.ReceiveContext))
    }
    for i, lor := range ctx.AllLogicalOr() {
        line := ctx.ARROW(i).GetSymbol().GetLine()
        rexpr := this.VisitLogicalOr(lor.(*parser.LogicalOrContext))
        if i == 0 && ctx.Participant() != nil {
            participant := this.VisitParticipant(ctx.Participant().(*parser.ParticipantContext))
            lexpr = newSendExpr(participant, rexpr, line)
        } else {
            lexpr = newSendExpr(lexpr, rexpr, line)
        }
    }
    return lexpr
}

func (this *sessionsVisitor) VisitOut(ctx *parser.OutContext) *outExpr {
    var ioConfig expression = nil
    //line := ctx.OUT().GetSymbol().GetLine()
    if ctx.Term() != nil {
        ioConfig = this.VisitTerm(ctx.Term().(*parser.TermContext))
        //ioConfig = newPort(ioConfig, line)
    }
    line := ctx.OUT().GetSymbol().GetLine()
    exprs := make([]expression, len(ctx.AllLogicalOr()))
    for i, lor := range ctx.AllLogicalOr() {
        exprs[i] = this.VisitLogicalOr(lor.(*parser.LogicalOrContext))
    }
    return newOutExpr(ioConfig, exprs, line)
}

func (this *sessionsVisitor) VisitInp(ctx *parser.InpContext) *inpExpr {
    var ioConfig expression = nil
    //line := ctx.INP().GetSymbol().GetLine()
    if ctx.Term() != nil {
        ioConfig = this.VisitTerm(ctx.Term().(*parser.TermContext))
        //ioConfig = newPort(ioConfig, line)
    }

    line := ctx.ARROW().GetSymbol().GetLine()
    variable := this.VisitVariableDef(ctx.VariableDef().(*parser.VariableDefContext))
    return newInpExpr(ioConfig, variable, line)
}

func (this *sessionsVisitor) VisitCall(ctx *parser.CallContext) ([]expression, expression, []*variableExpr, int) {
    exprs := make([]expression, len(ctx.AllLogicalOr()))
    for i, expr := range ctx.AllLogicalOr() {
        exprs[i] = this.VisitLogicalOr(expr.(*parser.LogicalOrContext))
    }

    adef := this.VisitTerm(ctx.Term().(*parser.TermContext))

    variables := make([]*variableExpr, len(ctx.AllVariableDef()))
    for i, variable := range ctx.AllVariableDef() {
        variables[i] = this.VisitVariableDef(variable.(*parser.VariableDefContext))
    }

    line := ctx.EQUALS().GetSymbol().GetLine()

    return exprs, adef, variables, line
}

func (this *sessionsVisitor) VisitLogicalOr(ctx *parser.LogicalOrContext) expression {
    var lexpr expression = nil
    line := 0
    if ctx.LogicalOr() != nil {
        lexpr = this.VisitLogicalOr(ctx.LogicalOr().(*parser.LogicalOrContext))
        line = ctx.LOR().GetSymbol().GetLine()
    }
    rexpr := this.VisitLogicalAnd(ctx.LogicalAnd().(*parser.LogicalAndContext))
    if lexpr != nil {
        return newLogicalExpr(lexpr, rexpr, ctx.LOR().GetText(), line)
    }
    return rexpr
}

func (this *sessionsVisitor) VisitLogicalAnd(ctx *parser.LogicalAndContext) expression {
    var lexpr expression = nil
    line := 0
    if ctx.LogicalAnd() != nil {
        lexpr = this.VisitLogicalAnd(ctx.LogicalAnd().(*parser.LogicalAndContext))
        line = ctx.LAND().GetSymbol().GetLine()
    }
    rexpr := this.VisitEqualityExpr(ctx.EqualityExpr().(*parser.EqualityExprContext))
    if lexpr != nil {
        return newLogicalExpr(lexpr, rexpr, ctx.LAND().GetText(), line)
    }
    return rexpr
}

func (this *sessionsVisitor) VisitEqualityExpr(ctx *parser.EqualityExprContext) expression {
    var lexpr expression = nil
    line := 0
    if ctx.EqualityExpr() != nil {
        lexpr = this.VisitEqualityExpr(ctx.EqualityExpr().(*parser.EqualityExprContext))
        line = ctx.EQOP().GetSymbol().GetLine()
    }
    rexpr := this.VisitRelationalExpr(ctx.RelationalExpr().(*parser.RelationalExprContext))
    if lexpr != nil {
        return newEqualityExpr(lexpr, rexpr, ctx.EQOP().GetText(), line)
    }
    return rexpr
}

func (this *sessionsVisitor) VisitRelationalExpr(ctx *parser.RelationalExprContext) expression {
    var lexpr expression = nil
    line := 0
    if ctx.RelationalExpr() != nil {
        lexpr = this.VisitRelationalExpr(ctx.RelationalExpr().(*parser.RelationalExprContext))
        line = ctx.RELOP().GetSymbol().GetLine()
    }
    rexpr := this.VisitSumExpr(ctx.SumExpr().(*parser.SumExprContext))
    if lexpr != nil {
        return newRelationalExpr(lexpr, rexpr, ctx.RELOP().GetText(), line)
    }
    return rexpr
}

func (this *sessionsVisitor) VisitSumExpr(ctx *parser.SumExprContext) expression {
    var lexpr expression = nil
    line := 0
    if ctx.SumExpr() != nil {
        lexpr = this.VisitSumExpr(ctx.SumExpr().(*parser.SumExprContext))
        line = ctx.SUMOP().GetSymbol().GetLine()
    }
    rexpr := this.VisitMultExpr(ctx.MultExpr().(*parser.MultExprContext))
    if lexpr != nil {
        return newSumExpr(lexpr, rexpr, ctx.SUMOP().GetText(), line)
    }
    return rexpr
}

func (this *sessionsVisitor) VisitMultExpr(ctx *parser.MultExprContext) expression {
    var lexpr expression = nil
    line := 0
    if ctx.MultExpr() != nil {
        lexpr = this.VisitMultExpr(ctx.MultExpr().(*parser.MultExprContext))
        line = ctx.MULOP().GetSymbol().GetLine()
    }
    rexpr := this.VisitUnary(ctx.Unary().(*parser.UnaryContext))
    if lexpr != nil {
        return newMultExpr(lexpr, rexpr, ctx.MULOP().GetText(), line)
    }
    return rexpr
}

func (this *sessionsVisitor) VisitUnary(ctx *parser.UnaryContext) expression {
    if ctx.Unary() != nil {
        lexpr := this.VisitUnary(ctx.Unary().(*parser.UnaryContext))
        line := 0
        if ctx.NOT() != nil {
            line = ctx.NOT().GetSymbol().GetLine()
            return newNotExpr(lexpr, ctx.NOT().GetText(), line)
        } else if ctx.SUMOP() != nil {
            line = ctx.SUMOP().GetSymbol().GetLine()
            return newSignExpr(lexpr, ctx.SUMOP().GetText(), line)
        }
    }
    return this.VisitTerm(ctx.Term().(*parser.TermContext))
}

func (this *sessionsVisitor) VisitTerm(ctx *parser.TermContext) expression {
    if ctx.Literal() != nil {
        return this.VisitLiteral(ctx.Literal().(*parser.LiteralContext))
    } else if ctx.Variable() != nil {
        id, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
        return newVariableExpr(id, nil, line)
    } else if ctx.ListExpr() != nil {
        return this.VisitListExpr(ctx.ListExpr().(*parser.ListExprContext))
    } else if ctx.LSQ() != nil {
        term := this.VisitTerm(ctx.Term().(*parser.TermContext))
        expr := this.VisitLogicalOr(ctx.LogicalOr().(*parser.LogicalOrContext))
        line := ctx.LSQ().GetSymbol().GetLine()
        return newListAccessExpr(term, expr, line)
    } else if ctx.ListSlice() != nil {
        term := this.VisitTerm(ctx.Term().(*parser.TermContext))
        // slice := this.VisitListSlice(ctx.ListSlice().(*parser.ListSliceContext))
        lexpr, rexpr := this.VisitListSlice(ctx.ListSlice().(*parser.ListSliceContext))
        line := ctx.Term().GetStart().GetLine()
        return newListSliceExpr(term, lexpr, rexpr, line)
    } else if ctx.ConditionalExpr() != nil {
        return this.VisitConditionalExpr(ctx.ConditionalExpr().(*parser.ConditionalExprContext))
    } else if ctx.Record() != nil {
        return this.VisitRecord(ctx.Record().(*parser.RecordContext))
    } else if ctx.Broker() != nil {
        return this.VisitBroker(ctx.Broker().(*parser.BrokerContext))
    } else if ctx.DOT() != nil{
        term := this.VisitTerm(ctx.Term().(*parser.TermContext))
        id := ctx.ID().GetText()
        line := ctx.DOT().GetSymbol().GetLine()
        return newRecordAccessExpr(term, id, line)
    } else if ctx.LPAR() != nil {
        return this.VisitLogicalOr(ctx.LogicalOr().(*parser.LogicalOrContext))
    } else if ctx.Abstraction() != nil {
        return this.VisitAbstraction(ctx.Abstraction().(*parser.AbstractionContext))
    }
	return nil
}

func (this *sessionsVisitor) VisitLiteral(ctx *parser.LiteralContext) expression {
    if ctx.INT_LIT() != nil {
        return newIntExpr(ctx.INT_LIT().GetText(), ctx.INT_LIT().GetSymbol().GetLine())
    } else if ctx.FLOAT_LIT() != nil {
        return newFloatExpr(ctx.FLOAT_LIT().GetText(), ctx.FLOAT_LIT().GetSymbol().GetLine())
    } else if ctx.TRUE() != nil {
        return newTrueExpr(ctx.TRUE().GetSymbol().GetLine())
    } else if ctx.FALSE() != nil {
        return newFalseExpr(ctx.FALSE().GetSymbol().GetLine())
    } else if ctx.STRING_LIT() != nil /*&& ctx.NOTHING() == nil*/ {
        str := ctx.STRING_LIT().GetText()
        str = str[1:len(str)-1]
        return newStringExpr(str, ctx.STRING_LIT().GetSymbol().GetLine())
    }/* else if ctx.NOTHING() != nil {
        str := ""
        if ctx.STRING_LIT() != nil {
            str = ctx.STRING_LIT().GetText()
            str = str[1:len(str)-1]
        }
        return newNothing(str, ctx.NOTHING().GetSymbol().GetLine())
    }*/
    return nil
}

func (this *sessionsVisitor) VisitLabelExpression(ctx *parser.LabelExpressionContext) *labelExpr {
    //var expr expression = nil
    //if ctx.Expression() != nil {
    //    expr = this.VisitExpression(ctx.Expression().(*parser.ExpressionContext))
    //}
    return newLabelExpr(ctx.LABEL().GetText(), /*expr, */ctx.LABEL().GetSymbol().GetLine())
}

func (this *sessionsVisitor) VisitListExpr(ctx *parser.ListExprContext) expression {
    expressions := make([]expression, len(ctx.AllLogicalOr()))
    for i, lor := range ctx.AllLogicalOr() {
        expressions[i] = this.VisitLogicalOr(lor.(*parser.LogicalOrContext))
    }
    return newListExpr(expressions, ctx.LSQ().GetSymbol().GetLine())
}

func (this *sessionsVisitor) VisitListSlice(ctx *parser.ListSliceContext) (expression, expression) {
    var lexpr expression = nil
    var rexpr expression = nil
    if ctx.GetLeft() != nil {
        lexpr = this.VisitLogicalOr(ctx.GetLeft().(*parser.LogicalOrContext))
    }
    if ctx.GetRight() != nil {
        rexpr = this.VisitLogicalOr(ctx.GetRight().(*parser.LogicalOrContext))
    }
    // line := ctx.COLON().GetSymbol().GetLine()
    // return newSliceExpr(lexpr, rexpr, line)
    return lexpr, rexpr
}

func (this *sessionsVisitor) VisitConditionalExpr(ctx *parser.ConditionalExprContext) expression {
    cond := this.VisitLogicalOr(ctx.LogicalOr(0).(*parser.LogicalOrContext))
    thenExpr := this.VisitLogicalOr(ctx.LogicalOr(1).(*parser.LogicalOrContext))
    elseExpr := this.VisitLogicalOr(ctx.LogicalOr(2).(*parser.LogicalOrContext))
    line := ctx.IF().GetSymbol().GetLine()
    return newConditionalExpr(cond, thenExpr, elseExpr, line)
}

func (this *sessionsVisitor) VisitRecord(ctx *parser.RecordContext) expression {
    labels := make([]string, len(ctx.AllID()))
    exprs := make([]expression, len(ctx.AllLogicalOr()))
    for i := range ctx.AllID() {
        labels[i] = ctx.ID(i).GetText()
        exprs[i] = this.VisitLogicalOr(ctx.LogicalOr(i).(*parser.LogicalOrContext))
    }
    line := ctx.LBRA().GetSymbol().GetLine()
    return newRecordExpr(labels, exprs, line)
}

func (this sessionsVisitor) VisitBroker(ctx *parser.BrokerContext) broker {
    gconfig := this.VisitConfigurationDef(ctx.ConfigurationDef().(*parser.ConfigurationDefContext))
    if ctx.REQUESTER() != nil {
        requester := this.VisitParticipant(ctx.Participant().(*parser.ParticipantContext))
        line := ctx.REQUESTER().GetSymbol().GetLine()
        return newReqBroker(requester, gconfig, line)
    }
    if ctx.SYNC() != nil {
        line := ctx.SYNC().GetSymbol().GetLine()
        return newRecBroker(gconfig, line)
    } else if ctx.ASYNC() != nil {
        line := ctx.ASYNC().GetSymbol().GetLine()
        return newFastRecBroker(gconfig, line)
    }
    return nil
}

func (this *sessionsVisitor) VisitPrimaryTerm(ctx *parser.PrimaryTermContext) expression {
    if ctx.Literal() != nil {
        return this.VisitLiteral(ctx.Literal().(*parser.LiteralContext))
    } else if ctx.PrimaryList() != nil {
        return this.VisitPrimaryList(ctx.PrimaryList().(*parser.PrimaryListContext))
    } else if ctx.PrimaryRecord() != nil {
        return this.VisitPrimaryRecord(ctx.PrimaryRecord().(*parser.PrimaryRecordContext))
    } else if ctx.Broker() != nil {
        return this.VisitBroker(ctx.Broker().(*parser.BrokerContext))
    }
    return nil
}

func (this *sessionsVisitor) VisitPrimaryList(ctx *parser.PrimaryListContext) expression {
    expressions := make([]expression, len(ctx.AllPrimaryTerm()))
    for i, pt := range ctx.AllPrimaryTerm() {
        expressions[i] = this.VisitPrimaryTerm(pt.(*parser.PrimaryTermContext))
    }
    return newListExpr(expressions, ctx.LSQ().GetSymbol().GetLine())
}

func (this *sessionsVisitor) VisitPrimaryRecord(ctx *parser.PrimaryRecordContext) expression {
    labels := make([]string, len(ctx.AllID()))
    exprs := make([]expression, len(ctx.AllPrimaryTerm()))
    for i := range ctx.AllID() {
        labels[i] = ctx.ID(i).GetText()
        exprs[i] = this.VisitPrimaryTerm(ctx.PrimaryTerm(i).(*parser.PrimaryTermContext))
    }
    line := ctx.LBRA().GetSymbol().GetLine()
    return newRecordExpr(labels, exprs, line)
}

// func (this *sessionsVisitor) VisitIoConfiguration(ctx *parser.IoConfigurationContext) expression {
//     if ctx.Variable() != nil {
//         id, line := this.VisitVariable(ctx.Variable().(*parser.VariableContext))
//         return newVariableExpr(id, nil, line)
//     } else if ctx.PrimaryRecord() != nil {
//         return this.VisitPrimaryRecord(ctx.PrimaryRecord().(*parser.PrimaryRecordContext))
//     }
//     return nil
// }

/*func (this *sessionsVisitor) VisitList_concat(ctx *parser.List_concatContext) expression {
  return newListConcatExpr(
        this.VisitExpression(ctx.Expression(0).(*parser.ExpressionContext)),
        this.VisitExpression(ctx.Expression(1).(*parser.ExpressionContext)),
        ctx.COLON().GetSymbol().GetLine())
}*/

// Session types

func (this *sessionsVisitor) VisitLocalAbstraction(ctx *parser.LocalAbstractionContext) *localAbstraction {
    args := make([]*participantType, len(ctx.AllID()))
    for i, arg := range ctx.AllID() {
        args[i] = newParticipantType(arg.GetText(), arg.GetSymbol().GetLine())
    }
    loc := this.VisitLocal(ctx.Local().(*parser.LocalContext))
    line := ctx.LOCAL().GetSymbol().GetLine()
    return newLocalAbstraction(args, loc, line)
}

func (this *sessionsVisitor) VisitType(ctx *parser.TypeContext) typedef {
    if ctx.PrimitiveType() != nil {
        return this.VisitPrimitiveType(ctx.PrimitiveType().(*parser.PrimitiveTypeContext))
    } else if ctx.LSQ() != nil {
        tdef := this.VisitType(ctx.Type_().(*parser.TypeContext))
        line := ctx.LSQ().GetSymbol().GetLine()
        return newListType(tdef, line)
    } else if ctx.RecordType() != nil {
        return this.VisitRecordType(ctx.RecordType().(*parser.RecordTypeContext))
    } else if ctx.BrokerType() != nil {
        return this.VisitBrokerType(ctx.BrokerType().(*parser.BrokerTypeContext))
    } else if ctx.Session() != nil {
        return this.VisitSession(ctx.Session().(*parser.SessionContext))
    } else if ctx.NameType() != nil {
        return this.VisitNameType(ctx.NameType().(*parser.NameTypeContext))
    } else if ctx.IoType() != nil {
        return this.VisitIoType(ctx.IoType().(*parser.IoTypeContext))
    } else if ctx.Type_() != nil {
        return this.VisitType(ctx.Type_().(*parser.TypeContext))
    }
    return nil
}

func (this *sessionsVisitor) VisitPrimitiveType(ctx *parser.PrimitiveTypeContext) typedef {
    if ctx.BOOLEAN() != nil {
        line := ctx.BOOLEAN().GetSymbol().GetLine()
        return newBoolType(line)
    } else if ctx.INT() != nil {
        line := ctx.INT().GetSymbol().GetLine()
        return newIntType(line)
    } else if ctx.FLOAT() != nil {
        line := ctx.FLOAT().GetSymbol().GetLine()
        return newFloatType(line)
    } else if ctx.STRING() != nil {
        line := ctx.STRING().GetSymbol().GetLine()
        return newStringType(line)
    }
    return nil
}

func (this *sessionsVisitor) VisitNameType(ctx *parser.NameTypeContext) *nameType {
    name, line := this.VisitName(ctx.Name().(*parser.NameContext))
    return newNameType(name, line)
}

func (this *sessionsVisitor) VisitIoType(ctx *parser.IoTypeContext) *ioType {
    line := ctx.IO().GetSymbol().GetLine()
    return newioType(line)
} 

func (this *sessionsVisitor) VisitRecordType(ctx *parser.RecordTypeContext) *recordType {
    labels := make([]string, len(ctx.AllID()))
    tdefs := make([]typedef, len(ctx.AllType_()))
    for i := range ctx.AllID() {
        labels[i] = ctx.ID(i).GetText()
        tdefs[i] = this.VisitType(ctx.Type_(i).(*parser.TypeContext))
    }
    line := ctx.RECORD().GetSymbol().GetLine()
    return newRecordType(labels, tdefs, line)
}

func (this *sessionsVisitor) VisitBrokerType(ctx *parser.BrokerTypeContext) *brokerType {
    gconfig := this.VisitConfigurationDef(ctx.ConfigurationDef().(*parser.ConfigurationDefContext))

    var requester *participantType
    if ctx.REQUESTER() != nil {
        id := ctx.ID().GetText()
        line := ctx.ID().GetSymbol().GetLine()
        requester = newParticipantType(id, line)
    }

    line := ctx.BROKER().GetSymbol().GetLine()
    return newBrokerType(gconfig, requester, line)
}

func (this *sessionsVisitor) VisitSession(ctx *parser.SessionContext) typedef {
    if ctx.LocalAbstraction() != nil {
        return this.VisitLocalAbstraction(ctx.LocalAbstraction().(*parser.LocalAbstractionContext))
    } else if ctx.Projection() != nil {
        return this.VisitProjection(ctx.Projection().(*parser.ProjectionContext))
    // } else if ctx.GlobalDef() != nil {
    //     return this.VisitGlobalDef(ctx.GlobalDef().(*parser.GlobalDefContext))
    }
    return nil
}

func (this *sessionsVisitor) VisitConfiguration(ctx *parser.ConfigurationContext) globalConfig {
    if ctx.GlobalDef() != nil {
        return this.VisitGlobalDef(ctx.GlobalDef().(*parser.GlobalDefContext))
    } else if ctx.LocalContext() != nil {
        return this.VisitLocalContext(ctx.LocalContext().(*parser.LocalContextContext))
    }
    return nil
}

func (this *sessionsVisitor) VisitSessionDef(ctx *parser.SessionDefContext) typedef {
    if ctx.Session() != nil {
        return this.VisitSession(ctx.Session().(*parser.SessionContext))
    } else if ctx.NameType() != nil {
        return this.VisitNameType(ctx.NameType().(*parser.NameTypeContext))
    }  
    return nil
}

func (this *sessionsVisitor) VisitConfigurationDef(ctx *parser.ConfigurationDefContext) typedef {
    if ctx.Configuration() != nil {
        return this.VisitConfiguration(ctx.Configuration().(*parser.ConfigurationContext))
    } else if ctx.NameType() != nil {
        return this.VisitNameType(ctx.NameType().(*parser.NameTypeContext))
    }
    return nil
}

func (this *sessionsVisitor) VisitLocal(ctx *parser.LocalContext) local {
    if ctx.SendLocal() != nil {
        return this.VisitSendLocal(ctx.SendLocal().(*parser.SendLocalContext))
    } else if ctx.ReceiveLocal() != nil {
        return this.VisitReceiveLocal(ctx.ReceiveLocal().(*parser.ReceiveLocalContext))
    } else if ctx.SelectLocal() != nil {
        return this.VisitSelectLocal(ctx.SelectLocal().(*parser.SelectLocalContext))
    } else if ctx.BranchLocal() != nil {
        return this.VisitBranchLocal(ctx.BranchLocal().(*parser.BranchLocalContext))
    } else if ctx.End() != nil {
        line := this.VisitEnd(ctx.End().(*parser.EndContext))
        return newEndLocal(line)
    }
    return nil
}

func (this *sessionsVisitor) VisitSendLocal(ctx *parser.SendLocalContext) *sendLocal {
    participant := ctx.ID().GetText()
    t := this.VisitType(ctx.Type_().(*parser.TypeContext))
    line := ctx.SEND().GetSymbol().GetLine()
    var loc local = nil
    if ctx.Local() != nil {
        loc = this.VisitLocal(ctx.Local().(*parser.LocalContext))
    } else {
        loc = newEndLocal(line)
    }
    ptype := newParticipantType(participant, line)
    return newSendLocal(ptype, t, loc, line)
}

func (this *sessionsVisitor) VisitReceiveLocal(ctx *parser.ReceiveLocalContext) *receiveLocal {
    participant := ctx.ID().GetText()
    t := this.VisitType(ctx.Type_().(*parser.TypeContext))
    line := ctx.RECEIVE().GetSymbol().GetLine()
    var loc local = nil
    if ctx.Local() != nil {
        loc = this.VisitLocal(ctx.Local().(*parser.LocalContext))
    } else {
        loc = newEndLocal(line)
    }
    ptype := newParticipantType(participant, line)
    return newReceiveLocal(ptype, t, loc, line)
}

func (this *sessionsVisitor) VisitSelectLocal(ctx *parser.SelectLocalContext) *selectLocal {
    participant := ctx.ID().GetText()
    line := ctx.ID().GetSymbol().GetLine()
    ptype := newParticipantType(participant, line)
    labels := make([]*labelType, len(ctx.AllLabelLocal()))
    locals := make([]local, len(ctx.AllLabelLocal()))
    for i, s := range ctx.AllLabelLocal() {
        labels[i], locals[i] = this.VisitLabelLocal(s.(*parser.LabelLocalContext))
    }
    line = ctx.SELECT().GetSymbol().GetLine()
    return newSelectLocal(ptype, labels, locals, line)
}

func (this *sessionsVisitor) VisitBranchLocal(ctx *parser.BranchLocalContext) *branchLocal {
    participant := ctx.ID().GetText()
    line := ctx.ID().GetSymbol().GetLine()
    ptype := newParticipantType(participant, line)
    labels := make([]*labelType, len(ctx.AllLabelLocal()))
    locals := make([]local, len(ctx.AllLabelLocal()))
    for i, s := range ctx.AllLabelLocal() {
        labels[i], locals[i] = this.VisitLabelLocal(s.(*parser.LabelLocalContext))
    }
    line = ctx.BRANCH().GetSymbol().GetLine()
    return newBranchLocal(ptype, labels, locals, line)
}

func (this *sessionsVisitor) VisitLabelLocal(ctx *parser.LabelLocalContext) (*labelType, local) {
    lab := this.VisitLabelType(ctx.LabelType().(*parser.LabelTypeContext))
    loc := this.VisitLocal(ctx.Local().(*parser.LocalContext))
    return lab, loc
}

func (this *sessionsVisitor) VisitLabelType(ctx *parser.LabelTypeContext) *labelType {
    label := ctx.LABEL().GetText()
    line := ctx.LABEL().GetSymbol().GetLine()
    return newLabelType(label, line)
}

func (this *sessionsVisitor) VisitEnd(ctx *parser.EndContext) int {
    return ctx.END().GetSymbol().GetLine()
}

func (this *sessionsVisitor) VisitLocalContext(ctx *parser.LocalContextContext) *localContext{
    args := make([]*participantType, len(ctx.AllID()))
    for i, arg := range ctx.AllID() {
        args[i] = newParticipantType(arg.GetText(), arg.GetSymbol().GetLine())
    }
    locals := make([]local, len(ctx.AllLocal()))
    for i, loc := range ctx.AllLocal() {
        locals[i] = this.VisitLocal(loc.(*parser.LocalContext))
    }

    line := ctx.CONTEXT().GetSymbol().GetLine()
    return newLocalContext(args, locals, line)
}

func (this *sessionsVisitor) VisitProjection(ctx *parser.ProjectionContext) *projection {
    conf := this.VisitConfigurationDef(ctx.ConfigurationDef().(*parser.ConfigurationDefContext))
    //name, _ := this.VisitName(ctx.Name().(*parser.NameContext))
    id := ctx.ID().GetText()
    line := ctx.ID().GetSymbol().GetLine()
    ptype := newParticipantType(id, line)
    line = ctx.LPAR().GetSymbol().GetLine()
    return newProjection(conf, ptype, line)
}

func (this *sessionsVisitor) VisitGlobalDef(ctx *parser.GlobalDefContext) *globalDef {
    args := make([]*participantType, len(ctx.AllID()))
    for i, arg := range ctx.AllID() {
        args[i] = newParticipantType(arg.GetText(), arg.GetSymbol().GetLine())
    }
    glob := this.VisitGlobal(ctx.Global().(*parser.GlobalContext))
    line := ctx.GLOBAL().GetSymbol().GetLine()
    return newGlobalDef(args, glob, line)
}

func (this *sessionsVisitor) VisitGlobal(ctx *parser.GlobalContext) global {
    if ctx.Pass() != nil {
        return this.VisitPass(ctx.Pass().(*parser.PassContext))
    } else if ctx.Choice() != nil {
      return this.VisitChoice(ctx.Choice().(*parser.ChoiceContext))
    } else if ctx.End() != nil {
        line := this.VisitEnd(ctx.End().(*parser.EndContext))
        return newEndGlobal(line)
    }
    return nil
}

func (this *sessionsVisitor) VisitPass(ctx *parser.PassContext) global {
    senderID := ctx.ID(0).GetText()
    receiverID := ctx.ID(1).GetText()
    senderLine := ctx.ID(0).GetSymbol().GetLine()
    receiverLine := ctx.ID(1).GetSymbol().GetLine()
    tdef := this.VisitType(ctx.Type_().(*parser.TypeContext))
    line := ctx.PASS().GetSymbol().GetLine()
    var glob global = nil
    if ctx.DOT() != nil {
        glob = this.VisitGlobal(ctx.Global().(*parser.GlobalContext))
    } else {
        glob = newEndGlobal(line)
    }
    sender := newParticipantType(senderID, senderLine)
    receiver := newParticipantType(receiverID, receiverLine)
    return newPassGlobal(sender, receiver, tdef, glob, line)
}

func (this *sessionsVisitor) VisitChoice(ctx *parser.ChoiceContext) global {
    senderID := ctx.ID(0).GetText()
    receiverID := ctx.ID(1).GetText()
    senderLine := ctx.ID(0).GetSymbol().GetLine()
    receiverLine := ctx.ID(1).GetSymbol().GetLine()
    line := ctx.PASS().GetSymbol().GetLine()
    labels := make([]*labelType, len(ctx.AllLabelGlobal()))
    globals := make([]global, len(ctx.AllLabelGlobal()))
    for i, s := range ctx.AllLabelGlobal() {
        labels[i], globals[i] = this.VisitLabelGlobal(s.(*parser.LabelGlobalContext))
    }
    sender := newParticipantType(senderID, senderLine)
    receiver := newParticipantType(receiverID, receiverLine)
    return newChoiceGlobal(sender, receiver, labels, globals, line)
}

func (this *sessionsVisitor) VisitLabelGlobal(ctx *parser.LabelGlobalContext) (*labelType, global) {
    lab := this.VisitLabelType(ctx.LabelType().(*parser.LabelTypeContext))
    glob := this.VisitGlobal(ctx.Global().(*parser.GlobalContext))
    return lab, glob
}
