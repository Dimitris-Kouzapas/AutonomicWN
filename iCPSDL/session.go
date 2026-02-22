package iCPSDL

import "fmt"
import "iCPSDL/util"

type session interface {
    hash() string
    clone() session
    setFilename(string)
    isEnd() bool
    participants() *util.Set
    rename(map[string]string) session
    project(*local, util.ErrorLog) *local
    typeCheck(*sessionContext, util.ErrorLog)
    prettyPrint(*util.Stream)

    //todo this is for the autonomic supervisor interface
    chann(participant string, channs map[string](map[string]string), counter int) int
}

/********************
 * end
 ********************/

type end struct {
    baseNode
}

func newEnd(line int) (this *end) {
    this = new(end)
    this.init(line)
    return
}

func (this *end) hash() string {
    return util.SequenceHash("end")
}

func (this *end) clone() session {
    return newEnd(0)
}

func (this *end) isEnd() bool {
    return true
}

func (this *end) participants() *util.Set {
    return util.NewSet()
}

func (this *end) rename(context map[string]string) session {
    return newEnd(0)
}

func (this *end) project(lc *local, log util.ErrorLog) *local {
    for participant := range lc.configuration {
        lc.configuration[participant] = newEnd(0)
    }
    return lc
}

func (this *end) typeCheck(context *sessionContext, log util.ErrorLog) {
}

func (this *end) prettyPrint(stream *util.Stream) {
    stream.Print("end")
}

/********************
 * pass
 ********************/

type pass struct {
    baseNode
    from string
    to string
    tdef imodel
    cont session
}

func (this *pass) initPass(from string, to string, tdef imodel, cont session, line int) {
    this.from = from
    this.to = to
    this.tdef = tdef
    this.cont = cont
    this.init(line)
    return
}

func newPass(from string, to string, tdef imodel, cont session, line int) (this *pass) {
    this = new(pass)
    this.initPass(from, to, tdef, cont, line)
    return
}

func (this *pass) hash() string {
    hashCont := this.cont.hash()
    return util.SequenceHash(this.from + "->" + this.to + ": " + this.tdef.getName() + ". ", hashCont)
}

func (this *pass) clone() session {
    return newPass(this.from, this.to, this.tdef, this.cont.clone(), this.line)
}

func (this *pass) participants() *util.Set {
    participants := this.cont.participants()
    participants.Add(this.from, this.to)
    return participants
}

func (this *pass) rename(context map[string]string) session {
    from, ok1 := context[this.from]
    if ok1 == false { from = this.from }
    to, ok2 := context[this.to]
    if ok2 == false { to = this.to }
    cont := this.cont.rename(context)
    return newPass(from, to, this.tdef, cont, this.line)
}

func (this *pass) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    this.tdef.setFilename(filename)
    this.cont.setFilename(filename)
}

func (this *pass) isEnd() bool {
    return false
}

func (this *pass) project(lc *local, log util.ErrorLog) *local {
    lc.configuration[this.from] = nil
    lc.configuration[this.to] = nil
    lc = this.cont.project(lc, log)
    if lc == nil { return nil }
    t1 := lc.configuration[this.from]
    t2 := lc.configuration[this.to]
    lc.configuration[this.from] = newSend(this.to, this.tdef, t1, 0)
    lc.configuration[this.to] = newReceive(this.from, this.tdef, t2, 0)
    return lc
}

func (this *pass) typeCheck(context *sessionContext, log util.ErrorLog) {
    if context.dmn != nil {
        this.tdef.typeCheck(context.dmn, log)
    }
    this.cont.typeCheck(context, log)
}

func (this *pass) prettyPrint(stream *util.Stream) {
    stream.Println(this.from + "->" + this.to + ": " + this.tdef.getName() + ". ")
    this.cont.prettyPrint(stream)
}

/********************
 * send
 ********************/

type send struct {
    pass
}

func newSend(to string, tdef imodel, cont session, line int) (this *send) {
    this = new(send)
    this.initPass(to, to, tdef, cont, line)
    return
}

func (this *send) hash() string {
    hashCont := this.cont.hash()
    return util.SequenceHash(this.to + "!" + this.tdef.getName() + ". ", hashCont)
}

func (this *send) rename(context map[string]string) session {
    to, ok2 := context[this.to]
    if ok2 == false { to = this.to }
    cont := this.cont.rename(context)
    return newSend(to, this.tdef, cont, this.line)
}

