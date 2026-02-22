package iCPSDL

import "fmt"
import "iCPSDL/util"

/***********************
 * state estimation graph node
 ***********************/

type isegNode interface {
    clone() isegNode
    getParent() isegNode
    getChild() isegNode
    addNeighbour(isegNode)
    addReverse(isegNode)
    configure(*repository, *local, string, util.ErrorLog) string
    prettyPrint(stream *util.Stream)
    mermaid(stream *util.Stream)
    toString() string
}

type argNode struct {
    component icomponent
    neighbours map[string]isegNode
    reverse map[string]isegNode
}

func (this *argNode) newArgNode(component icomponent) {
    this.component = component
    this.neighbours = make(map[string]isegNode)
    this.reverse = make(map[string]isegNode)
}

func (this *argNode) getParent() isegNode {
    for _, next := range this.neighbours {
        return next
    }
    return nil
}

func (this *argNode) getChild() isegNode {
    for _, next := range this.reverse {
        return next
    }
    return nil
}

func (this *argNode) prettyPrint(stream *util.Stream) {
    stream.Print(this.component.getName())
}

func (this *argNode) toString() string {
    return this.component.getName()
}

/***********************
 * arg sensor node
 ***********************/

type sensorNode struct {
    argNode
}

func newSensorNode(component icomponent) (this *sensorNode) {
    this = new(sensorNode)
    this.newArgNode(component)
    return
}

func (this *sensorNode) clone() isegNode {
    return newSensorNode(this.component)
}

func (this *sensorNode) addNeighbour(n isegNode) {
    this.neighbours[n.toString()] = n
    n.addReverse(this)
}

func (this *sensorNode) addReverse(n isegNode) {
    this.reverse[n.toString()] = n
}

func (this *sensorNode) configure(rep *repository, lc *local, controller string, log util.ErrorLog) string {
    name := this.toString()
    comp := this.component.(*sensorComp)
    class := comp.property.getName()
    sess := rep.getSensorClass(class)
    if sess == nil {
        error := fmt.Sprintf("Sensor class %v not found.", class)
        log.Add(util.NewExecutionException(error))
        return ""
    }
    grandParent := this.getParent().getParent()
    adv := controller
    if grandParent != nil {
        adv = grandParent.toString()
    }
    substitutionMap := make(map[string]string)
    substitutionMap["consumer1"] = adv
    sess = sess.rename(substitutionMap)
    lc.add(name, sess)
    return name
}

func (this *sensorNode) mermaid(stream *util.Stream) {
    this.prettyPrint(stream)
    stream.Print("((" + this.toString() + "))")
}

/***********************
 * state estimation graph component node
 ***********************/

type componentNode struct {
    argNode
    model imodel
}

func newComponentNode(component icomponent, model imodel) (this *componentNode) {
    this = new(componentNode)
    this.newArgNode(component)
    this.model = model
    return
}

func (this *componentNode) clone() isegNode {
    return newComponentNode(this.component, this.model)
}

func (this *componentNode) addNeighbour(n isegNode) {
    this.neighbours[n.toString()] = n
    n.addReverse(this)
}

func (this *componentNode) addReverse(n isegNode) {
    this.reverse[n.toString()] = n
}

func (this *componentNode) configure(rep *repository, lc *local, controller string, log util.ErrorLog) string {
    _, ok := this.model.(*property)
    if ok == true {
        child := this.getChild()
        name := ""
        if child != nil {
            name = child.configure(rep, lc, controller, log)
        }
        return name
    }
    name := this.toString()
    class := this.model.getName()
    sess := rep.getEstimatorClass(class)
    if sess == nil {
        error := fmt.Sprintf("Estimator class %v not found.", class)
        log.Add(util.NewExecutionException(error))
        return ""
    }
    i := 1
    substitutionMap := make(map[string]string)

    for _, child := range this.reverse {
        child.configure(rep, lc, controller, log)
        if log.HasErrors() {
            return ""
        }
        grandChild := child.getChild()
        if grandChild != nil {
            producer := fmt.Sprintf("producer%v", i)
            substitutionMap[producer] = grandChild.toString()
            i += 1
        }
    }
    grandParent := this.getParent().getParent()
    adv := controller
    if grandParent != nil {
        adv = grandParent.toString()
    }
    substitutionMap["consumer1"] = adv
    sess = sess.rename(substitutionMap)
    lc.add(name, sess)
    return name
}

func (this *componentNode) prettyPrint(stream *util.Stream) {
    this.argNode.prettyPrint(stream)
    stream.Print("." + this.model.getName())
}

func (this *componentNode) toString() string {
  return this.argNode.toString() + "." + this.model.getName()
}

func (this *componentNode) mermaid(stream *util.Stream) {
  this.prettyPrint(stream)
  stream.Print(this.model.mermaid(this.toString()))
}

/******************************************************************************
 * state estimation graph edge
 ******************************************************************************/

type argEdge struct {
    node1, node2 isegNode
}

func newArgEdge(node1 isegNode, node2 isegNode) (this *argEdge) {
    this = new(argEdge)
    this.node1 = node1
    this.node2 = node2
    return
}

