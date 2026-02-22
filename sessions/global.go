package ast

import (
    "fmt"
    "sessions/util"
)

type globalConfig interface {
    typedef
    channels(bool) (channelMap)
    project(*participantType) *localAbstraction
    parameterTypes() []*participantType
    participant(int) string
    participants() []string
}

type globalConfigImpl struct {
    baseNode
    parameters []*participantType
}

func (gc *globalConfigImpl) setFilename(filename string) {
    gc.baseNode.setFilename(filename)
    for _, p := range gc.parameters {
        p.setFilename(filename)
    }
}

func (gc *globalConfigImpl) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    for _, p := range gc.parameters {
        p.typeCheck(ctx, log)
    }
}

func (gc *globalConfigImpl) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, _ util.ReportLog) {
    for _, p := range gc.parameters {
        if !ctx.addParticipant(p) {
            p.reportErrorf(elog, "duplicated participant definition: %q.", p.String())
        }
    }
}

func (gc *globalConfigImpl) channels(active bool) (channels channelMap) {
    parts := gc.participants()

    channels = make(channelMap, len(parts))
    for _, p := range parts {
        //channels[p] = make(map[string]chan expression, len(parts)-1)
        channels[p] = make(map[string]util.Channel[expression], len(parts)-1)
    }

    for i, p1 := range parts {
        for j, p2 := range parts {
            if i == j { continue }
            channels[p1][p2] = util.NewBasicChannel[expression](active) //make(chan expression)
        }
    }
    return
}

func (gc *globalConfigImpl) parameterTypes() []*participantType {
    return gc.parameters
}

func (gc *globalConfigImpl) participant(i int) string {
    return gc.parameters[i].participant
}

func (gc *globalConfigImpl) participants() []string {
    parts := make([]string, len(gc.parameters))
    for i, p := range gc.parameters {
        parts[i] = p.participant
    }
    return parts
}

func (gc *globalConfigImpl) substitute(ptypes []*participantType) (map[string]string, []*participantType) {
    substitution := make(map[string]string, len(gc.parameters))
    parameters := make([]*participantType, len(gc.parameters))

    for i := range gc.parameters {
        participant := gc.participant(i)
        if i < len(ptypes) {
            substitution[participant] = ptypes[i].participant
            parameters[i] = ptypes[i]
        } else {
            substitution[participant] = participant
            parameters[i] = gc.parameters[i]
        }
    }
    return substitution, parameters
}

func (gc *globalConfigImpl) defaultValue() []*participantExpr {
    parameters := make([]*participantExpr, len(gc.parameters))
    for i, par := range gc.parameters {
        parameters[i] = par.defaultValue().(*participantExpr)
    }
    return parameters
}

/******************************************************************************
 * local ctx
 ******************************************************************************/

type localContext struct {
    globalConfigImpl
    locals []local
    localMap map[string]local
}

func newLocalContext(parameters []*participantType, locals []local, line int) *localContext {
    lc := &localContext {
        globalConfigImpl: globalConfigImpl {
            baseNode: baseNode{line: line},
            parameters: parameters,
        },
        locals: locals,
        localMap: make(map[string]local, len(locals)),
    }
    for i, loc := range locals {
        lc.localMap[lc.participant(i)] = loc
    }
    return lc
}

func (lc *localContext) setFilename(filename string) {
    lc.globalConfigImpl.setFilename(filename)
    for _, loc := range lc.locals {
        loc.setFilename(filename)
    }
}

func (lc *localContext) getType() typedef { return lc }

func (lc *localContext) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    lc.globalConfigImpl.typeCheck(ctx, log)
    for _, loc := range lc.locals {
        loc.typeCheck(ctx, log)
    }
}

