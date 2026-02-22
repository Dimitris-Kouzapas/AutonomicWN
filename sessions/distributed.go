package ast

import (
	"fmt"
	"os"
	"sync"
	"reflect"
	"sync/atomic"
	"context"
    "sessions/util"
)


type channelMap map[string](map[string] util.Channel[expression])
type brokerChannel chan channelMap

/******************************************************************************
 * broker interface
 ******************************************************************************/

type broker interface {
	expression
	close()
}

type recBrokerI interface {
	recurse(string) (channelMap, bool)
}

/******************************************************************************
 * broker Impl
 ******************************************************************************/

type brokerImpl struct {
	baseNode
	tdef 		typedef
	requester 	*participantExpr
	syncType 	string

	gconfig 	globalConfig
	channels 	channelMap

	closed		atomic.Bool
}


func (b *brokerImpl) setFilename(filename string) {
	b.baseNode.setFilename(filename)
	b.tdef.setFilename(filename)
}

func (b *brokerImpl) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
	b.tdef.typeCheck(ctx, log)
	if tdef := b.tdef.getType(); tdef != nil {
		ok := true
		if b.gconfig, ok = tdef.(globalConfig); !ok {
			b.reportErrorf(log, "Expecting global configuration type. Found %s instead", b.tdef.String())
		}
	}
}

func (b *brokerImpl) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	b.tdef.projectionCheck(ctx, elog, rlog)
}

func (_ *brokerImpl) sessionCheck(_ *sessionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}
func (_ *brokerImpl) operation(_ *evaluationContext, _ expression, _ string) expression { return nil }
func (b *brokerImpl) prettyPrint(iw util.IndentedWriter) {
	iw.Println(b.syncType + " { ")
    iw.Inc()
    iw.Print("chan: ")
    iw.Inc()
    b.gconfig.prettyPrint(iw)
    iw.Dec()
    iw.Println()
    if b.requester != nil {
        iw.Print("requester: ")
        iw.Inc()
        b.requester.prettyPrint(iw)
    }
    iw.Println()
    iw.Dec()
    iw.Dec()
    iw.Println("}")
}

func (b *brokerImpl) String() string {
	s := b.syncType + " { chan: " + b.gconfig.String() + ";" 
    if b.requester != nil {
        s += "requester " + b.requester.String() + ";"
    }
    s += " }"
    return s
}

func (_ *brokerImpl) goCode(_ util.IndentedWriter) {}

func (b *brokerImpl) close() {
	b.closed.Store(true)
}

const ctxKey = ":ctx"
func tryReceiveFromChannels(ctx context.Context, chans map[string]chan brokerChannel) (string, brokerChannel, bool) {
	cases := make([]reflect.SelectCase, 0, len(chans)+1)
	keys  := make([]string, 0, len(chans)+1)

	for key, ch := range chans {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
		keys = append(keys, key)
	}

	// Add the context as an extra receive case (if provided)
	if ctx != nil {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()), // <-chan struct{}
		})
		keys = append(keys, ctxKey)
	} else if len(cases) == 0 {
		return "", nil, false
	}

	// Optional non-blocking default:
	// cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	// keys  = append(keys, ":default")

	chosen, recv, ok := reflect.Select(cases)
	key := keys[chosen]

	// Context selected: Done() is a closed chan → ok == false.
	if key == ctxKey {
		return ctxKey, nil, false
	}

	// A real channel closed → tell caller to drop it.
	if !ok {
		return key, nil, false
	}

	val, _ := recv.Interface().(brokerChannel)
	return key, val, true
}

/******************************************************************************
 * fast recursion broker 
 ******************************************************************************/

type fastRecBroker struct {
	brokerImpl
	participants map[string]bool
}

func newFastRecBroker(tdef typedef, line int) *fastRecBroker {
	return &fastRecBroker {
		brokerImpl : brokerImpl {
			baseNode: baseNode{line: line},
			tdef: tdef,
			syncType: "async",
		},
		participants: make(map[string]bool),
	}
}

func (b *fastRecBroker) getType() typedef {
	tdef := b.tdef.getType()
	return newBrokerType(tdef, nil, b.line)
}

