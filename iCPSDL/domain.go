package iCPSDL

import "fmt"
import "iCPSDL/util"

/*******************************************************************************
 * domain
 ******************************************************************************/

type domain struct {
    baseNode
    decls []domainDecl
    properties map[string]imodel
    models map[string]*model
    physicals map[string]*physical
    actuators map[string]*actuator
    translations []*translation
}

func newDomain(line int) (this *domain) {
    this = new(domain)
    this.init(line)
    this.decls = make([]domainDecl, 0)
    this.properties = make(map[string]imodel)
    this.physicals = make(map[string]*physical)
    this.actuators = make(map[string]*actuator)
    this.models = make(map[string]*model)
    this.translations = make([]*translation, 0)
    return
}

func (this *domain) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    for _, decl := range this.decls {
        decl.setFilename(filename)
    }
}

func (this *domain) add(decl domainDecl) {
    this.decls = append(this.decls, decl)
}

func (this *domain) getClass(name string) (*class, bool) {
    phy, ok1 := this.physicals[name]
    if ok1 == true {
        return &phy.class, true
    }
    act, ok2 := this.actuators[name]
    if ok2 == true {
        return &act.class, true
    }
    return nil, false
}

func (this *domain) typeCheck(log util.ErrorLog) {
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
            case *property:
                this.properties[decl.getName()] = d
            case *enumProperty:
                this.properties[decl.getName()] = d
            case *model:
                this.models[d.getName()] = d
            case *physical:
                this.physicals[decl.getName()] = d
            case *actuator:
                this.actuators[decl.getName()] = d
            case *translation:
                this.translations = append(this.translations, d)
        }
    }

    for _, decl := range this.decls {
        decl.typeCheck(this, log)
    }
    return
}

func (this *domain) execute(context *iCPSDL) (expression, util.ErrorLog) {
    log := util.NewErrorLog()
    this.typeCheck(log)
    return this, log
}

func (this *domain) render() {}

func (this *domain) prettyPrint(stream *util.Stream) {
    stream.Println("domain {")
    stream.Inc()
    for _, v := range this.decls {
        v.prettyPrint(stream)
        stream.Println("")
    }
    stream.Dec()
    stream.Println("}")
}

func (this *domain) getType() string {
    return "domain"
}

func (this *domain) getModule() string {
    return "*"
}

func (this *domain) info() string {
    s := "Type: domain"
    s += fmt.Sprintf(", Property classes: %v", len(this.properties))
    s += fmt.Sprintf(", Model classes: %v", len(this.models))
    s += fmt.Sprintf(", Physical component classes: %v", len(this.physicals))
    s += fmt.Sprintf(", Actuator classes: %v", len(this.actuators))
    s += fmt.Sprintf(", Translations: %v.", len(this.translations))
    return s
}

/*******************
 * domainDecl
 *******************/

type domainDecl interface {
    getName() string
    setFilename(string)
    getFilename() string
    getLine() int
    reportError(string, util.ErrorLog)
    typeCheck(*domain, util.ErrorLog)
    prettyPrint(*util.Stream)
}

/*******************
 * model + property
 *******************/
type imodel interface {
    getName() string
    setFilename(string)
    getFilename() string
    getLine() int
    reportError(string, util.ErrorLog)
    typeCheck(*domain, util.ErrorLog)
    prettyPrint(*util.Stream)
    toString() string
    mermaid(string) string
}

type modelImpl struct {
    baseNode
    name string
    kind string
}

func (this *modelImpl) newModelImpl(kind string, name string, line int) {
    this.init(line)
    this.name = name
    this.kind = kind
    return
}

func (this *modelImpl) getName() string {
    return this.name
}

func (this *modelImpl) typeCheck(dmn *domain, log util.ErrorLog) {
}

func (this *modelImpl) prettyPrint(stream *util.Stream) {
    stream.Print(this.kind + " " + this.name)
}

func (this *modelImpl) toString() (s string) {
    s += this.kind + " " + this.name
    return
}

func (this *modelImpl) mermaid(text string) string {
    return "(" + text + ")"
}

/*******************
 * property
 *******************/

type property struct {
    modelImpl
}

func newProperty(name string, line int) (this *property) {
    this = new(property)
    this.newModelImpl("property", name, line)
    return
}

/*******************
 * enumerationType
 *******************/

type enumProperty struct {
    property
    labels []string
}

func newEnumerationProperty(name string, line int) (this *enumProperty) {
    this = new(enumProperty)
    this.name = name
    this.labels = make([]string, 0)
    this.kind = "property"
    this.init(line)
    return
}

