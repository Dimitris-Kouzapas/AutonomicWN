package iCPSDL

import "fmt"
import "iCPSDL/util"

/*******************************************************************************
 * process
 ******************************************************************************/

type process struct {
    baseNode
    domainName string
    decls []processDecl
    dmn *domain
    devices map[string]*hardware
    physicals map[string]*physicalComp
    actuators map[string]*actuatorComp
    sensors map[string]*sensorComp
    connections map[string]*connection
    arg *arGraph
}

func newProcess(domainName string, line int) (this *process) {
    this = new(process)
    this.init(line)
    this.domainName = domainName
    this.decls = make([]processDecl, 0)
    this.devices = make(map[string]*hardware)
    this.physicals = make(map[string]*physicalComp)
    this.actuators = make(map[string]*actuatorComp)
    this.sensors = make(map[string]*sensorComp)
    this.connections = make(map[string]*connection)
    this.arg = nil
    return
}

func (this *process) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    for _, decl := range this.decls {
        decl.setFilename(filename)
    }
}

func (this *process) add(d processDecl) {
    this.decls = append(this.decls, d)
}

func (this *process) getComponent(name string) (icomponent, bool) {
    ph, ok1 := this.physicals[name]
    if ok1 == true {
        return ph, ok1
    }
    ac, ok2 := this.actuators[name]
    if ok2 == true {
        return ac, ok2
    }
    return nil, false
}

func (this *process) getComponentOrSensor(name string) (icomponent, bool) {
    c, ok1 := this.getComponent(name)
    if ok1 == true {
        return c, ok1
    }
    s, ok2 := this.sensors[name]
    if ok2 == true {
        return s, ok2
    }
    return nil, false
}

func (this *process) typeCheck(context *iCPSDL, log util.ErrorLog) {
    dmn, ok := context.expressions[this.domainName]
    if ok == false {
        error := fmt.Sprintf("Domain %v is not defined", this.domainName)
        this.reportError(error, log)
        return
    }
    this.dmn, ok = dmn.(*domain)
    if ok == false {
        error := fmt.Sprintf("Expression %v is not a domain", this.domainName)
        this.reportError(error, log)
        return
    }

    size := len(this.decls)
    for i := 0; i < size; i++ {
        for j := i + 1; j < size; j++ {
            if this.decls[i].getName() == this.decls[j].getName() {
                error := fmt.Sprintf(   "Duplicate definition of %v. Originally defined in %v:%v ",
                                        this.decls[j].getName(), this.decls[i].getFilename(), this.decls[i].getLine())
                this.decls[j].reportError(error, log)
            }
        }
    }

    for _, decl := range this.decls {
        switch d := decl.(type) {
            case *hardware:
                this.devices[decl.getName()] = d
            case *physicalComp:
                this.physicals[decl.getName()] = d
            case *actuatorComp:
                this.actuators[decl.getName()] = d
            case *sensorComp:
                this.sensors[decl.getName()] = d
            case *connection:
                this.connections[d.getName()] = d
      }
  }

  for _, decl := range this.decls {
      decl.typeCheck(this, log)
  }
  return
}

func (this *process) translate() *arGraph {
    if this.arg != nil {
        return this.arg
    }

    this.arg = newArGraph(this)

    for _, phy := range this.physicals {
        for _, model := range phy.class.nodeMap {
            node := newComponentNode(phy, model)
            this.arg.addNode(node)
        }
    }
    for _, act := range this.actuators {
        for _, model := range act.class.nodeMap {
            node := newComponentNode(act, model)
            this.arg.addNode(node)
        }
    }

    for _, sensor := range this.sensors {
        node := newSensorNode(sensor)
        this.arg.addNode(node)
    }

    for _, conn := range this.connections {
        switch snd := conn.component2.(type) {
            case *sensorComp:
                e1 := newSensorNode(snd)
                e2 := newComponentNode(conn.component1, snd.property)
                this.arg.addEdge(e1, e2)
            default:
                for _, ie := range conn.component1.getClass().internalEdges {
                    e1 := newComponentNode(conn.component1, ie.model1)
                    e2 := newComponentNode(conn.component1, ie.model2)
                    this.arg.addEdge(e1, e2)
                }
                for _, ie := range conn.component2.getClass().internalEdges {
                    e1 := newComponentNode(conn.component2, ie.model1)
                    e2 := newComponentNode(conn.component2, ie.model2)
                    this.arg.addEdge(e1, e2)
                }
                //todo: if no trans is found then error
                for _, trans := range this.dmn.translations {
                    if conn.component1.getClass().name == trans.node1  && conn.component2.getClass().name == trans.node2 {
                        for _, ee := range trans.targets {
                            var c1, c2 icomponent
                            if conn.component1.getClass().name == ee.c1 {
                                c1 = conn.component1
                                c2 = conn.component2
                            } else {
                                c1 = conn.component2
                                c2 = conn.component1
                            }
                            e1 := newComponentNode(c1, ee.model1)
                            e2 := newComponentNode(c2, ee.model2)
                            this.arg.addEdge(e1, e2)
                        }
                        break
                    }
                }
        }
    }
    return this.arg
}