func (this *send) prettyPrint(stream *util.Stream) {
    stream.Print(this.to + "!" + this.tdef.getName() + ". ")
    this.cont.prettyPrint(stream)
}

/********************
 * receive
 ********************/

type receive struct {
    pass
}

func newReceive(from string, tdef imodel, cont session, line int) (this *receive) {
    this = new(receive)
    this.initPass(from, from, tdef, cont, line)
    return
}

func (this *receive) hash() string {
    hashCont := this.cont.hash()
    return util.SequenceHash(this.from + "?" + this.tdef.getName() + ". ", hashCont)
}

func (this *receive) rename(context map[string]string) session {
    from, ok1 := context[this.from]
    if ok1 == false { from = this.from }
    cont := this.cont.rename(context)
    return newReceive(from, this.tdef, cont, this.line)
}


func (this *receive) prettyPrint(stream *util.Stream) {
    stream.Print(this.from + "?" + this.tdef.getName() + ". ")
    this.cont.prettyPrint(stream)
}

/********************
 * loop
 ********************/

type loop struct {
    baseNode
    label string
    body session
}

func newLoop(label string, body session, line int) (this *loop) {
    this = new(loop)
    this.label = label
    this.body = body
    this.init(line)
    return
}

//label needs to be irrelevant
func (this *loop) hash() string {
    hashCont := this.body.hash()
    return util.SequenceHash(this.label + ". ", hashCont)
}

func (this *loop) clone() session {
    return newLoop(this.label, this.body.clone(), this.line)
}

func (this *loop) participants() *util.Set {
    return this.body.participants()
}

func (this *loop) rename(context map[string]string) session {
    body := this.body.rename(context)
    return newLoop(this.label, body, this.line)
}

func (this *loop) isEnd() bool {
    return false
}

func (this *loop) project(lc *local, log util.ErrorLog) *local {
    lc = this.body.project(lc, log)
    if lc == nil { return nil }
    for participant, ses := range lc.configuration {
        switch t := ses.(type) {
              case *labelVar:
                  if this.label == t.label {
                      lc.configuration[participant] = newEnd(0)
                  }
              case *end:
                  break
              default:
                  lc.configuration[participant] = newLoop(this.label, ses, 0)
          }
    }
    return lc
}

func (this *loop) typeCheck(context *sessionContext, log util.ErrorLog) {
    context.varMap[this.label] = this
    this.body.typeCheck(context, log)
    delete(context.varMap, this.label)
}

func (this *loop) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    this.body.setFilename(filename)
}

func (this *loop) prettyPrint(stream *util.Stream) {
    stream.Println(this.label + ". ")
    stream.Inc()
    this.body.prettyPrint(stream)
    stream.Dec()
}

/********************
 * labelVar
 ********************/

type labelVar struct {
    baseNode
    label string
}

func newLabel(label string, line int) (this *labelVar) {
    this = new(labelVar)
    this.label = label
    this.init(line)
    return
}

func (this *labelVar) hash() string {
    return util.SequenceHash(this.label)
}

func (this *labelVar) clone() session {
    return newLabel(this.label, this.line)
}

func (this *labelVar) participants() *util.Set {
    return util.NewSet()
}

func (this *labelVar) rename(context map[string]string) session {
    return newLabel(this.label, this.line)
}

func (this *labelVar) isEnd() bool {
    return false
}

func (this *labelVar) project(lc *local, log util.ErrorLog) *local {
    for participant := range lc.configuration {
        lc.configuration[participant] = newLabel(this.label, 0)
    }
    return lc
}

func (this *labelVar) typeCheck(context *sessionContext, log util.ErrorLog) {
    _, ok := context.varMap[this.label]
    if ok == false {
        error := fmt.Sprintf("Loop label %v is not defined.", this.label)
        this.reportError(error, log)
    }
}

func (this *labelVar) prettyPrint(stream *util.Stream) {
    stream.Print(this.label)
}

/********************
 * choice
 ********************/

type choice struct {
    baseNode
    from string
    to string
    enum string
    sessions map[string]session
}

func newChoice(from string, to string, enum string, line int) (this *choice) {
    this = new(choice)
    this.initChoice(from, to, enum, line)
    return
}

func (this *choice) hash() string {
    hashes := make([]string, 0)
    for _, sess := range this.sessions {
        hashes = append(hashes, sess.hash())
    }
    hashCont := util.CommutativeHash(hashes)
    return util.SequenceHash(this.from + "->" + this.to + ": " + this.enum + " ", hashCont)
}