func (this *argEdge) getName() string {
    return this.toString()
}

func (this *argEdge) prettyPrint(stream *util.Stream) {
    this.node1.prettyPrint(stream)
    stream.Print("->")
    this.node2.prettyPrint(stream)
}

func (this *argEdge) toString() string {
    return  this.node1.toString() + "->" + this.node2.toString()
}

func (this *argEdge) mermaid(stream *util.Stream) {
    this.node1.mermaid(stream)
    stream.Print("--->")
    this.node2.mermaid(stream)
}

/******************************************************************************
 * state estimation graph
 ******************************************************************************/

type arGraph struct {
    baseNode
    proc *process
    nodes map[string]isegNode
    edges map[string]*argEdge
}

func newArGraph(proc *process) (this *arGraph) {
    this = new(arGraph)
    this.proc = proc
    this.nodes = make(map[string]isegNode)
    this.edges = make(map[string]*argEdge)
    return
}

func (this *arGraph) addNode(node isegNode) {
    _, ok := this.nodes[node.toString()]
    if ok == false {
        this.nodes[node.toString()] = node
    }
}

func (this *arGraph) addEdge(n1 isegNode, n2 isegNode) {
    edge := newArgEdge(n1, n2)
    this.edges[edge.toString()] = edge
    node1 := this.nodes[n1.toString()]
    node2 := this.nodes[n2.toString()]
    node1.addNeighbour(node2)
}

func (this *arGraph) execute(context *iCPSDL) (expression, util.ErrorLog) {
    return this, util.NewErrorLog()
}

func (this *arGraph) traverse(rootName string, log util.ErrorLog) *forest {
    node, ok := this.nodes[rootName]
    if ok == false {
      error := fmt.Sprintf("Graph does not have node %v.", rootName)
      log.Add(util.NewExecutionException(error))
      return nil
    }
    comp, ok2 := node.(*componentNode)
    if ok2 == false {
        error := fmt.Sprintf("Node %v is a sensor node.", rootName)
        log.Add(util.NewExecutionException(error))
        return nil
    }
    _, ok3 := comp.model.(*property)
    if ok3 == false {
      error := fmt.Sprintf("Node %v is not a property node.", rootName)
      log.Add(util.NewExecutionException(error))
      return nil
    }

    tree := newArTree(comp, this.proc)
    return this.traverseProperty(tree, rootName, log)
}

func (this *arGraph) traverseProperty(tree *arTree, rootName string, log util.ErrorLog) *forest {
    forest := newForest()
    r := this.nodes[rootName]
    switch root := r.(type) {
        case *sensorNode:
            error := fmt.Sprintf("Expecting a property node, but found sensor node: %v.", root.toString())
            log.Add(util.NewExecutionException(error))
            return nil
        case *componentNode:
            switch root.model.(type) {
                case *property:
                    tree.addNode(root)
                    for _, nei := range root.reverse {
                        tmpTree := tree.clone()
                        _, ok := tmpTree.nodes[nei.toString()]
                        if ok == true { continue }
                        tmpTree.addNode(nei)
                        tmpTree.addEdge(nei, root)
                        tmpForest := this.traverseModel(tmpTree, root, nei.toString(), log)
                        if log.HasErrors() { return nil }
                        forest.addForest(tmpForest)
                    }
                    if len(root.reverse) == 0 {
                        forest.addTree(tree)
                    }
                case *enumProperty:
                    error := fmt.Sprintf("Expecting a property node, but found enumeration node: %v.", root.toString())
                    log.Add(util.NewExecutionException(error))
                    return nil
                case *model:
                    error := fmt.Sprintf("Expecting a property node, but found model node: %v.", root.toString())
                    log.Add(util.NewExecutionException(error))
                    return nil
            }
    }
    return forest
}

func (this *arGraph) traverseModel(tree *arTree, c *componentNode, rootName string, log util.ErrorLog) *forest {
    forest := newForest()
    r := this.nodes[rootName]
    switch root := r.(type) {
        case *sensorNode:
            forest.addTree(tree)
        case *componentNode:
            switch root.model.(type) {
                case *property:
                    error := fmt.Sprintf("Expecting a model node, but found property node: %v.", root.toString())
                    log.Add(util.NewExecutionException(error))
                    return nil
                case *enumProperty:
                    error := fmt.Sprintf("Expecting a model node, but found enumeration node: %v.", root.toString())
                    log.Add(util.NewExecutionException(error))
                    return nil
                case *model:
                    forest.addTree(tree)
                    for _, nei := range root.reverse {
                        if nei.toString() == c.toString() { continue }
                        tmpForest := newForest()
                        for _, tmpTree := range forest.trees {
                            _, ok := tmpTree.nodes[nei.toString()]
                            if ok == true { return newForest() }
                            tmpTree.addNode(nei)
                            tmpTree.addEdge(nei, root)
                            tmp := this.traverseProperty(tmpTree, nei.toString(), log)
                            if log.HasErrors() { return nil }
                            tmpForest.addForest(tmp)
                        }
                        forest = tmpForest
                    }
            }
    }
    return forest
}