func (b *fastRecBroker) expressionCheck(ctx *expressionCheckContext, log util.ErrorLog) typedef {
	for _, p := range b.gconfig.participants() {
		b.participants[p] = true
	}
	b.channels = b.gconfig.channels(true)
	return newBrokerType(b.gconfig, nil, b.line)
}

func (b *fastRecBroker) evaluate(_ *evaluationContext) expression { return b }

func (b *fastRecBroker) recurse(participant string) (channelMap, bool) {
	if _, ok := b.participants[participant]; !ok {
		msg := b.runtimeErrorf("unknown participant: %q.", participant)
		fmt.Fprintln(os.Stderr, msg)
		return b.channels, false
	}

	return b.channels, true
}

/******************************************************************************
 * recursion broker 
 ******************************************************************************/

type recBroker struct {
	fastRecBroker
	registered  map[string]bool

	running atomic.Bool

	mu          sync.Mutex
	ready		*sync.Cond
}

func newRecBroker(tdef typedef, line int) *recBroker {
	b := &recBroker {
		fastRecBroker: fastRecBroker {
			brokerImpl : brokerImpl {
				baseNode: baseNode{line: line},
				tdef: tdef,
				syncType: "sync",
			},
			participants: make(map[string]bool),
		},
		registered: make(map[string]bool),
	}
	b.ready = sync.NewCond(&b.mu)
	return b
}

func (b *recBroker) expressionCheck(ctx *expressionCheckContext, log util.ErrorLog) typedef {
	b.running.Store(true)
	return b.fastRecBroker.expressionCheck(ctx, log)
}

func (b *recBroker) evaluate(_ *evaluationContext) expression { return b }

func (b *recBroker) recurse(participant string) (channelMap, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.registered) == 0 {
		if b.closed.Load() {
			b.running.Store(false)
		}
		for _, p := range b.gconfig.participants() {
			b.registered[p] = true
		}
		b.ready.Broadcast()
	}

	if _, ok := b.participants[participant]; !ok {
		msg := b.runtimeErrorf("unknown participant: %q.", participant)
		fmt.Fprintln(os.Stderr, msg)
		return b.channels, false
	}

	_, ok := b.registered[participant]
	if !ok {
		b.ready.Wait()
	}

	delete(b.registered, participant)
	return b.channels, b.running.Load()
}

/******************************************************************************
 * broker
 ******************************************************************************/

type reqBroker struct {
	brokerImpl

	responseChannels map[string]brokerChannel
	cMap map[string](chan brokerChannel)

	ctx context.Context
	cancel context.CancelFunc
}

func newReqBroker(requester *participantExpr, tdef typedef, line int) *reqBroker {
	return &reqBroker {
		brokerImpl : brokerImpl {
			baseNode:	baseNode{line: line},
			requester:	requester,
			tdef:		tdef,
			syncType:	"sync",
		},
		responseChannels: make(map[string]brokerChannel),
	}
}

func (b *reqBroker) setFilename(filename string) {
	b.brokerImpl.setFilename(filename)
	b.requester.setFilename(filename)
}

func (b *reqBroker) getType() typedef {
	tdef := b.tdef.getType()
	return newBrokerType(tdef, b.requester.pType(), b.line)
}

func (b *reqBroker) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
	b.brokerImpl.typeCheck(ctx, log)
	b.requester.typeCheck(ctx, log)
}

func (b *reqBroker) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	b.brokerImpl.projectionCheck(ctx, elog, rlog)
	b.requester.projectionCheck(ctx, elog, rlog)
}

func (b *reqBroker) expressionCheck(ctx *expressionCheckContext, log util.ErrorLog) typedef {
	participants := b.gconfig.participants()

	found := false
	for _, participant := range participants {
		if participant == b.requester.id {
			found = true
			break
		}
	}

	if !found {
		stream := util.NewStream().Inc().Inc()
		b.gconfig.prettyPrint(stream)
		b.requester.reportErrorf(
			log,
			"broker expression: qlobal configuration:\n%s\n\tdoes not define request participant %q",
			stream.String(), b.requester.String(),
		)
		return nil
	}
	b.channels = b.gconfig.channels(true)
	b.ctx, b.cancel = context.WithCancel(context.Background())
	go b.run()
	return newBrokerType(b.gconfig, b.requester.pType(), b.line)
}