func (lc *localContext) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    ctx.push()
    defer ctx.pop()
    lc.globalConfigImpl.projectionCheck(ctx, elog, rlog)

    if glob, ok := lc.liveness(); !ok {
        stream := util.NewStream().Inc().Inc()
        lc.prettyPrint(stream)
        lc.reportErrorf(elog, "non-live session roles\n%v", stream.String())
    } else if ctx.analysis() {
        gdef := newGlobalDef(lc.parameters, glob, lc.line)
        stream := util.NewStream()
        stream.Printf("local context liveness analysis at line %v produces global type:\n", lc.line)
        stream.Inc().Inc()
        gdef.prettyPrint(stream)
        lc.reportf(rlog, stream.String())
    }
}

func (lc *localContext) subtypeOf(tdef typedef) bool {
    visited := util.NewHashSet[typePair]()
    return lc.subtypeOf_(tdef, visited)
}

func (lc *localContext) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
    tPair := typePair{a: lc, b:tdef}

    if visited.Contains(tPair) {
        return true
    }

    visited.Add(tPair)

    switch gconfig := tdef.getType().(type) {
        case *nothingType:
            return true
        case *localContext:
            if (len(lc.parameters) != len(gconfig.parameters)) {
                return false
            }
            for i, loc1 := range lc.locals {
                if !lc.parameters[i].subtypeOf(gconfig.parameters[i]) {
                    return false
                }
                loc2 := gconfig.locals[i]
                if !loc1.subtypeOf(loc2, visited) {
                    return false
                }
            }
            return true
        case *globalDef:
            lcont := gconfig.projectionSet()
            return lc.subtypeOf_(lcont, visited)
        default:
            return false
    }
}

func (lc *localContext) join(tdef typedef) (typedef, bool) {
    switch gconfig := tdef.getType().(type) {
        case *nothingType:
            return lc, true
        case *localContext:
            if len(lc.parameters) != len(gconfig.parameters) {
                return lc, false
            }
            locals := make([]local, len(lc.locals))
            result := true
            // Do not sustitute names here.
            for i, loc1 := range lc.locals {
                if !lc.parameters[i].subtypeOf(gconfig.parameters[i]) {
                    result = false
                    locals[i] = loc1
                    continue
                }
                loc2:= gconfig.locals[i]
                ok := true
                locals[i], ok = loc1.join(loc2)
                if !ok {
                    result = false
                    locals[i] = loc1
                }
            }
            return newLocalContext(lc.parameters, locals, lc.line), result
        case *globalDef:
            lcont := gconfig.projectionSet()
            return lc.join(lcont)
        default:
            return lc, false
    }
}

func (lc *localContext) substitute(ptypes []*participantType) *localContext {
    substitution, parameters := lc.globalConfigImpl.substitute(ptypes)
    locals := make([]local, len(lc.locals))
    for i := range lc.locals {
        locals[i] = lc.locals[i].substitute(substitution)
    }
    return newLocalContext(parameters, locals, lc.line)
}

func (lc *localContext) project(ptype *participantType) *localAbstraction {
    var loc local = newEndLocal(lc.line)
    for i, pt := range lc.parameters {
        if ptype.subtypeOf(pt) {
            loc = lc.locals[i]
            break
        }
    }
    parameters := make([]*participantType, 0, len(lc.parameters))
    for _, p := range lc.parameters {
        if loc.hasParticipant(p) {
            parameters = append(parameters, p)
        }
    }
    return newLocalAbstraction(parameters, loc, lc.line)
}


func (lc *localContext) liveness() (global, bool) {
    return liveness(lc.localMap)
}

func (lc *localContext) defaultValue() expression {
    return newNothingf(lc.line, "unexpected local context type")
}

func (lc *localContext) prettyPrint(iw util.IndentedWriter) {
    iw.Println("context {")
    iw.Inc()
    for i := range lc.parameters {
        lc.parameters[i].prettyPrint(iw)
        iw.Print(":\t")
        iw.Inc()
        lc.locals[i].prettyPrint(iw)
        iw.Dec()
        iw.Println()
    }
    iw.Dec()
    iw.Println("}")
    return
}

func (lc *localContext) String() string {
    s := "context {"
    for i := range lc.parameters {
        s += " " + lc.parameters[i].String() + ": " + lc.locals[i].String()
    }
    return s + " }"
}