func (this *choice) initChoice(from string, to string, enum string, line int) {
    this.from = from
    this.to = to
    this.enum = enum
    this.sessions = make(map[string]session)
    this.init(line)
    return
}

func (this *choice) clone() session {
    tmp := newChoice(this.from, this.to, this.enum, this.line)
    for label, sess := range this.sessions {
        tmp.addSession(label, sess.clone())
    }
    return tmp
}

func (this *choice) rename(context map[string]string) session {
    from, ok1 := context[this.from]
    if ok1 == false { from = this.from }
    to, ok2 := context[this.to]
    if ok2 == false { to = this.to }
    choice := newChoice(from, to, this.enum, this.line)
    for label, sess := range this.sessions {
        rsess := sess.rename(context)
        choice.addSession(label, rsess)
    }
    return choice
}

func (this *choice) participants() *util.Set {
    participants := util.NewSet()
    for _, sess := range this.sessions {
        parts := sess.participants()
        participants.Add(parts.Elements()...)
    }
    participants.Add(this.from, this.to)
    return participants
}

func (this *choice) addSession(label string, s session) {
    this.sessions[label] = s
}

func (this *choice) isEnd() bool {
    return false
}

func (this *choice) project(lc *local, log util.ErrorLog) *local {
    lc.configuration[this.from] = nil
    lc.configuration[this.to] = nil
    from := newSelect(this.to, this.enum, 0)
    to := newBranch(this.from, this.enum, 0)

    lctmp := make([]*local, 0)
    i := 0
    for label, sess := range this.sessions {
        lctmp = append(lctmp, lc.clone())
        lctmp[i] = sess.project(lctmp[i], log)
        if lctmp[i] == nil { return nil  }
        t1 := lctmp[i].configuration[this.from]
        t2 := lctmp[i].configuration[this.to]
        from.addSession(label, t1)
        to.addSession(label, t2)
        i += 1
    }
    for j := 1; j < len(lctmp); j++ {
        if lctmp[0].hash(this.from, this.to) != lctmp[j].hash(this.from, this.to) {
            error := fmt.Sprintf("Global type cannot be projected.")
            log.Add(util.NewExecutionException(error))
            return nil
        }
    }
    lctmp[0].configuration[this.from] = from
    lctmp[0].configuration[this.to] = to
    return lctmp[0]
}

func (this *choice) typeCheck(context *sessionContext, log util.ErrorLog) {
    if context.dmn != nil {
        t, ok1 := context.dmn.properties[this.enum]
        ok2 := false
        var enum *enumProperty = nil
        if ok1 == true {
            enum, ok2 = t.(*enumProperty)
        }
        if ok1 == false || ok2 == false {
            error := fmt.Sprintf("Enumeration type %v is not defined", this.enum)
            this.reportError(error, log)
        } else {
            for label := range this.sessions {
                flag := false
                for _, item := range enum.labels {
                    if item == label {
                        flag = true
                        break
                    }
                }
                if flag == false {
                    error := fmt.Sprintf("Label %v is not defined for enumeration type %v.", label, this.enum)
                    this.reportError(error, log)
                }
            }
            for _, label := range enum.labels {
                _, ok := this.sessions[label]
                if ok == false {
                    error := fmt.Sprintf("Enumeration label %v is not defined by the selection", label)
                    this.reportError(error, log)
                }
            }
        }
    }
    for _, s := range this.sessions {
        s.typeCheck(context, log)
    }
}

func (this *choice) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    for _, s := range this.sessions {
        s.setFilename(filename)
    }
}

func (this *choice) prettyPrint(stream *util.Stream) {
    stream.Print(this.from + "->" + this.to + ": " + this.enum + " ")
    this.prettyPrint_(stream)
}

func (this *choice) prettyPrint_(stream *util.Stream) {
    i := 0
    for label, s := range this.sessions {
        if i != 0 {
            stream.Print(" or ")
        }
        stream.Println("{ ")
        stream.Inc()
        stream.Print(label + ": ")
        s.prettyPrint(stream)
        stream.Println("")
        stream.Dec()
        stream.Print(" }")
        i++
    }
}

/********************
 * select
 ********************/

type sel struct {
    choice
}

func newSelect(to string, enum string, line int) (this *sel) {
    this = new(sel)
    this.initChoice(to, to, enum, line)
    return
}