func (b *reqBroker) sessionCheck(ctx *sessionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	b.requester.sessionCheck(ctx, elog, rlog)
}

func (b *reqBroker) evaluate(_ *evaluationContext) expression { return b }
// func (b *reqBroker) prettyPrint(iw util.IndentedWriter) {
// 	b.brokerImpl.prettyPrint(iw)
// 	iw.Print(" on ")
// 	b.req.prettyPrint(iw)
// }

// func (b *reqBroker) String() string {
// 	return b.brokerImpl.String() + " on " + b.req.String()
// }

func (b *reqBroker) close() {
	b.brokerImpl.close()
	if b.cancel != nil {
	 	b.cancel()
	}
}

func (b *reqBroker) run() {
	b.cMap = make(map[string](chan brokerChannel), len(b.gconfig.participants()))
	for _, participant := range b.gconfig.participants() {
		b.cMap[participant] = make(chan brokerChannel)
	}
	// defer b.cleanup() <- TODO how to call cleanup without selector problems

	cMap := make(map[string]chan brokerChannel, len(b.cMap))
	for k, v := range b.cMap { cMap[k] = v }

	for {
		key, val, ok := tryReceiveFromChannels(b.ctx, cMap)
		if key == ctxKey {
			return
		}
		if !ok {
			// Either empty (shouldn't happen here) or chosen chan closed: drop it
			delete(cMap, key)
			continue
		}
		// Remove comment when the selector in tryReceiveFromChannels(cMap) has a "default" case.
		// Need to be careful, whenever a channel is named as "default" 
		// if key == "default" { } 

		b.responseChannels[key] = val
		delete(cMap, key)

		if len(cMap) == 0 {
			if b.closed.Load() {
				return
			}
			//TODO: what if the broker closes here... <- it is a concurrency bug
			// channels := b.gconfig.channels(true)

			for _, rc := range b.responseChannels {
				go func(c brokerChannel) { 
					select {
						case c <- b.channels:
						case <- b.ctx.Done():
					}
				}(rc)
			}

			// Reset state for next round
			b.responseChannels = make(map[string]brokerChannel, len(b.cMap))
			for k, v := range b.cMap { cMap[k] = v }
		}		
	}
}

func (b *reqBroker) cleanup() {
	for _, v := range b.cMap {
		close(v)
	}
}

func (b *reqBroker) accept(participant string, c brokerChannel) func() (channelMap, bool) {
	if ch, ok := b.cMap[participant]; ok && b.requester.id != participant {
		select {
			case ch <- c:
			case <- b.ctx.Done():
		}
	} else {
		msg := b.runtimeErrorf("accept on unknown participant: %q.", participant)
		fmt.Fprintln(os.Stderr, msg)
	}

	return func() (channelMap, bool) {
		select {
			case channels, open := <- c:
				return channels, open
			case <- b.ctx.Done(): 
				return nil, false
		}
	}
}

func (b *reqBroker) request(participant string, c brokerChannel) func() (channelMap, bool) {
	if b.requester.id == participant {
		select {
			case b.cMap[participant] <- c:
			case <- b.ctx.Done():
		}
	} else {
		msg := b.runtimeErrorf("request on unknown participant: %q.", participant)
		fmt.Fprintln(os.Stderr, msg)
	}

	return func() (channelMap, bool) {
		select {
			case channels, open := <- c:
				return channels, open
			case <- b.ctx.Done(): 
				return nil, false
		}
	}

}

/******************************************************************************
 * init process
 ******************************************************************************/

type initProc struct {
	baseNode
	channel *variableExpr
	participants []*participantExpr
	applications []*application
	btype *brokerType
	gconfig globalConfig
	locals map[string]local

	req *participantExpr
}