func (this *enumProperty) addLabel(label string) {
    this.labels = append(this.labels, label)
}

func (this *enumProperty) prettyPrint(stream *util.Stream) {
    this.property.prettyPrint(stream)
    stream.Print(" { ")
    for i, label := range this.labels {
        if i != 0 {
            stream.Print(", ")
        }
        stream.Print(label)
    }
    stream.Print(" }")
}

func (this *enumProperty) toString() (s string) {
    s = this.property.toString() + "{ "
    for i, label := range this.labels {
        if i != 1 {
            s += ", "
        }
        s += label
    }
    s += " }"
    return
}

/*******************
 * model
 *******************/

type model struct {
    modelImpl
}

func newModel(name string, line int) (this *model) {
    this = new(model)
    this.newModelImpl("model", name, line)
    return
}

func (this *model) mermaid(text string) string {
    return "[" + text + "]"
}

/*********************
 * class
 *********************/

type class struct {
    baseNode
    name string
    kind string
    nodes []string
    nodeMap map[string]imodel
    internalEdges []*internalEdge
}

func (this *class) newClass(kind string, name string, nodes []string, internalEdges []*internalEdge, line int) {
    this.init(line)
    this.kind = kind
    this.name = name
    this.nodes = nodes
    this.nodeMap = make(map[string]imodel)
    this.internalEdges = internalEdges
    return
}

func (this *class) getName() string {
    return this.name
}

func (this *class) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    for _, ie := range this.internalEdges {
        ie.setFilename(filename)
    }
}

func (this *class) hasNode(n string) bool {
    for _, nn := range this.nodes {
        if n == nn {
          return true
        }
    }
    return false
}

func (this *class) typeCheck(dmn *domain, log util.ErrorLog) {
    for _, n := range this.nodes {
        prop, ok1 := dmn.properties[n]
        mod, ok2 := dmn.models[n]
        if ok1 == true {
            this.nodeMap[n] = prop
        } else if ok2 == true {
            this.nodeMap[n] = mod
        } else if (ok1 == false && ok2 == false) {
            error := fmt.Sprintf("Node %v is neither a property nor a model", n)
            this.reportError(error, log)
        }
    }
    for _, ie := range this.internalEdges {
        ie.typeCheck(this, log)
    }
}

func (this *class) prettyPrint(stream *util.Stream) {
    stream.Print(this.kind + " " + this.name + "(")
    for i, node := range this.nodes {
        if i != 0 {
          stream.Print(", ")
        }
        stream.Print(node)
    }
    stream.Print(") : ")
    for i, ie := range this.internalEdges {
        if i != 0 {
          stream.Print(", ")
        }
        ie.prettyPrint(stream)
    }
}

func (this *class) toString() (s string) {
    s += this.kind + " " + this.name
    return
}

/*********************
 * physical
 *********************/

type physical struct {
    class
}

func newPhysical(name string, nodes []string, internalEdges []*internalEdge, line int) (this *physical) {
    this = new(physical)
    this.newClass("physical", name, nodes, internalEdges, line)
    return
}

/*********************
 * actuator
 *********************/

type actuator struct {
    class
}

func newActuator(name string, nodes []string, internalEdges []*internalEdge, line int) (this *actuator) {
    this = new(actuator)
    this.newClass("actuator", name, nodes, internalEdges, line)
    return
}

/***********************
 * internalEdge
 ***********************/

type internalEdge struct {
    baseNode
    node1, node2 string
    model1, model2 imodel
}

func newInternalEdge(node1 string, node2 string, line int) (this *internalEdge) {
    this = new(internalEdge)
    this.init(line)
    this.node1 = node1
    this.node2 = node2
    this.model1 = nil
    this.model2 = nil
    return
}

func (this *internalEdge) getName() string {
    return this.node1 + "->" + this.node2
}

func (this *internalEdge) typeCheck(c *class, log util.ErrorLog) {
    ok := false
    this.model1, ok = c.nodeMap[this.node1]
    if ok == false {
        error := fmt.Sprintf("Undefined node: %v.", this.node1)
        this.reportError(error, log)
    }
    this.model2, ok = c.nodeMap[this.node2]
    if ok == false {
        error := fmt.Sprintf("Undefined node: %v.", this.node2)
        this.reportError(error, log)
    }

    _, ok1 := this.model1.(*property)
    _, ok2 := this.model1.(*model)
    _, ok3 := this.model2.(*property)
    _, ok4 := this.model2.(*model)

    type1 := ""
    type2 := ""
    if ok1 == true { type1 = "property" } else if ok2 == true { type1 = "model" }
    if ok3 == true { type2 = "property" } else if ok4 == true { type2 = "model" }

    if ((ok1 && ok4) || (ok2 && ok3)) == false {
        error := fmt.Sprintf("Translation edge does not implement a property - model pair: %v is a %v, %v is a %v.", this.node1, type1, this.node2, type2)
        this.reportError(error, log)
    }
    return
}