func (this *process) execute(context *iCPSDL) (expression, util.ErrorLog) {
    _, ok := context.expressions[this.domainName]
    if ok == false {
        expr, log := newLoad(this.domainName, 0).execute(context)
        if log.HasErrors() {
            return nil, log
        }
        context.expressions[this.domainName] = expr
    }

    var log util.ErrorLog = util.NewErrorLog()

    this.typeCheck(context, log)
    return this, log
}

func (this *process) render() {}

func (this *process) prettyPrint(stream *util.Stream) {
    stream.Println("process " + this.domainName + " {" )
    stream.Inc()
    for _, v := range this.decls {
        _, ok := v.(*connection)
        if ok == true { stream.Print("conn ") }
        v.prettyPrint(stream)
        stream.Println("")
    }
    stream.Dec()
    stream.Println("}")
}

func (this *process) mermaid(stream *util.Stream) {
    stream.Println("graph TD;")
    stream.Inc()
    for _, conn := range this.connections {
        conn.mermaid(stream)
        stream.Println(";")
    }
    stream.Dec()
}

func (this *process) getType() string {
    return "process"
}

func (this *process) getModule() string {
    return "*"
}

func (this *process) info() string {
    s := "Type: process"
    s += fmt.Sprintf(", Domain: %v", this.domainName)
    s += fmt.Sprintf(", Devices: %v", len(this.devices))
    s += fmt.Sprintf(", Physical components: %v", len(this.physicals))
    s += fmt.Sprintf(", Actuators: %v", len(this.actuators))
    s += fmt.Sprintf(", Sensors: %v", len(this.sensors))
    s += fmt.Sprintf(", Connections: %v.", len(this.connections))
    return s
}

/*******************
 * processDecl
 *******************/

type processDecl interface {
    getName() string
    setFilename(string)
    getFilename() string
    getLine() int
    reportError(string, util.ErrorLog)
    typeCheck(*process, util.ErrorLog)
    prettyPrint(*util.Stream)
}

/*******************
 * node
 *******************/

type hardware struct {
    baseNode
    name string
}

func newHardware(name string, line int) (this *hardware) {
    this = new(hardware)
    this.init(line)
    this.name = name
    return
}

func (this *hardware) getName() string {
    return this.name
}

func (this *hardware) setFilename(filename string) {
    this.baseNode.setFilename(filename)
}

func (this *hardware) typeCheck(context *process, log util.ErrorLog) {
}

func (this *hardware) prettyPrint(stream *util.Stream) {
    stream.Print("device " + this.name)
}

func (this *hardware) toString() (s string) {
    s = "device " + this.name
    return
}

/***********************
 * component
 ***********************/

type icomponent interface {
    getName() string
    getClass() *class
    prettyPrint(*util.Stream)
    toString() string
}

type component struct {
    baseNode
    name string
    device string
    className string
    class *class
}

func (this *component) newComponent(name string, device string, className string, line int) {
    this.name = name
    this.device = device
    this.className = className
    this.init(line)
    return
}

func (this *component) getName() string {
    return this.name
}

func (this *component) getClass() *class {
    return this.class
}

func (this *component) typeCheck(context *process, log util.ErrorLog) {
    if this.device != "" {
      _, ok := context.devices[this.device]
      if ok == false {
          error := fmt.Sprintf("Device %v is undefined", this.device)
          this.reportError(error, log)
      }
    }
}