/******************************************************************************
 * global def
 ******************************************************************************/

type globalDef struct {
    *localContext
    glob global
}

func newGlobalDef(parameters []*participantType, glob global, line int) *globalDef {
    lc := &localContext{
        globalConfigImpl: globalConfigImpl{
            baseNode:   baseNode{line: line},
            parameters: parameters,
        },
    }
    return &globalDef{
        localContext: lc,
        glob:         glob,
    }
}

func (gd *globalDef) setFilename(filename string) {
    gd.localContext.setFilename(filename)
    gd.glob.setFilename(filename)
}

func (gd *globalDef) getType() typedef { return gd }

func (gd *globalDef) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    gd.globalConfigImpl.typeCheck(ctx, log)
    gd.glob.typeCheck(ctx, log)
}

func (gd *globalDef) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    ctx.push()
    defer ctx.pop()
    gd.globalConfigImpl.projectionCheck(ctx, elog, rlog)
    gd.localMap = gd.glob.projectionCheck(ctx, elog, rlog)
    gd.locals = make([]local, len(gd.parameters))
    for i := range gd.parameters {
        gd.locals[i] = gd.localMap[gd.participant(i)]
    }
}

func (gd *globalDef) subtypeOf(tdef typedef) bool {
    visited := util.NewHashSet[typePair]()
    return gd.subtypeOf_(tdef, visited)
}

func (gd *globalDef) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
    tPair := typePair{a: gd, b:tdef}

    if visited.Contains(tPair) {
        return true
    }

    visited.Add(tPair)

    switch gconfig := tdef.getType().(type) {
        case *nothingType:
            return true
        case *globalDef, *localContext:
            lcont := gd.projectionSet()
            return lcont.subtypeOf_(gconfig, visited)
        default:
            return false
    }
}

func (gd *globalDef) join(tdef typedef) (typedef, bool) {
    switch gconfig := tdef.getType().(type) {
        case *nothingType:
            lcont := gd.projectionSet()
            return lcont, true
        case *globalDef, *localContext:
            lcont := gd.projectionSet()
            return lcont.join(gconfig)
        default:
            return gd, false
    }
}

func (gd *globalDef) projectionSet() *localContext {
    return gd.localContext
}

func (gd *globalDef) substitute(ptypes []*participantType) *globalDef {
    substitution, parameters := gd.globalConfigImpl.substitute(ptypes)
    glob := gd.glob.substitute(substitution)
    return newGlobalDef(parameters, glob, gd.line)
}

func (gd *globalDef) defaultValue() expression {
    return newNothingf(gd.line, "unexpected global type")
}

func (gd *globalDef) prettyPrint(iw util.IndentedWriter) {
    iw.Print("global")
    for _, participant := range gd.parameters {
        iw.Print(" ")
        participant.prettyPrint(iw)
    }
    iw.Println(". ")
    iw.Inc()
    gd.glob.prettyPrint(iw)
    iw.Dec()
}

func (gd *globalDef) String() string {
    s := "global"
    for _, participant := range gd.parameters {
        s += " " + participant.String()
    }
    return s + ". " + gd.glob.String()
}

/******************************************************************************
 * global
 ******************************************************************************/
type global interface {
    setFilename(string)
    typeCheck(*typeCheckContext, util.ErrorLog) // map[string]local
    projectionCheck(*projectionCheckContext, util.ErrorLog, util.ReportLog) map[string]local
    reportErrorf(util.ErrorLog, string, ...interface{})
    substitute(map[string]string) global
    subtypeOf(global, *util.HashSet[typePair]) bool
    join(global) (global, bool)
    prettyPrint(util.IndentedWriter)
    fmt.Stringer
}

/******************************************************************************
 * end global
 ******************************************************************************/

type endGlobal struct {
    baseNode
}

func newEndGlobal(line int) *endGlobal {
    return &endGlobal {
        baseNode: baseNode{line: line},
    }
}

func (e *endGlobal) setFilename(filename string) {
    e.baseNode.setFilename(filename)
}

func (_ *endGlobal) typeCheck(_ *typeCheckContext, _ util.ErrorLog) {}

