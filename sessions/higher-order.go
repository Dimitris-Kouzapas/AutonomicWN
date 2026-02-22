package ast

import (
    "sessions/util"
)

/******************************************************************************
 * abstraction
 ******************************************************************************/
type abstraction struct {
    baseNode
    body process
    parameters []*participantExpr
}

func newAbstraction(body process, parameters []*participantExpr, line int) *abstraction {
    return &abstraction {
        baseNode: baseNode{line: line},
        body: body,
        parameters: parameters,
    }
}

func (a *abstraction) setFilename(filename string) {
    for _, parameter := range a.parameters {
        parameter.setFilename(filename)
    }
    a.body.setFilename(filename)
    a.baseNode.setFilename(filename)
}

func (a *abstraction) getType() typedef {
    // cannot decide on the type during the typecheck phase
    return nil
}

func (a *abstraction) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    for _, parameter := range a.parameters {
          parameter.typeCheck(ctx, log)
    }
    a.body.typeCheck(ctx, log)
}

func (a *abstraction) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    for _, parameter := range a.parameters {
          parameter.projectionCheck(ctx, elog, rlog)
    }
    a.body.projectionCheck(ctx, elog, rlog) 
}

func (a *abstraction) expressionCheck(ctx *expressionCheckContext, log util.ErrorLog) typedef {
    ctx.pushFrame()
    defer ctx.popFrame()

    for _, parameter := range a.parameters {
        if !ctx.addParticipant(parameter) {
            parameter.reportErrorf(log, "duplicate definition of parameter: %q.", parameter.String())
        }
    }

    parameters := make([]*participantType, len(a.parameters))
    for i, p := range a.parameters {
        td := p.expressionCheck(ctx, log)
        pt, ok := td.(*participantType)
        if !ok || pt == nil {
            p.reportErrorf(log, "expecting participant type; instead found: %q.", td.String())
            continue
        }
        parameters[i] = pt
    }

    lin := newLinearContext(a.parameters, a.line)
    lin = a.body.expressionCheck(ctx, lin, log)
    loc := lin.removeSession(a.parameters, a.line)
    return newLocalAbstraction(parameters, loc, a.line)
}

func (a *abstraction) sessionCheck(ctx *sessionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    a.body.sessionCheck(ctx, elog, rlog)
}

func (a *abstraction) runtimeAbstraction(ctx *evaluationContext) *runtimeAbstraction {
    return &runtimeAbstraction {
        abstraction: abstraction{
            baseNode: baseNode{line: a.line, filename: a.filename},
            body: a.body,
            parameters: a.parameters,
        },
        variables: ctx.cloneVariables(),
    }
}

func (a *abstraction) evaluate(ctx *evaluationContext) expression {
    return a.runtimeAbstraction(ctx)
    //return a
}

func (_ *abstraction) operation(_ *evaluationContext, _ expression, _ string) expression { return nil }

func (a *abstraction) prettyPrint(iw util.IndentedWriter) {
    iw.Print("proc")
    for _, parameter := range a.parameters {
        iw.Print(" ")
        parameter.prettyPrint(iw)
    }
    iw.Println(".")
    iw.Inc()
    a.body.prettyPrint(iw)
    iw.Dec()
}

func (a *abstraction) String() string {
    s := "proc"
    for _, parameter := range a.parameters {
        s += " " + parameter.String()
    }
    return s + ". " + a.body.String()
}

func (a *abstraction) goCode(iw util.IndentedWriter) {}

/******************************************************************************
 * runtime abstraction
  ******************************************************************************/

type runtimeAbstraction struct {
    abstraction
    variables map[string]expression
}

func (ra *runtimeAbstraction) spawn(sendChannels []util.Channel[expression], receiveChannels []util.Channel[expression], ctx *evaluationContext) {
    defer ctx.done()
    //TODO channel clean up should not be done whenever we have true recursion
    // defer ctx.cleanup() // <- closes channels and releases resources
    for i := range ra.parameters {
        ctx.addParticipantChannel(ra.parameters[i], sendChannels[i], receiveChannels[i])
    }
    ctx.variables = ra.variables
    ra.body.evaluate(ctx)
}

func (ra *runtimeAbstraction) runtimeAbstraction() *runtimeAbstraction {
    variables := make(map[string]expression)
    for k, v := range ra.variables {
       variables[k] = v
    }
    return &runtimeAbstraction {
        abstraction: abstraction{
            baseNode: baseNode{line: ra.line, filename: ra.filename},
            body: ra.body,
            parameters: ra.parameters,
        },
        variables: variables,
    }
}

func (ra *runtimeAbstraction) evaluate(_ *evaluationContext) expression {
    return ra.runtimeAbstraction()
}

/******************************************************************************
 * name abstraction
 ******************************************************************************/

// type nameAbstraction struct {
//     baseNode
//     variable *variableExpr
// }

// func newNameAbstraction(name string, line int) *nameAbstraction {
//     return &nameAbstraction {
//         baseNode: baseNode{line: line},
//         variable: newVariableExpr(name, nil, line),
//     }
// }

// func (na *nameAbstraction) setFilename(filename string) {
//     na.baseNode.setFilename(filename)
//     na.variable.setFilename(filename)
// }