func (ip *initProc) setFilename(filename string) {
	ip.baseNode.setFilename(filename)
	ip.channel.setFilename(filename)
	for i := range ip.participants {
		ip.participants[i].setFilename(filename)
		ip.applications[i].setFilename(filename)
	}
	if ip.req != nil {
		ip.req.setFilename(filename)
	}
}

func (ip *initProc) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
	ip.channel.typeCheck(ctx, log)
	for i := range ip.participants {
		ip.participants[i].typeCheck(ctx, log)
		ip.applications[i].typeCheck(ctx, log)
	}
	if ip.req != nil {
		ip.req.typeCheck(ctx, log)
	}
}

func (ip *initProc) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	for i := range ip.applications {
		ip.participants[i].projectionCheck(ctx, elog, rlog)
		ip.applications[i].projectionCheck(ctx, elog, rlog)
	}
	if ip.req != nil {
		ip.req.projectionCheck(ctx, elog, rlog)
	}
}

func (ip *initProc) brokerCheck(ctx *expressionCheckContext, log util.ErrorLog) bool {
	tdef := ip.channel.expressionCheck(ctx, log)
	if tdef == nil {
		return false
	}
	var ok bool
	ip.btype, ok = tdef.getType().(*brokerType)
	if !ok {
		ip.reportErrorf(log, "expected broker type; instead found: %s", tdef.String())
		return ok
	}
	ip.channel.setType(ip.btype)

	ip.gconfig, ok = ip.btype.gconfig.getType().(globalConfig)
	if !ok {
		ip.reportErrorf(log, "exppecting global configuration for broker type; instead found %s", ip.btype.gconfig.String())
		return ok
	}

	for _, participant := range ip.participants {
		if participant.id == ip.btype.requester.participant {
			participant.reportErrorf(log, "cannot accept on request participant: %q", participant)
		}
	}

	return ok
}

func (ip *initProc) applicationCheck(ctx *expressionCheckContext, log util.ErrorLog) {
	parSet := util.NewHashSet[string]()
	for _, participant := range ip.gconfig.participants() {
		parSet.Add(participant)
	}

	for _, p := range ip.participants {
		if !parSet.Contains(p.id) {
			stream := util.NewStream().Inc().Inc()
			stream.Print(ip.channel.id + ": ")
			ip.gconfig.prettyPrint(stream)
			p.reportErrorf(log, "broker channel:\n%s\n\tdoes not define participant %q.", stream.String(), p.String())
		}
	}

	ip.locals = make(map[string]local, len(ip.applications))
	for i := range ip.applications {
		ip.locals[ip.participants[i].id] = ip.applications[i].expressionCheck(ip.participants[i], parSet, ctx, log)
	}
}

func (ip *initProc) sessionCheck(ctx *sessionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	for i := range ip.applications {
		ip.participants[i].sessionCheck(ctx, elog, rlog)
		ip.applications[i].sessionCheck(ctx, elog, rlog)
	}
	var participants []*participantExpr
	if ip.req != nil {
		participants = append(ip.participants, ip.req)
	} else {
		participants = ip.participants
	}
	for _, pexpr := range participants {
		loc, present := ip.locals[pexpr.id]
		if !present {
			// prior error likely; skip to avoid nil deref
			continue
		}
		locAbstr := ip.gconfig.project(pexpr.pType())
		if !loc.subtypeOf(locAbstr.loc, util.NewHashSet[typePair]()) {
			stream1 := util.NewStream().Inc().Inc()
			stream1.Print(pexpr.String() + ": ")
			loc.prettyPrint(stream1)

			stream2 := util.NewStream().Inc().Inc()
			stream2.Print(ip.channel.id + ": ")
			ip.gconfig.prettyPrint(stream2)

			pexpr.reportErrorf(elog, "role:\n%v\n\tis not projectable from broker channel:\n%v", stream1.String(), stream2.String())
		}
	}
}