func (this *internalEdge) prettyPrint(stream *util.Stream) {
    stream.Print(this.node1 + "->" + this.node2)
}

/************************************************
* translation and externalEdge
************************************************/

type externalEdge struct {
    baseNode
    c1, c2 string
    m1, m2 string
    class1, class2 *class
    model1, model2 imodel
}

func newExternalEdge(c1, m1, c2, m2 string, line int) (this *externalEdge) {
    this = new(externalEdge)
    this.init(line)
    this.c1 = c1
    this.c2 = c2
    this.m1 = m1
    this.m2 = m2
    this.class1 = nil
    this.class2 = nil
    this.model1 = nil
    this.model2 = nil
    return
}

func (this *externalEdge) getName() string {
    return this.c1 + "." + this.m1 + "->" + this.c2 + "." + this.m2
}

func (this *externalEdge) typeCheck(t *translation, log util.ErrorLog) {
    if this.c1 == t.node1 {
        this.class1 = t.class1
    } else if this.c1 == t.node2 {
        this.class1 = t.class2
    }

    if this.c2 == t.node1 {
        this.class2 = t.class1
    } else if this.c2 == t.node2 {
        this.class2 = t.class2
    }

    if this.class1 == nil {
        error := fmt.Sprintf("Invalid class: %v.", this.c1)
        this.reportError(error, log)
    } else {
        ok := false
        this.model1, ok = this.class1.nodeMap[this.m1]
        if ok == false {
          error := fmt.Sprintf("Invalid property or model: %v.", this.m1)
          this.reportError(error, log)
        }
    }

    if this.class2 == nil {
        error := fmt.Sprintf("Invalid class: %v.", this.c2)
        this.reportError(error, log)
    } else {
        ok := false
        this.model2, ok = this.class2.nodeMap[this.m2]
        if ok == false {
          error := fmt.Sprintf("Invalid property or model: %v.", this.m2)
          this.reportError(error, log)
        }
    }

    _, ok1 := this.model1.(*property)
    _, ok2 := this.model1.(*model)
    _, ok3 := this.model2.(*property)
    _, ok4 := this.model2.(*model)

    type1 := ""
    type2 := ""
    if ok1 == true { type1 = "property" } else if ok2 == true { type1 = "model" }
    if ok3 == true { type2 = "property" } else if ok4 == true { type2 = "model" }

    if ((ok1 && ok4) || (ok2 && ok3)) == false {
        error := fmt.Sprintf("Translation edge does not implement a property - model pair: %v is a %v, %v is a %v.", this.m1, type1, this.m2, type2)
        this.reportError(error, log)
    }
}

func (this *externalEdge) prettyPrint(stream *util.Stream) {
  stream.Print(this.c1 + "." + this.m1 + "->" + this.c2 + "." + this.m2)
}


/****************
* translation
*****************/

type translation struct {
    baseNode
    node1, node2 string
    class1, class2 *class
    targets []*externalEdge
}

func newTranslation(node1 string, node2 string, targets []*externalEdge, line int) (this *translation) {
    this = new(translation)
    this.node1 = node1
    this.node2 = node2
    this.targets = targets
    this.init(line)
    return
}

func (this *translation) getName() string {
    return this.node1 + "->" + this.node2
}

func (this *translation) setFilename(filename string) {
    this.baseNode.setFilename(filename)
    for _, target := range this.targets {
        target.setFilename(filename)
    }
}

func (this *translation) typeCheck(dmn *domain, log util.ErrorLog) {
    ok := false
    this.class1, ok = dmn.getClass(this.node1)
    if ok == false {
      error := fmt.Sprintf("Class %v is undefined.", this.node1)
      this.reportError(error, log)
    }

    this.class2, ok = dmn.getClass(this.node2)
    if ok == false {
      error := fmt.Sprintf("Class %v is undefined.", this.node2)
      this.reportError(error, log)
    }

    for _, target := range this.targets {
        target.typeCheck(this, log)
    }
}

func (this *translation) prettyPrint(stream *util.Stream) {
    stream.Print("translation " + this.node1 + "->" + this.node2)
    stream.Print(" : ")
    for i, target := range this.targets {
        if i != 0 {
            stream.Print(", ")
        }
        target.prettyPrint(stream)
    }
}