func (this *component) prettyPrint(stream *util.Stream) {
    device := ""
    if this.device != "" {
        device = "@" + this.device
    }
    stream.Print(" " + this.name + device + " " + this.className)
}

func (this *component) toString() string {
    device := ""
    if device != "" {
        device = "@" + this.device
    }
    return " " + this.name + device + " " + this.className
}

/***********************
 * physical component
 ***********************/

type physicalComp struct {
    component
}

func newPhysicalComp(name string, class string, line int) (this *physicalComp) {
    this = new(physicalComp)
    this.component.newComponent(name, "", class, line)
    return
}

func (this *physicalComp) typeCheck(context *process, log util.ErrorLog) {
    phy, ok := context.dmn.physicals[this.className]
    if ok == false {
        error := fmt.Sprintf("Name %v is not defined as a physical component class", this.className)
        this.reportError(error, log)
        return
    }
    this.class = &phy.class
    this.component.typeCheck(context, log)
}

func (this *physicalComp) prettyPrint(stream *util.Stream) {
    stream.Print("physical")
    this.component.prettyPrint(stream)
}

func (this *physicalComp) toString() string {
    return "physical" + this.component.toString()
}

/***********************
 * actuator component
 ***********************/

type actuatorComp struct {
    component
}

func newActuatorComp(name string, device string, class string, line int) (this *actuatorComp) {
    this = new(actuatorComp)
    this.component.newComponent(name, device, class, line)
    return
}

func (this *actuatorComp) typeCheck(context *process, log util.ErrorLog) {
    act, ok := context.dmn.actuators[this.className]
    if ok == false {
        error := fmt.Sprintf("Name %v is not defined as an actuator component class", this.className)
        this.reportError(error, log)
        return
    }
    this.class = &act.class
    this.component.typeCheck(context, log)
}

func (this *actuatorComp) prettyPrint(stream *util.Stream) {
    stream.Print("actuator")
    this.component.prettyPrint(stream)
}

func (this *actuatorComp) toString() string {
    return "actuator" + this.component.toString()
}

/***********************
 * actuator component
 ***********************/

type sensorComp struct {
    component
    propertyName string
    property imodel
}

func newSensorComp(name string, device string, propertyName string, line int) (this *sensorComp) {
    this = new(sensorComp)
    this.component.newComponent(name, device, "", line)
    this.propertyName = propertyName
    this.property = nil
    return
}

func (this *sensorComp) typeCheck(context *process, log util.ErrorLog) {
    ok := false
    this.property, ok = context.dmn.properties[this.propertyName]
    if ok == false {
      error := fmt.Sprintf("Name %v is not defined as a property", this.className)
      this.reportError(error, log)
      return
    }
    this.component.typeCheck(context, log)
}

func (this *sensorComp) prettyPrint(stream *util.Stream) {
    stream.Print("sensor")
    this.component.prettyPrint(stream)
    stream.Print(" " + this.property.getName())
}

func (this *sensorComp) toString() string {
    return "sensor" + this.component.toString() + " " + this.property.getName()
}

/******************************************************************************
 * connection
 ******************************************************************************/

type connection struct {
    baseNode
    node1, node2 string
    component1, component2 icomponent
}

func newConnection(node1 string, node2 string, line int) (this *connection) {
    this = new(connection)
    this.init(line)
    this.node1 = node1
    this.node2 = node2
    this.component1 = nil
    this.component2 = nil
    return
}

func (this *connection) getName() string {
    return this.node1 + "->" + this.node2
}

func (this *connection) typeCheck(context *process, log util.ErrorLog) {
    ok := false
    this.component1, ok = context.getComponent(this.node1)
    if ok == false {
        error := fmt.Sprintf("Component %v is undefined.", this.node1)
        this.reportError(error, log)
    }
    this.component2, ok = context.getComponentOrSensor(this.node2)
    if ok == false {
        error := fmt.Sprintf("Component %v is undefined.", this.node2)
        this.reportError(error, log)
    }
}

func (this *connection) prettyPrint(stream *util.Stream) {
    stream.Print(this.node1 + "->" + this.node2)
}

func (this *connection) mermaid(stream *util.Stream) {
    stream.Print(this.node1 + "--->" + this.node2)
}

func (this *connection) toString() string {
    return  this.node1 + "->" + this.node2
}