func (e *endGlobal) projectionCheck(ctx *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) map[string]local {
    parts := ctx.participants()
    localMap := make(map[string]local, len(parts))
    for _, participant := range parts {
        localMap[participant] = newEndLocal(e.line)
    }
    return localMap
}

func (_ *endGlobal) subtypeOf(glob global, _ *util.HashSet[typePair]) bool {
    _, ok := glob.(*endGlobal)
    return ok
}

func (e *endGlobal) join(glob global) (global, bool) {
    _, ok := glob.(*endGlobal)
    return e, ok
}

func (e *endGlobal) substitute(substitution map[string]string) global {
    return newEndGlobal(e.line)
}

func (e *endGlobal) prettyPrint(iw util.IndentedWriter) { iw.Print(e.String()) }
func (_ *endGlobal) String() string { return "end" }

/******************************************************************************
 * pass global
 ******************************************************************************/

type passGlobal struct {
    baseNode
    sender *participantType
    receiver *participantType
    tdef typedef
    cont global
}

func newPassGlobal(sender *participantType, receiver *participantType, tdef typedef, cont global, line int) *passGlobal {
    return &passGlobal {
        baseNode: baseNode{line: line},
        sender: sender,
        receiver: receiver,
        tdef: tdef,
        cont: cont,
    }
}

func (p *passGlobal) setFilename(filename string) {
    p.sender.setFilename(filename)
    p.receiver.setFilename(filename)
    p.tdef.setFilename(filename)
    p.cont.setFilename(filename)
    p.baseNode.setFilename(filename)
}

func (p *passGlobal) typeCheck(ctx *typeCheckContext, log util.ErrorLog) /*map[string]local*/ {
    p.tdef.typeCheck(ctx, log)
    p.cont.typeCheck(ctx, log)
}

func (p *passGlobal) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) map[string]local {
    if !ctx.containsParticipant(p.sender) {
        p.reportErrorf(elog, "undefined participant: %q", p.sender.String())
    }
    if !ctx.containsParticipant(p.receiver) {
        p.reportErrorf(elog, "undefined participant: %q", p.receiver.String())
    }

    p.tdef.projectionCheck(ctx, elog, rlog)
    localMap := p.cont.projectionCheck(ctx, elog, rlog)
    loc := localMap[p.sender.participant]
    localMap[p.sender.participant] = newSendLocal(p.receiver, p.tdef, loc, p.line)
    loc = localMap[p.receiver.participant]
    localMap[p.receiver.participant] = newReceiveLocal(p.sender, p.tdef, loc, p.line)
    return localMap
}


func (p *passGlobal) subtypeOf(glob global, visited *util.HashSet[typePair]) bool {
    pglob, ok := glob.(*passGlobal)
    if !ok {
        return false
    }
    if !p.sender.subtypeOf(pglob.sender) {
        return false
    }
    if !p.receiver.subtypeOf(pglob.receiver) {
        return false
    }

    if !(p.tdef.subtypeOf_(pglob.tdef, visited) && pglob.tdef.subtypeOf_(p.tdef, visited)){
        return false
    }
    return p.cont.subtypeOf(pglob.cont, visited)
}

func (p *passGlobal) join(glob global) (global, bool) {
    pglob, ok := glob.(*passGlobal)
    if !ok {
        return p, false
    }
    if !p.sender.subtypeOf(pglob.sender) {
        return p, false
    }
    if !p.receiver.subtypeOf(pglob.receiver) {
        return p, false
    }

    if !(p.tdef.subtypeOf(pglob.tdef) && pglob.tdef.subtypeOf(p.tdef)) {
        return p, false
    }
    cont, ok2 := p.cont.join(pglob.cont)
    return newPassGlobal(p.sender, p.receiver, p.tdef, cont, p.line), ok2
}