func (ip *initProc) evaluate(ctx *evaluationContext) channelMap {
	v := ip.channel.evaluate(ctx)
	broker, ok := v.(*reqBroker)
	if !ok {
		msg := ip.runtimeErrorf("undefined broker: %q.", ip.channel.id)
		fmt.Fprintln(os.Stderr, msg)
		return nil
	}

	chans := make([]brokerChannel, 0, len(ip.participants))
	var running bool
	funcs := make([]func() (channelMap, bool), 0, len(ip.participants))
	for _, pexpr := range ip.participants {
		c := make(brokerChannel)
		chans = append(chans, c)
		// running = broker.accept(pexpr.id, c)
		funcs = append(funcs, broker.accept(pexpr.id, c))
	}

	var c brokerChannel
	var reqFunc func() (channelMap, bool)
	if ip.req != nil {
		c = make(brokerChannel)
		reqFunc = broker.request(ip.req.id, c)
	}

	var channels channelMap
	for i := range funcs {//chans {
		// channels = <-chans[i]
		channels, running = funcs[i]()
	}

	if ip.req != nil {
		// channels = <- c
		channels, running = reqFunc()
	}
	if running || channels != nil {
		for i := range ip.applications {
			ip.applications[i].evaluate(ip.participants[i], channels, ctx)
		}
	} else {
		channels = ip.gconfig.channels(false)
	}

	// TODO again when to release resources.
	for i := range chans { close(chans[i]) }
	if ip.req != nil { close(c) }

	return channels
}

func (ip *initProc) prettyPrint(iw util.IndentedWriter) {
	iw.Print("(" + ip.channel.id + ")")
	if ip.req != nil {
		iw.Print(" as role " + ip.req.id)
	}
	iw.Println(" with")
	iw.Inc()
	for i := range ip.applications {
		iw.Print("role ")
		ip.participants[i].prettyPrint(iw)
		iw.Print(" as ")
		ip.applications[i].prettyPrint(iw)
		iw.Println(";")
	}
	iw.Dec()
	// iw.Println()
}

func (ip *initProc) String() string {
	s := "(" + ip.channel.id + ")"
	if ip.req != nil {
		s += " role as " + ip.req.id
	}
	s += " with "
	for i := range ip.applications {
		s += "role " + ip.participants[i].String() + " as " + ip.applications[i].String() + "; "
	}
	// s += " }"
	return s
}

func (_ *initProc) goCode(_ util.IndentedWriter) {}

/******************************************************************************
 * accept process
 ******************************************************************************/

type acceptProc struct {
	initProc
}

func newAcceptProc(channel *variableExpr, participants []*participantExpr, applications []*application, line int) *acceptProc {
	return &acceptProc {
		initProc: initProc {
			baseNode: baseNode{line: line},
			channel: channel,
			participants: participants,
			applications: applications,
		},
	}
}

func (ap *acceptProc) expressionCheck(ctx *expressionCheckContext, lin *linearContext, log util.ErrorLog) *linearContext {
	ok := ap.initProc.brokerCheck(ctx, log)

	set := util.NewHashSet[string]()
	for _, pexpr := range ap.participants {
		if set.Contains(pexpr.String()) {
			pexpr.reportErrorf(log, "duplicate definition of participant: %q.", pexpr.String())
		}
		set.Add(pexpr.String())
	}

	if !ok {
		return lin.newEndLocal(ap.line)
	}

	ap.applicationCheck(ctx, log)
	return lin.newEndLocal(ap.line)
}

func (ap *acceptProc) evaluate(ctx *evaluationContext) {
	ap.initProc.evaluate(ctx)
}

func (ap *acceptProc) prettyPrint(iw util.IndentedWriter) {
	iw.Print("acc ")
	ap.initProc.prettyPrint(iw)
}

func (ap *acceptProc) String() string {
	return "acc" + ap.initProc.String()
}

/******************************************************************************
 * request process
 ******************************************************************************/

type requestProc struct {
	initProc
	cont process
}

func newRequestProc(channel *variableExpr, req *participantExpr, participants []*participantExpr, applications []*application, cont process, line int) *requestProc {
	return &requestProc {
		initProc: initProc {
			baseNode: baseNode{line: line},
			channel: channel,
			participants: participants,
			applications: applications,
			req: req,
		},
		cont: cont,
	}
}

func (rp *requestProc) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	rp.initProc.projectionCheck(ctx, elog, rlog)
	rp.cont.projectionCheck(ctx, elog, rlog)
}