func (this *sel) hash() string {
    hashes := make([]string, 0)
    for _, sess := range this.sessions {
        hashes = append(hashes, sess.hash())
    }
    hashCont := util.CommutativeHash(hashes)
    return util.SequenceHash(this.choice.to + "!" + this.choice.enum + " ", hashCont)
}

func (this *sel) rename(context map[string]string) session {
    to, ok2 := context[this.to]
    if ok2 == false { to = this.to }
    sel := newSelect(to, this.enum, this.line)
    for label, sess := range this.sessions {
        rsess := sess.rename(context)
        sel.addSession(label, rsess)
    }
    return sel
}

func (this *sel) prettyPrint(stream *util.Stream) {
    stream.Print(this.choice.to + "!" + this.choice.enum + " ")
    this.choice.prettyPrint_(stream)
}

/********************
 * branch
 ********************/

type branch struct {
    choice
}

func newBranch(from string, enum string, line int) (this *branch) {
    this = new(branch)
    this.initChoice(from, from, enum, line)
    return
}

func (this *branch) hash() string {
    hashes := make([]string, 0)
    for _, sess := range this.sessions {
        hashes = append(hashes, sess.hash())
    }
    hashCont := util.CommutativeHash(hashes)
    return util.SequenceHash(this.choice.from + "?" + this.choice.enum + " ", hashCont)
}

func (this *branch) rename(context map[string]string) session {
    from, ok1 := context[this.from]
    if ok1 == false { from = this.from }
    branch := newBranch(from, this.enum, this.line)
    for label, sess := range this.sessions {
        rsess := sess.rename(context)
        branch.addSession(label, rsess)
    }
    return branch
}

func (this *branch) prettyPrint(stream *util.Stream) {
    stream.Print(this.from + "?" + this.enum + " ")
    this.choice.prettyPrint_(stream)
}

/******************************************************************************
 * typecheck context
 ******************************************************************************/

type sessionContext struct {
    dmn *domain
    varMap map[string]*loop
}

func newSessionContext(dmn *domain) (this *sessionContext) {
    this = new(sessionContext)
    this.dmn = dmn
    this.varMap = make(map[string]*loop)
    return
}

/******************************************************************************
 * local configuration
 ******************************************************************************/

type local struct {
    baseNode
    configuration map[string]session
}

func newLocal(line int) (this *local) {
    this = new(local)
    this.configuration = make(map[string]session)
    this.init(line)
    return
}

func (this *local) hash(from string, to string) string {
    hashes := make([]string, 0)
    for participant, sess := range this.configuration {
        if participant == from || participant == to { continue }
        hashes = append(hashes, participant + sess.hash())
    }
    return util.CommutativeHash(hashes)
}

func (this *local) clone() *local {
    clone := newLocal(this.line)
    for participant, sess := range this.configuration {
        if sess != nil {
            clone.configuration[participant] = sess.clone()
        } else {
            clone.configuration[participant] = sess
        }
    }
    return clone
}

func (this *local) add(id string, sess session) *local {
    this.configuration[id] = sess
    return this
}

func (this *local) compose(log util.ErrorLog) *global {
    sess := compose(this.configuration, log)
    if log.HasErrors() {
        return nil
    }
    return newGlobal(sess, 0)
}

func (this *local) typeCheck(log util.ErrorLog) {
    for _, ses := range this.configuration {
        c := newSessionContext(nil)
        ses.typeCheck(c, log)
    }
}

func (this *local) execute(context *iCPSDL) (expression, util.ErrorLog) {
    log := util.NewErrorLog()
    this.typeCheck(log)
    return this, log
}

func (this *local) render() {}

func (this *local) prettyPrint(stream *util.Stream) {
    stream.Println("local {")
    stream.Inc()
    for p, role := range this.configuration {
        stream.Println(p + " = ")
        stream.Inc()
        role.prettyPrint(stream)
        stream.Dec()
        stream.Println("")
        stream.Println("")
    }
    stream.Dec()
    stream.Println("}")
    return
}

func (this *local) getType() string {
    return "local configuration"
}

func (this *local) getModule() string {
    return "*"
}

func (this *local) info() string {
    return fmt.Sprintf("Type: local configuration, Roles: %v", len(this.configuration))
}