func (p *passGlobal) substitute(substitution map[string]string) global {
    senderId, ok := substitution[p.sender.participant]
    if !ok {
        senderId = p.sender.participant
    }
    receiverId, ok := substitution[p.receiver.participant]
    if !ok {
        receiverId = p.receiver.participant
    }

    sender := newParticipantType(senderId, p.sender.line)
    receiver := newParticipantType(receiverId, p.receiver.line)
    glob := p.cont.substitute(substitution)

    return newPassGlobal(sender, receiver, p.tdef, glob, p.line)
}

func (p *passGlobal) prettyPrint(iw util.IndentedWriter) {
    p.sender.prettyPrint(iw)
    iw.Print(" -> ")
    p.receiver.prettyPrint(iw)
    iw.Print(": (")
    p.tdef.prettyPrint(iw)
    iw.Println("). ")
    p.cont.prettyPrint(iw)
}

func (p *passGlobal) String() string {
    s := p.sender.String() + " -> " + p.receiver.String()
    s += ": (" + p.tdef.String() + "). "
    s += p.cont.String()
    return s
}

/******************************************************************************
 * choice global
 ******************************************************************************/

type choiceGlobal struct {
    baseNode
    sender *participantType
    receiver *participantType
    labels []*labelType
    globals []global
    labelMap map[string]global
}

func newChoiceGlobal(sender *participantType, receiver *participantType, labels []*labelType, globals []global, line int) *choiceGlobal {
    labelMap := make(map[string]global, len(labels))
    for i, lab := range labels {
        labelMap[lab.label] = globals[i]
    }
    return &choiceGlobal{
        baseNode: baseNode{line: line},
        sender:   sender,
        receiver: receiver,
        labels:   labels,
        globals:  globals,
        labelMap: labelMap,
    }
}

func (cg *choiceGlobal) setFilename(filename string) {
    cg.sender.setFilename(filename)
    cg.receiver.setFilename(filename)
    for i, glob := range cg.globals {
        glob.setFilename(filename)
        cg.labels[i].setFilename(filename)
    }
    cg.baseNode.setFilename(filename)
}

func (cg *choiceGlobal) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    ctx.resetLabels()
    for _, lab := range cg.labels {
        if !ctx.addLabel(lab) {
            lab.reportErrorf(log, "duplicate definition of select label: %v.", lab.label)
        }
    }
    for _, glob := range cg.globals {
        glob.typeCheck(ctx, log)
    }
}


func (cg *choiceGlobal) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) map[string]local {
    if !ctx.containsParticipant(cg.sender) {
        cg.reportErrorf(elog, "undefined participant: %q", cg.sender.String())
    }
    if !ctx.containsParticipant(cg.receiver) {
        cg.reportErrorf(elog, "undefined participant: %q", cg.receiver.String())
    }

    localMaps := make([]map[string]local, len(cg.globals))
    for i, glob := range cg.globals {
        localMaps[i] = glob.projectionCheck(ctx, elog, rlog)
    }
    parts := ctx.participants()
    localMap := make(map[string]local, len(parts))
    for _, participant := range parts {
        if cg.sender.participant == participant {
            locals := make([]local, len(localMaps))
            for j := range localMaps {
                locals[j] = localMaps[j][participant]
            }
            localMap[participant] = newSelectLocal(cg.receiver, cg.labels, locals, cg.line)
        } else if cg.receiver.participant == participant {
            locals := make([]local, len(localMaps))
            for j := range localMaps {
                locals[j] = localMaps[j][participant]
            }
            localMap[participant] = newBranchLocal(cg.sender, cg.labels, locals, cg.line)
        } else {
            var max local = nil
            for j := range localMaps {
                ok := true
                loc := localMaps[j][participant]
                if j == 0 {
                    max, ok = loc, true
                } else {
                    max, ok = max.merge(loc)
                }
                if !ok {
                    stream1 := util.NewStream().Inc().Inc()
                    stream2 := util.NewStream().Inc().Inc()
                    stream1.Print(participant + ": ")
                    max.prettyPrint(stream1)
                    stream2.Print(participant + ": ")
                    loc.prettyPrint(stream2)
                    msg := fmt.Sprintf("choice on participant %q (label:%q); expecting type: \n", participant, cg.labels[j].String())
                    msg += fmt.Sprintf("%v\n", stream1.String())
                    msg += fmt.Sprintf("\tbut found:\n")
                    msg += fmt.Sprintf("%v", stream2.String())
                    cg.globals[j].reportErrorf(elog, msg)
                }
            }
            localMap[participant] = max
        }
    }
    return localMap
}