func (rp *requestProc) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
	rp.initProc.typeCheck(ctx, log)
	rp.cont.typeCheck(ctx, log)
}

func (rp *requestProc) expressionCheck(ctx *expressionCheckContext, lin *linearContext, log util.ErrorLog) *linearContext {
	ok := rp.initProc.brokerCheck(ctx, log)

	if !ok {
		return lin.newEndLocal(rp.line)
	}

	if rp.req.id != rp.btype.requester.participant {
		rp.req.reportErrorf(
			log, 
			"Invalid request participant %q; expecting participant %q.",
			rp.req.String(), rp.btype.requester.String(),
		)
	}

	participants := make([]*participantExpr, len(rp.gconfig.parameterTypes()))
	for i, parameter := range rp.gconfig.parameterTypes() {
		participants[i] = newParticipantExpr(parameter.participant, parameter.line)
	}

	for _, participant := range participants {
		if rp.req.id == participant.id {
			continue
		}
		if !ctx.addParticipant(participant) {
			rp.reportErrorf(log, "duplicate definition of participant: %q.", participant.String())
		}
	}

	rp.initProc.applicationCheck(ctx, log)

	lin = lin.newSession(participants, rp.line)
	lin = rp.cont.expressionCheck(ctx, lin, log)
	rp.locals[rp.req.id] = lin.removeSession(participants, rp.line)

	for _, participant := range participants {
		if rp.req.id == participant.id {
			continue
		}
		ctx.removeParticipant(participant)
	}

	return lin
}

func (rp *requestProc) sessionCheck(ctx *sessionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	rp.initProc.sessionCheck(ctx, elog, rlog)
	rp.cont.sessionCheck(ctx, elog, rlog)
}

func (rp *requestProc) evaluate(ctx *evaluationContext) {
	channels := rp.initProc.evaluate(ctx)

	participants := make([]*participantExpr, len(rp.gconfig.parameterTypes()))
	for i, parameter := range rp.gconfig.parameterTypes() {
		participants[i] = newParticipantExpr(parameter.participant, parameter.line)
	}

	for _, p := range participants {
		if rp.req.id == p.id {
			continue
		}
		ctx.addParticipantChannel(p, channels[rp.req.id][p.id], channels[p.id][rp.req.id])
	}

	rp.cont.evaluate(ctx)
}


func (rp *requestProc) prettyPrint(iw util.IndentedWriter) {
	iw.Print("req ")
	rp.initProc.prettyPrint(iw)
	rp.cont.prettyPrint(iw)
	iw.Println()
}

func (rp *requestProc) String() string {
	return "req" + rp.initProc.String()
}

/******************************************************************************
 * recurse process
 ******************************************************************************/

type recurseProc struct {
	baseNode

	channel         *variableExpr
	rec             *participantExpr
	app             *application
	participants    []*participantExpr
	applications    []*application

	btype *brokerType
	gconfig globalConfig
	locals map[string]local
}

func newRecurseProc(channel *variableExpr, rec *participantExpr, app *application, participants []*participantExpr, applications []*application, line int) *recurseProc {
	return &recurseProc {
		baseNode: baseNode{line: line},
		channel: channel,
		rec: rec,
		app: app,
		participants: participants,
		applications: applications,
	}
}

func (rp *recurseProc) setFilename(filename string) {
	rp.baseNode.setFilename(filename)
	rp.channel.setFilename(filename)
	rp.rec.setFilename(filename)
	rp.app.setFilename(filename)
	for i := range rp.participants {
		rp.participants[i].setFilename(filename)
		rp.applications[i].setFilename(filename)
	}
}

func (rp *recurseProc) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
	rp.channel.typeCheck(ctx, log)
	rp.rec.typeCheck(ctx, log)
	rp.app.typeCheck(ctx, log)
	for i := range rp.participants {
		rp.participants[i].typeCheck(ctx, log)
		rp.applications[i].typeCheck(ctx, log)
	}
}

