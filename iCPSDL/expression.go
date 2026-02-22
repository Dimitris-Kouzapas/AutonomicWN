package iCPSDL

import "fmt"
import "iCPSDL/util"

/******************************************************************************
 * expression
 *****************************************************************************/

type expression interface {
    setFilename(string)
    reportError(string, util.ErrorLog)
    execute(*iCPSDL) (expression, util.ErrorLog)
    prettyPrint(*util.Stream)
    getType() string
    getModule() string
    info() string
}

type expressionImpl struct {
    baseNode
    expr expression
    executed expression
}

func (this *expressionImpl) newExpression(expr expression, line int) {
    this.expr = expr
    this.init(line)
}

func (this *expressionImpl) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    if this.expr != nil {
        this.expr.setFilename(filename)
    }
}

func (this *expressionImpl) execute(context *iCPSDL) (expression, util.ErrorLog) {
    var log util.ErrorLog
    this.executed, log = this.expr.execute(context)
    return this.executed, log
}

func (this *expressionImpl) prettyPrint(stream *util.Stream) {
    this.expr.prettyPrint(stream)
}

func (this *expressionImpl) reportCommandLineError(err string, log util.ErrorLog) {
    log.Add(util.NewCommandLineError(err, this.getFilename(), this.getLine()))
}

func (this *expressionImpl) getType() string {
    return this.executed.getType()
}

func (this *expressionImpl) getModule() string {
    return this.expr.getModule()
}

func (this *expressionImpl) info() string {
    return this.executed.info()
}

/*****************************************************************************
 * load
 *****************************************************************************/

type loadExpression struct {
    baseNode
    id string
}

func newLoad(id string, line int) (this *loadExpression) {
    this = new(loadExpression)
    this.id = id
    this.init(line)
    return
}

func (this *loadExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    filename := this.id + context.suffix
    log := newCustomErrorListener(filename)
    expressions := parseFile(filename, log)
    if log.HasErrors() {
        return nil, log
    }
    for _, stmt := range expressions {
        stmt.setFilename(filename)
    }
    return context.executeStack(expressions)
}

func (this *loadExpression) prettyPrint(stream *util.Stream) {
    stream.Print("load " + this.id)
}

func (this *loadExpression) getType() string {
    return "File"
}

func (this *loadExpression) getModule() string {
    //".dss" should be a parameter
    return this.id + ".dss"
}

func (this *loadExpression) info() string {
    return "load"
}

/*****************************************************************************
 * show
 *****************************************************************************/

type showExpression struct {
    expressionImpl
}

func newShow(expr expression, line int) (this *showExpression) {
    this = new(showExpression)
    this.newExpression(expr, line)
    return
}

func (this *showExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    if this.expr == nil {
        for id, expr := range context.expressions {
            fmt.Println("   " + id + ": " + expr.getType())
        }
        return this, util.NewErrorLog()
    }
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() == true {
        return nil, log
    }
    stream := util.NewStream()
    executed.prettyPrint(stream)
    fmt.Println( stream.ToString() )

    return executed, log
}

func (this *showExpression) prettyPrint(stream *util.Stream) {
    stream.Print("show")
    if this.expr == nil {
        return
    }
    stream.Print(" ")
    this.expressionImpl.prettyPrint(stream)
}

/*****************************************************************************
 * assignment
 *****************************************************************************/

type assignmentExpression struct {
    expressionImpl
    name string
}

func newAssignment(name string, expr expression, line int) (this *assignmentExpression) {
    this = new(assignmentExpression)
    this.name = name
    this.newExpression(expr, line)
    return
}

func (this *assignmentExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    context.expressions[this.name] = executed
    return executed, log
}

func (this *assignmentExpression) prettyPrint(stream *util.Stream) {
    stream.Print(this.name + ":= ")
    this.expressionImpl.prettyPrint(stream)
}

func (this *assignmentExpression) getModule() string {
    return this.name
}

/*****************************************************************************
 * remove
 *****************************************************************************/

type removeExpression struct {
    expressionImpl
    name string
}

func newRemove(name string, line int) (this *removeExpression) {
    this = new(removeExpression)
    this.name = name
    this.newExpression(nil, line)
    return
}