func compose(roles map[string]session, log util.ErrorLog) session {
    endFlag := false
    loopLabel := ""
    loop:
    for p, t := range roles {
        switch s := t.(type) {
            case *end:
                endFlag = true
            case *labelVar:
                if loopLabel == "" {
                    loopLabel = s.label
                } else if loopLabel != s.label {
                    error := fmt.Sprintf("Participant % v ends with label %v, instead of label %v.", p, s.label, loopLabel)
                    log.Add(util.NewExecutionException(error))
                    return nil
                }
            default:
                endFlag = false
                loopLabel = ""
                break loop
        }
    }
    if endFlag == true && loopLabel == "" {
        return newEnd(0)
    }
    if loopLabel != "" {
        return newLabel(loopLabel, 0)
    }
    for p, t1 := range roles {
        switch fst := t1.(type) {
            case *receive:
            case *branch:
            case *end:
            case *loop:
            case *send:
                for q, t2 := range roles {
                    snd, ok := t2.(*receive)
                    if ok == false {
                        continue
                    }
                    if (p == snd.from && fst.to == q) == false {
                        continue
                    }
                    if fst.tdef.toString() != snd.tdef.toString() {
                        error := fmt.Sprintf("Mismatched carried types. Participant %v: %v, participant %v: %v.", p, fst.tdef.getName(), q, snd.tdef.getName())
                        log.Add(util.NewExecutionException(error))
                        return nil
                    }
                    roles[p] = fst.cont
                    roles[q] = snd.cont
                    sess := compose(roles, log)
                    roles[p] = t1
                    roles[q] = t2
                    return newPass(p, q, fst.tdef, sess, 0)
                }
            case *sel:
                for q, t2 := range roles {
                    snd, ok := t2.(*branch)
                    if ok == false {
                        continue
                    }
                    if (p == snd.from && fst.to == q) == false {
                        continue
                    }
                    if fst.enum != snd.enum {
                        error := fmt.Sprintf("Mismatched enum types. Participant %v: %v, participant %v: %v.", p, fst.enum, q, snd.enum)
                        log.Add(util.NewExecutionException(error))
                        return nil
                    }
                    ch := newChoice(p, q, fst.enum, 0)
                    for label, session1 := range fst.sessions {
                        session2, ok := snd.sessions[label]
                        if ok == false {
                            error := fmt.Sprintf("Mismatched enum types. Participant %v: %v.", q, label)
                            log.Add(util.NewExecutionException(error))
                            return nil
                        }
                        roles[p] = session1
                        roles[q] = session2
                        sess := compose(roles, log)
                        ch.addSession(label, sess)
                        roles[p] = t1
                        roles[q] = t2
                        if log.HasErrors() {
                            return nil
                        }
                    }
                    return ch
                }
        }
    }
    loopLabel = ""
    context := make(map[string]session)
    for p, t := range roles {
        context[p] = t
        if t.isEnd() == true {
            context[p] = t
            continue
        }
        l, ok := t.(*loop)
        if ok == false {
            error := fmt.Sprintf("No matching input output between participants.")
            log.Add(util.NewExecutionException(error))
            return nil
        }
        context[p] = l.body
        if loopLabel == "" {
            loopLabel = l.label
        } else if loopLabel != l.label {
            error := fmt.Sprintf("Participant %v does not implement a loop with label. Implements label %v instead.", p, loopLabel, l.label)
            log.Add(util.NewExecutionException(error))
            return nil
        }
    }
    sess := compose(context, log)
    if log.HasErrors() {
        return nil
    }
    return newLoop(loopLabel, sess, 0)
}

/******************************************************************************
 * global configuration
 ******************************************************************************/

type global struct {
    baseNode
    configuration session
}

func newGlobal(configuration session, line int) (this *global) {
    this = new(global)
    this.configuration = configuration
    this.init(line)
    return
}

func (this *global) project(log util.ErrorLog) *local {
    lc := newLocal(this.line)
    return this.configuration.project(lc, log)
}

func (this *global) typeCheck(log util.ErrorLog) {
    c := newSessionContext(nil)
    this.configuration.typeCheck(c, log)
}

func (this *global) execute(context *iCPSDL) (expression, util.ErrorLog) {
    log := util.NewErrorLog()
    this.typeCheck(log)
    return this, log
}

func (this *global) render() {}

func (this *global) prettyPrint(stream *util.Stream) {
    stream.Println("global ")
    stream.Inc()
    this.configuration.prettyPrint(stream)
    stream.Dec()
}

func (this *global) getType() string {
    return "global configuration"
}

func (this *global) getModule() string {
    return "*"
}

func (this *global) info() string {
    return "Type: global configuration."
}