func (rp *recurseProc) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
	rp.channel.projectionCheck(ctx, elog, rlog)
	rp.rec.projectionCheck(ctx, elog, rlog)
	rp.app.projectionCheck(ctx, elog, rlog)
	for i := range rp.applications {
		rp.participants[i].projectionCheck(ctx, elog, rlog)
		rp.applications[i].projectionCheck(ctx, elog, rlog)
	}
}

func (rp *recurseProc) expressionCheck(ctx *expressionCheckContext, lin *linearContext, log util.ErrorLog) *linearContext {
	tdef := rp.channel.expressionCheck(ctx, log)
	var ok bool
	if tdef != nil {
		rp.btype, ok = tdef.getType().(*brokerType)
		if !ok || rp.btype.requester != nil {
			rp.reportErrorf(log, "expected recursion broker type; instead found: %s", tdef.String())
		} else {
			rp.channel.setType(rp.btype)
			rp.gconfig, ok = rp.btype.gconfig.getType().(globalConfig)
			if !ok {
				rp.reportErrorf(log, "exppecting global configuration for broker type; instead found %s", rp.btype.gconfig.String())
			}
		}
	}

	participants := util.NewHashSet[string]()
	participants.Add(rp.rec.String())
	for _, pexpr := range rp.participants {
		if participants.Contains(pexpr.String()) {
			pexpr.reportErrorf(log, "duplicate definition of participant: %q.", pexpr.String())
		}
		participants.Add(pexpr.String())
	}

	if !ok || rp.btype.requester != nil {
		return lin.newEndLocal(rp.line)
	}

	gconfigParticipants := util.NewHashSet[string]()
	for _, participant := range rp.gconfig.participants() {
		gconfigParticipants.Add(participant)
	}

	for _, p := range participants.Slice() {
		if !gconfigParticipants.Contains(p) {
			stream := util.NewStream().Inc().Inc()
			rp.gconfig.prettyPrint(stream)
			rp.reportErrorf(log, "broker channel:\n%s\n\tdoes not define participant %q.", stream.String(), p)
		}
	}

	rp.locals = make(map[string]local, len(rp.applications))

	rp.locals[rp.rec.id] = rp.app.expressionCheck(rp.rec, gconfigParticipants, ctx, log)
	for i := range rp.applications {
		rp.locals[rp.participants[i].id] = rp.applications[i].expressionCheck(rp.participants[i], participants, ctx, log)
	}

	return lin.newEndLocal(rp.line)
}