func (this *removeExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    _, ok := context.expressions[this.name]
    log := util.NewErrorLog()
    if ok == false {
        error := fmt.Sprintf("Module %v does not exists", this.name)
        log.Add(util.NewExecutionException(error))
        return nil, log
    }
    this.expr = context.expressions[this.name]
    this.executed = this.expr
    delete(context.expressions, this.name)
    return this.executed, log
}

func (this *removeExpression) prettyPrint(stream *util.Stream) {
    stream.Print("rm " + this.name)
}

/*****************************************************************************
 * translate
 *****************************************************************************/

type translateExpression struct {
    expressionImpl
}

func newTranslate(expr expression, line int) (this *translateExpression) {
    this = new(translateExpression)
    this.newExpression(expr, line)
    return
}

func (this *translateExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    exec, ok := executed.(*process)
    if ok == false {
        error := fmt.Sprintf("%v %v is not a process.", util.CapitaliseFirst(executed.getType()), this.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }
    this.executed = exec.translate()
    return this.executed, log
}

func (this *translateExpression) prettyPrint(stream *util.Stream) {
    stream.Print("translate ")
    this.expressionImpl.prettyPrint(stream)
}

/*****************************************************************************
 * traverse
 *****************************************************************************/

type traverseExpression struct {
    expressionImpl
    node string
}

func newTraverse(node string, expr expression, line int) (this *traverseExpression) {
    this = new(traverseExpression)
    this.node = node
    this.newExpression(expr, line)
    return
}

func (this *traverseExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    arg, ok := executed.(*arGraph)
    if ok == false {
        error := fmt.Sprintf("%v %v is not a state estimation graph.", util.CapitaliseFirst(executed.getType()), this.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }

    this.executed = arg.traverse(this.node, log)
    if log.HasErrors() {
        return nil, log
    }
    return this.executed, log
}

func (this *traverseExpression) prettyPrint(stream *util.Stream) {
    stream.Print("traverse ")
    this.expressionImpl.prettyPrint(stream)
}

/*****************************************************************************
 * configure
 *****************************************************************************/

type configureExpression struct {
    expressionImpl
    rep expression
    controller string
    actuator string
}

func newConfigure(expr expression, rep expression, controller string, actuator string, line int) (this *configureExpression) {
  this = new(configureExpression)
  this.newExpression(expr, line)
  this.rep = rep
  this.controller = controller
  this.actuator = actuator
  return
}

func (this *configureExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    art, ok := executed.(*arTree)
    if ok == false {
        error := fmt.Sprintf("%v %v is not a state estimation redundancy tree.", util.CapitaliseFirst(executed.getType()), this.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }
    executed, log = this.rep.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    rep, ok := executed.(*repository)
    if ok == false {
        error := fmt.Sprintf("%v %v is not a agent repository.", util.CapitaliseFirst(this.rep.getType()), this.rep.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }

    this.executed = art.configure(rep, this.controller, this.actuator, log)
    return this.executed, log
}

func (this *configureExpression) prettyPrint(stream *util.Stream) {
    stream.Print("configure ")
    this.expressionImpl.prettyPrint(stream)
    stream.Print(" ")
    this.rep.prettyPrint(stream)
}

/*****************************************************************************
 * compose
 *****************************************************************************/

type composeExpression struct {
    expressionImpl
}

func newCompose(expr expression, line int) (this *composeExpression) {
    this = new(composeExpression)
    this.newExpression(expr, line)
    return
}

func (this *composeExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    lc, ok := executed.(*local)
    if ok == false {
        error := fmt.Sprintf("%v %v is not a local configuration.", util.CapitaliseFirst(executed.getType()), this.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }

    this.executed = lc.compose(log)
    return this.executed, log
}

func (this *composeExpression) prettyPrint(stream *util.Stream) {
    stream.Print("compose ")
    this.expressionImpl.prettyPrint(stream)
}

/*****************************************************************************
 * project
 *****************************************************************************/

type projectExpression struct {
    expressionImpl
}

func newProject(expr expression, line int) (this *projectExpression) {
    this = new(projectExpression)
    this.newExpression(expr, line)
    return
}

func (this *projectExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    glob, ok := executed.(*global)
    if ok == false {
        error := fmt.Sprintf("%v %v is not a global configuration.", util.CapitaliseFirst(executed.getType()), this.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }

    this.executed = glob.project(log)
    return this.executed, log
}

func (this *projectExpression) prettyPrint(stream *util.Stream) {
    stream.Print("project ")
    this.expressionImpl.prettyPrint(stream)
}

/*****************************************************************************
 * project
 *****************************************************************************/

type infoExpression struct {
    expressionImpl
}

func newInfo(expr expression, line int) (this *infoExpression) {
    this = new(infoExpression)
    this.newExpression(expr, line)
    return
}

func (this *infoExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    if this.expr == nil {
        log := util.NewErrorLog()
        output, err := util.ReadFile(context.infoFile)
        if err != nil {
            error := fmt.Sprintf("%v", err)
            log.Add(util.NewExecutionException(error))
            return nil, log
        }
        fmt.Println( output )
        return nil, log
    }
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    fmt.Println( executed.info() )
    return this.executed, log
}

func (this *infoExpression) prettyPrint(stream *util.Stream) {
    stream.Print("info ")
    this.expressionImpl.prettyPrint(stream)
}

/*****************************************************************************
 * id expression
 *****************************************************************************/

type idExpression struct {
    expressionImpl
    id string
}

func newIdExpression(id string, line int) (this *idExpression) {
    this = new(idExpression)
    this.expr = nil
    this.id = id
    this.init(line)
    return
}

func (this *idExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    ok := false
    this.expr, ok = context.expressions[this.id]
    log := util.NewErrorLog()
    if ok == false {
        error := fmt.Sprintf("Module %v not found.", this.id)
        this.reportCommandLineError(error, log)
        return nil, log
    }
    this.executed = this.expr
    return this.expr, log
}

func (this *idExpression) prettyPrint(stream *util.Stream) {
    stream.Print(this.id)
}

func (this *idExpression) getModule() string {
    return this.id
}

/*****************************************************************************
 * access expression
 *****************************************************************************/

type accessExpression struct {
    expressionImpl
    index int
}

func newAccess(expr expression, index int, line int) (this *accessExpression) {
    this = new(accessExpression)
    this.newExpression(expr, line)
    this.index = index
    return
}

func (this *accessExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    executed, log := this.expressionImpl.execute(context)
    if log.HasErrors() {
        return nil, log
    }
    forest, ok := executed.(*forest)
    if ok == false {
        error := fmt.Sprintf("%v %v is not accessible.", util.CapitaliseFirst(this.getType()), this.getModule())
        this.reportCommandLineError(error, log)
        return nil, log
    }
    if len(forest.trees) <= this.index {
        error := fmt.Sprintf("Index out of bounds. State estimation forest %v has %v trees.", this.getModule()/*this.id*/, len(forest.trees))
        this.reportCommandLineError(error, log)
        return nil, log
    }
    this.executed = forest.trees[this.index]
    return this.executed, log
}

func (this *accessExpression) getModule() string {
    return fmt.Sprintf("%v_%v", this.expressionImpl.getModule(), this.index)
}

func (this *accessExpression) prettyPrint(stream *util.Stream) {
    this.expressionImpl.prettyPrint(stream)
    s := fmt.Sprintf("[%v]", this.index)
    stream.Print(s)
}

/*****************************************************************************
 * exit
 *****************************************************************************/

type exitExpression struct {
    baseNode
}

func newExit(line int) (this *exitExpression) {
    this = new(exitExpression)
    this.init(line)
    return
}

func (this *exitExpression) execute(context *iCPSDL) (expression, util.ErrorLog) {
    context.quit = true
    return this, util.NewErrorLog()
}

func (this *exitExpression) prettyPrint(stream *util.Stream) {
    stream.Print("exit")
}

func (this *exitExpression) getType() string {
    return "exit"
}

func (this *exitExpression) getModule() string {
    return "exit"
}

func (this *exitExpression) info() string {
    return "exit"
}