func (cg *choiceGlobal) subtypeOf(glob global, visited *util.HashSet[typePair]) bool {
    cglob, ok := glob.(*choiceGlobal)
    if !ok {
        return false
    }

    if !cg.sender.subtypeOf(cglob.sender) {
        return false
    }

    if !cg.receiver.subtypeOf(cglob.receiver) {
        return false
    }

    for i := range cg.globals {
        glob, ok := cglob.labelMap[cg.labels[i].label]
        if !ok {
            return false
        }
        if !cg.globals[i].subtypeOf(glob, visited) {
            return false
        }
    }

    for i := range cglob.globals {
        glob, ok := cg.labelMap[cglob.labels[i].label]
        if !ok {
            return false
        }
        if !cglob.globals[i].subtypeOf(glob, visited) {
            return false
        }
    }
    return true
}

func (cg *choiceGlobal) join(glob global) (global, bool) {
    cglob, ok := glob.(*choiceGlobal)
    if !ok {
        return cg, false
    }

    if !cg.sender.subtypeOf(cglob.sender) {
        return cg, false
    }

    if !cg.receiver.subtypeOf(cglob.receiver) {
        return cg, false
    }

    status := true
    labels := make([]*labelType, 0)
    globals := make([]global, 0)

    for i, lab := range cg.labels {
        glob1, ok := cglob.labelMap[lab.label]
        var glob2 global
        ok3 := false
        if ok {
            glob2, ok3 = cg.globals[i].join(glob1)
            if !ok3 {
                status = false
                glob2 = cg.globals[i]
            }
        } else {
            glob2 = cg.globals[i]
        }
        labels = append(labels, lab)
        globals = append(globals, glob2)
    }

    for i, lab := range cglob.labels {
        if _, ok := cg.labelMap[lab.label]; !ok {
            labels = append(labels, lab)
            globals = append(globals, cglob.globals[i])
        }
    }

    if len(labels) == 0 {
        return cg, false
    }
    return newChoiceGlobal(cg.sender, cg.receiver, labels, globals, cg.line), status
}

func (cg *choiceGlobal) substitute(substitution map[string]string) global {
    senderID, ok := substitution[cg.sender.participant]
    if !ok {
        senderID = cg.sender.participant
    }
    receiverID, ok := substitution[cg.receiver.participant]
    if !ok {
        receiverID = cg.receiver.participant
    }
    sender := newParticipantType(senderID, cg.sender.line)
    receiver := newParticipantType(receiverID, cg.receiver.line)

    globals := make([]global, len(cg.globals))
    for i, glob := range cg.globals {
        globals[i] = glob.substitute(substitution)
    }
    return newChoiceGlobal(sender, receiver, cg.labels, globals, cg.line)
}

func (cg *choiceGlobal) prettyPrint(iw util.IndentedWriter) {
    cg.sender.prettyPrint(iw)
    iw.Print(" -> ")
    cg.receiver.prettyPrint(iw)
    iw.Println(" {")
    for i := range cg.globals {
        if i != 0 {
            iw.Println(" or {")
        }
        iw.Inc()
        cg.labels[i].prettyPrint(iw)
        iw.Print(":\t")
        iw.Inc()
        cg.globals[i].prettyPrint(iw)
        iw.Dec()
        iw.Dec()
        iw.Println()
        iw.Print("}")
    }
}

func (cg *choiceGlobal) String() string {
    s := cg.sender.String() + " -> " + cg.receiver.String() + " { "
    for i := range cg.globals {
        if i != 0 {
            s += " or { "
        }
        s += cg.labels[i].String() + ": "
        s += cg.globals[i].String()
        s += " }"
    }
    return s
}