func (rp *recurseProc) sessionCheck(ctx *sessionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {

	for i := range rp.applications {
		rp.participants[i].sessionCheck(ctx, elog, rlog)
		rp.applications[i].sessionCheck(ctx, elog, rlog)
	}

	participants := append(rp.participants, rp.rec)
	for _, pexpr := range participants {
		loc, present := rp.locals[pexpr.id]
		if !present {
			// prior error likely; skip to avoid nil deref
			continue
		}
		locAbstr := rp.gconfig.project(pexpr.pType())
		if !loc.subtypeOf(locAbstr.loc, util.NewHashSet[typePair]()) {
			stream1 := util.NewStream().Inc().Inc()
			stream1.Print(pexpr.String() + ": ")
			loc.prettyPrint(stream1)

			stream2 := util.NewStream().Inc().Inc()
			rp.gconfig.prettyPrint(stream2)

			pexpr.reportErrorf(elog, "role:\n%v\n\tis not projectable from broker channel:\n%v", stream1.String(), stream2.String())
		}
	}
}


// func (rp *recurseProc) evaluate(ctx *evaluationContext) {
// 	v := rp.channel.evaluate(ctx)
// 	broker, ok := v.(*recBroker)
// 	if !ok {
// 		msg := rp.runtimeErrorf("undefined broker: %q.", rp.channel.id)
// 		fmt.Fprintln(os.Stderr, msg)
// 		return
// 	}

// 	chans := make([]brokerChannel, 0, len(rp.participants))
// 	defer func() { for i := range chans { close(chans[i]) } }()

// 	fs := make([]func()(channelMap, bool), 0, len(rp.participants))

// 	for _, pexpr := range rp.participants {
// 		c := make(brokerChannel)
// 		chans = append(chans, c)
// 		fs = append(fs, broker.recurse(pexpr.id, c))
// 	}

// 	c := make(brokerChannel)
// 	rf := broker.recurse(rp.rec.id, c)

// 	var brokerChannels channelMap
// 	closed := false

// 	for i := range chans {
// 		if fs[i] != nil {
// 			brokerChannels, closed = fs[i]()
// 		}
// 	}

// 	if rf == nil {
// 		return
// 	}

// 	brokerChannels, closed = rf()

// 	// copy to concurrent map read and map write runtime error
// 	channels := make(channelMap)
// 	for k1, ichans := range brokerChannels {
// 		channels[k1] = make(map[string]util.Channel[expression])
// 		for k2, ichan := range ichans {
// 			channels[k1][k2] = ichan
// 		}
// 	}

// 	for i, p1 := range rp.participants {
// 		for j, p2 := range rp.participants {
// 			if i == j { continue }
// 			channels[p1.id][p2.id] = util.NewBasicChannel[expression](true)
// 		}
// 		channels[p1.id][rp.rec.id] = util.NewBasicChannel[expression](true)
// 		channels[rp.rec.id][p1.id] = util.NewBasicChannel[expression](true)
// 	}

// 	if !closed {
// 		rp.app.evaluate(rp.rec, channels, ctx)
// 		for i := range rp.applications {
// 			rp.applications[i].evaluate(rp.participants[i], channels, ctx)
// 		}
// 	}
// }

func (rp *recurseProc) evaluate(ctx *evaluationContext) {
	v := rp.channel.evaluate(ctx)
	broker, ok := v.(recBrokerI)
	if !ok {
		msg := rp.runtimeErrorf("undefined broker: %q.", rp.channel.id)
		fmt.Fprintln(os.Stderr, msg)
		return
	}

	for _, pexpr := range rp.participants {
		broker.recurse(pexpr.id)
	}

	// brokerChannels, running := broker.recurse(rp.rec.id)
	channels, running := broker.recurse(rp.rec.id)

	if !running {
		return
	}

	// copy to concurrent map read and map write runtime error
	// channels := make(channelMap)
	// for k1, ichans := range brokerChannels {
	// 	channels[k1] = make(map[string]util.Channel[expression])
	// 	for k2, ichan := range ichans {
	// 		channels[k1][k2] = ichan
	// 	}
	// }

	// for i, p1 := range rp.participants {
	// 	for j, p2 := range rp.participants {
	// 		if i == j { continue }
	// 		channels[p1.id][p2.id] = util.NewBasicChannel[expression](true)
	// 	}
	// 	channels[p1.id][rp.rec.id] = util.NewBasicChannel[expression](true)
	// 	channels[rp.rec.id][p1.id] = util.NewBasicChannel[expression](true)
	// }

	rp.app.evaluate(rp.rec, channels, ctx)
	for i := range rp.applications {
		rp.applications[i].evaluate(rp.participants[i], channels, ctx)
	}
}


func (rp *recurseProc) prettyPrint(iw util.IndentedWriter) {
	iw.Print("rec")

	iw.Print("(" + rp.channel.id + ") role ")
	rp.rec.prettyPrint(iw)
	iw.Print(" as ")
	rp.app.prettyPrint(iw)
	if len(rp.applications) != 0 {
		iw.Println(" with ")
		iw.Inc()
		for i := range rp.applications {
			iw.Print("role ")
			rp.participants[i].prettyPrint(iw)
			iw.Print(" as ")
			rp.applications[i].prettyPrint(iw)
			iw.Println(";")
		}
		iw.Dec()
	}
	// iw.Println()
}

func (rp *recurseProc) String() string {
	s := "rec (" + rp.channel.id + ")"
	s += " role " + rp.rec.String() + " as " + rp.app.String()
	if len(rp.applications) != 0 {
		s += " with "
		for i := range rp.applications {
			s += "role " + rp.participants[i].String() + " as " + rp.applications[i].String() + "; "
		}
	}
	// s += " }"
	return s
}

func (_ *recurseProc) goCode(_ util.IndentedWriter) {}