//func (this *arGraph) render() {}

func (this *arGraph) mermaid(stream *util.Stream) {
    stream.Println("graph TD;")
    stream.Inc()
    for _, e := range this.edges {
        e.mermaid(stream)
        stream.Println(";")
    }
    stream.Dec()

}

func (this *arGraph) prettyPrint(stream *util.Stream) {
    stream.Println("State estimation Graph {")
    stream.Inc()
    stream.Println("Nodes:")
    stream.Inc()
    for _, n := range this.nodes {
        stream.Println(n.toString())
    }
    stream.Dec()
    stream.Println("Edges:")
    stream.Inc()
    for _, e := range this.edges {
        e.prettyPrint(stream)
        stream.Println("")
    }
    stream.Dec()
    stream.Println("}")
    stream.Dec()
}

func (this *arGraph) getType() string {
    return "state estimation graph"
}

func (this *arGraph) getModule() string {
    return "*"
}

func (this *arGraph) info() string {
    s := "Type: state estimation graph"
    s += fmt.Sprintf(", Nodes: %v", len(this.nodes))
    s += fmt.Sprintf(", Edges: %v.", len(this.edges))
    return s
}

/******************************************************************************
 * state estimation tree
 ******************************************************************************/
type arTree struct {
   arGraph
   root isegNode
   graph *arGraph
}

func newArTree(root isegNode, proc *process) (this *arTree) {
   this = new(arTree)
   this.proc = proc
   this.root = root.clone()
   this.nodes = make(map[string]isegNode)
   this.edges = make(map[string]*argEdge)
   this.addNode(this.root)
   return
}

func (this *arTree) clone() *arTree {
    tree := newArTree(this.root, this.proc)
    for _, n := range this.nodes {
        tree.addNode(n.clone())
    }
    for _, e := range this.edges {
       tree.addEdge(e.node1, e.node2)
    }
    return tree
}

func (this *arTree) configure(rep *repository, controller string, actuator string, log util.ErrorLog) *local {
    lc := newLocal(0)

    act, ok := this.proc.actuators[actuator]
    if ok == false {
        error := fmt.Sprintf("Actuator %v not found.", actuator)
        log.Add(util.NewExecutionException(error))
        return nil
    }

    actuatorClass := act.getClass().name
    //TODO: add additional roles, as they arise from roles that are already added
    actSession := rep.getActuatorClass(actuatorClass)
    if actSession == nil {
        error := fmt.Sprintf("Actuator class %v not found.", actuatorClass)
        log.Add(util.NewExecutionException(error))
        return nil
    }
    substitutionMap := make(map[string]string)
    substitutionMap["producer1"] = controller
    actSession = actSession.rename(substitutionMap)
    lc.add(actuator, actSession)

    controlSession := rep.getController(controller)
    if controlSession == nil {
        error := fmt.Sprintf("Controller %v not found.", controller)
        log.Add(util.NewExecutionException(error))
        return nil
    }
    substitutionMap = make(map[string]string)
    substitutionMap["consumer1"] = actuator
    producer := this.root.configure(rep, lc, controller, log)
    if log.HasErrors() {
        return nil
    }
    substitutionMap["producer1"] = producer
    controlSession = controlSession.rename(substitutionMap)
    lc.add(controller, controlSession)
    return lc
}

func (this *arTree) getType() string {
    return "state estimation tree"
}

func (this *arTree) info() string {
    s := "Type: state estimation tree"
    s += ", Root: " + this.root.toString()
    s += fmt.Sprintf(", Nodes: %v", len(this.nodes))
    s += fmt.Sprintf(", Edges: %v.", len(this.edges))
    return s
}

/******************************************************************************
 * Forest
 ******************************************************************************/

type forest struct {
    baseNode
    trees []*arTree
}

func newForest() (this *forest) {
    this = new(forest)
    this.trees = make([]*arTree, 0)
    return
}

func (this *forest) addTree(tree *arTree) {
    this.trees = append(this.trees, tree)
}

func (this *forest) addForest(f *forest) {
    this.trees = append(this.trees, f.trees...)
}

func (this *forest) execute(context *iCPSDL) (expression, util.ErrorLog){
    return this, util.NewErrorLog()
}

func (this *forest) render() {}

func (this *forest) mermaid(stream *util.Stream) {
    stream.Inc()
    for _, tree := range this.trees {
        tree.mermaid(stream)
        stream.Println("")
    }
    stream.Dec()
}

func (this *forest) prettyPrint(stream *util.Stream) {
    stream.Inc()
    for i, tree := range this.trees {
        msg := fmt.Sprintf("%v.", i)
        stream.Print(msg)
        tree.prettyPrint(stream)
        stream.Println("")
    }
    stream.Dec()
}

func (this *forest) getType() string {
    return "state estimation forest"
}

func (this *forest) getModule() string {
    return "*"
}

func (this *forest) info() string {
    s := "Type: state estimation forest"
    s += fmt.Sprintf(", Trees: %v.", len(this.trees))
    return s
}