// func (na *nameAbstraction) getType() typedef { return na.variable.getType() }

// func (na *nameAbstraction) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
//     na.variable.typeCheck(ctx, log)
// }

// func (na *nameAbstraction) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
//     na.variable.projectionCheck(ctx, elog, rlog) 
// }

// func (na *nameAbstraction) expressionCheck(ctx *expressionCheckContext, log util.ErrorLog) typedef {
//     tdef := na.variable.expressionCheck(ctx, log)
//     if tdef != nil {
//         if r, ok := tdef.getType().(*localAbstraction); ok && r != nil {
//             return r
//         }
//         stream := util.NewStream().Inc().Inc()
//         tdef.prettyPrint(stream)
//         na.reportErrorf(log, "expecting local abstraction type; instead found:\n%v", stream.String())
//     }
//     return nil
// }

// func (_ *nameAbstraction) sessionCheck(_ *sessionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}
// func (na *nameAbstraction) evaluate(ctx *evaluationContext) expression {
//     //expressionChecking ensures that na.variable.evaluate(ctx) will return abstraction
//     abstr := na.variable.evaluate(ctx)
//     return abstr.evaluate(ctx)
// }

// func (_ *nameAbstraction) operation(_ *evaluationContext, _ expression, _ string) expression  { return nil }
// func (na *nameAbstraction) prettyPrint(iw util.IndentedWriter)          { iw.Print(na.String()) }
// func (na *nameAbstraction) String() string                              { return na.variable.String() }
// func (_ *nameAbstraction) goCode(_ util.IndentedWriter)                 {}

/******************************************************************************
 * application
 ******************************************************************************/

type application struct {
    baseNode
    adef expression
    arguments []*participantExpr
}

func newApplication(adef expression, arguments []*participantExpr, line int) *application {
    return &application {
        baseNode: baseNode{line: line},
        adef: adef,
        arguments: arguments,
    }
}

func (a *application) setFilename(filename string) {
    a.baseNode.setFilename(filename)
    a.adef.setFilename(filename)
    for _, argument := range a.arguments {
        argument.setFilename(filename)
    }
}

func (a *application) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    a.adef.typeCheck(ctx, log)
    for _, argument := range a.arguments {
        argument.typeCheck(ctx, log)
    }
}

func (a *application) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    a.adef.projectionCheck(ctx, elog, rlog)
}

func (a *application) expressionCheck(pexpr *participantExpr, pset *util.HashSet[string], ctx *expressionCheckContext, log util.ErrorLog) local {
    tdef := a.adef.expressionCheck(ctx, log)

    seen := util.NewHashSet[string]()
    for _, argument := range a.arguments {
        if seen.Contains(argument.id) {
            argument.reportErrorf(log, "cannot use participant %q twice as argument.", argument.String())
        } else {
            seen.Add(argument.id)
        }
        if argument.id == pexpr.id {
            a.reportErrorf(log, "self-reference participant argument: %q", argument.String())
        }
        if !pset.Contains(argument.id) {
            a.reportErrorf(log, "unknown participant: %q", argument.String())
        }
    }

    if tdef == nil {
        return nil
    }
    r, ok := tdef.(*localAbstraction)
    if !ok || r == nil {
        a.reportErrorf(log, "expecting local abstraction type; instead found: %q.", tdef)
        return newEndLocal(a.line)
    }

    if len(r.parameters) != len(a.arguments) {
        a.reportErrorf(
            log,
            "wrong number of arguments; expecting %d, but found %d.", 
            len(r.parameters),
            len(a.arguments),
        )
    }
    return r.apply(a.arguments)
}

func (a *application) sessionCheck(ctx *sessionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    a.adef.sessionCheck(ctx, elog, rlog)
}

func (a *application) evaluate(pexpr *participantExpr, channels channelMap, ctx *evaluationContext) {
    // sendChannels := make([]chan<- expression, len(a.arguments))
    // receiveChannels := make([]<-chan expression, len(a.arguments))
    sendChannels := make([]util.Channel[expression], len(a.arguments))
    receiveChannels := make([]util.Channel[expression], len(a.arguments))
    for i, pexpr2 := range a.arguments {
        sendChannels[i] = channels[pexpr.id][pexpr2.id]
        receiveChannels[i] = channels[pexpr2.id][pexpr.id]
    }
    val := a.adef.evaluate(ctx)
    var ra *runtimeAbstraction
    switch v := val.(type) {
        case *abstraction:
            ra = v.runtimeAbstraction(ctx)
        case *runtimeAbstraction:
            ra = v
        default:
            return
    }
    // ra, ok := val.(*runtimeAbstraction)
    // if !ok || ra == nil {
    //     return
    // }
    newContext := ctx.emptyCopy()
    ctx.add()
    go ra.spawn(sendChannels, receiveChannels, newContext)
}

func (a *application) prettyPrint(iw util.IndentedWriter) {
    a.adef.prettyPrint(iw)
    for _, argument := range a.arguments {
        iw.Print(" ")
        argument.prettyPrint(iw)
    }
}

func (a *application) String() string {
    s := a.adef.String()
    for _, argument := range a.arguments {
        s += " " + argument.String()
    }
    return s
}
