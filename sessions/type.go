package ast

import (
    "fmt"
    "sessions/util"
)

//TODO: Move it somewhere else?
type typePair struct {
    a typedef
    b typedef
}

type typedef interface {
    setFilename(string)
    reportErrorf(util.ErrorLog, string, ...interface{})
    getType() typedef
    subtypeOf(typedef) bool
    subtypeOf_(typedef, *util.HashSet[typePair]) bool
    join(typedef) (typedef, bool)
    // meet(typedef) (typedef, bool)
    // submatch(typedef, substitutionContext) bool
    // supermatch(typedef, substitutionContext) bool
    // substitute(substitutionContext) typedef
    typeCheck(*typeCheckContext, util.ErrorLog)
    projectionCheck(*projectionCheckContext, util.ErrorLog, util.ReportLog)
    defaultValue() expression
    prettyPrint(util.IndentedWriter)
    fmt.Stringer
}

/******************************************************************************
 * nothing type
 ******************************************************************************/

type nothingType struct {
    baseNode
}

func newNothingType() *nothingType {
    return &nothingType {
        baseNode: baseNode {
            line: 0,
        },
    }
}

func (n *nothingType) getType() typedef { return n }
func (_ *nothingType) typeCheck(_ *typeCheckContext, _ util.ErrorLog) {}
func (_ *nothingType) projectionCheck(_ *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}
func (n *nothingType) subtypeOf(tdef typedef) bool { return n.subtypeOf_(tdef, nil) }
func (_ *nothingType) subtypeOf_(_ typedef, _ *util.HashSet[typePair]) bool { return true }
func (n *nothingType) join(tdef typedef) (typedef, bool) { return tdef, true }
// func (n *nothingType) meet(tdef typedef) (typedef, bool) { return n, true }
func (n *nothingType) defaultValue() expression { return newNothingf(n.line, "nothing expression") }
func (n *nothingType) prettyPrint(iw util.IndentedWriter) { iw.Print(n.String()) }
func (_ *nothingType) String() string { return "nothing" }

// func (this *nothingType) submatch(td typedef, ctx substitutionContext) bool {
//     return true
// }
//
// func (this *nothingType) supermatch(td typedef, ctx substitutionContext) bool {
//     return td.submatch(this, ctx)
// }
//
// func (this *nothingType) substitute(ctx substitutionContext) typedef {
//     return this
// }

/******************************************************************************
 * primitive types
 ******************************************************************************/

/******************
 * primitive type
 ******************/

type primitiveType struct {
    baseNode
    typeString string
}

func (_ *primitiveType) typeCheck(_ *typeCheckContext, _ util.ErrorLog) {}
func (_ *primitiveType) projectionCheck(_ *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}
func (pt *primitiveType) prettyPrint(iw util.IndentedWriter) { iw.Print(pt.String()) }
func (pt *primitiveType) String() string { return pt.typeString }

/******************
 * boolean type
 ******************/

type boolType struct {
    primitiveType
}

func newBoolType(line int) *boolType {
    return &boolType {
        primitiveType: primitiveType {
            typeString: "bool",
            baseNode: baseNode {
                line: line,
            },
        },
    }
}

func (bt *boolType) getType() typedef               { return bt }
func (bt *boolType) subtypeOf(tdef typedef) bool    { return bt.subtypeOf_(tdef, nil) }

func (_ *boolType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch tdef.getType().(type) {
        case *nothingType, *boolType:
            return true
        default:
            return false
    }
}

func (bt *boolType) join(tdef typedef) (typedef, bool) {
    switch tdef.getType().(type) {
        case *nothingType, *boolType:
            return bt, true
        default:
            return tdef, false
    }
}

func (bt *boolType) defaultValue() expression {
    return newFalseExpr(bt.line)
}

// func (bt *boolType) meet(tdef typedef) (typedef, bool) {
//     switch tdef.getType().(type) {
//         case *nothingType, *boolType:
//             return tdef, true
//         default:
//             return bt, false
//     }
// }

// func (this *boolType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return this.subtypeOf(tdef)
// }
//
// func (this *boolType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return td2.subtypeOf(this)
//     }
//     return tdef.subtypeOf(this)
// }

/******************
 * integer type
 ******************/

type intType struct {
    primitiveType
}

func newIntType(line int) *intType {
    return &intType {
        primitiveType: primitiveType {
            typeString: "int",
            baseNode: baseNode {
                line: line,
            },
        },
    }
}

func (it *intType) getType() typedef { return it }
func (it *intType) subtypeOf(tdef typedef) bool { return it.subtypeOf_(tdef, nil) }

func (_ *intType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch tdef.getType().(type) {
        //case *anyTC, *eq, *num, *ord, *intType, *floatType:
        case *nothingType, *intType, *floatType:
            return true
        default:
            return false
    }
}

func (it *intType) join(tdef typedef) (typedef, bool) {
    switch tdef.getType().(type) {
        case *nothingType, *intType:
            return it, true
        case *floatType:
            return tdef, true
        default:
            return tdef, false
    }
}

func (it *intType) defaultValue() expression {
    return newIntExpr("0", it.line)
}

// func (it *intType) meet(tdef typedef) (typedef, bool) {
//     switch tdef.getType().(type) {
//         case *nothingType, *intType:
//             return tdef, true
//         case *floatType:
//             return it, true
//         default:
//             return it, false
//     }
// }

// func (this *intType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return this.subtypeOf(tdef)
// }
//
// func (this *intType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return td2.subtypeOf(this)
//     }
//     return tdef.subtypeOf(this)
// }

/******************
 * float type
 ******************/

type floatType struct {
    primitiveType
}

func newFloatType(line int) *floatType {
    return &floatType {
        primitiveType: primitiveType {
            typeString: "float",
            baseNode: baseNode {
                line: line,
            },
        },
    }
}

func (ft *floatType) getType() typedef            { return ft }
func (ft *floatType) subtypeOf(tdef typedef) bool { return ft.subtypeOf_(tdef, nil) }

func (_ *floatType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch tdef.getType().(type) {
        case *nothingType, *floatType:
            return true
        default:
            return false
    }
}

func (ft *floatType) join(tdef typedef) (typedef, bool) {
    switch tdef.getType().(type) {
        case *nothingType, *intType, *floatType:
            return ft, true
        default:
            return tdef, false
    }
}
 
 func (ft *floatType) defaultValue() expression {
    return newFloatExpr("0", ft.line)
 }

// func (ft *floatType) meet(tdef typedef) (typedef, bool) {
//     switch tdef.getType().(type) {
//         case *nothingType, *intType, *floatType:
//             return tdef, true
//         default:
//             return ft, false
//     }
// }

// func (this *floatType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return this.subtypeOf(tdef)
// }
//
// func (this *floatType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return td2.subtypeOf(this)
//     }
//     return tdef.subtypeOf(this)
// }

/******************
 * string type
 ******************/

type stringType struct {
    primitiveType
}

func newStringType(line int) *stringType {
    return &stringType {
        primitiveType: primitiveType {
            typeString: "string",
            baseNode: baseNode {
                line: line,
            },
        },
    }
}

func (st *stringType) getType() typedef            { return st }
func (st *stringType) subtypeOf(tdef typedef) bool { return st.subtypeOf_(tdef, nil) }

func (_ *stringType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch tdef.getType().(type) {
        case *nothingType, *stringType:
            return true
        default:
            return false
    }
}

func (st *stringType) join(tdef typedef) (typedef, bool) {
    switch tdef.getType().(type) {
        case *nothingType, *stringType:
            return st, true
        default:
            return tdef, false
    }
}

func (st *stringType) defaultValue() expression {
    return newStringExpr("", st.line)
}

// func (st *stringType) meet(tdef typedef) (typedef, bool) {
//     switch tdef.getType().(type) {
//         case *nothingType, *stringType:
//             return tdef, true
//         default:
//             return st, false
//     }
// }

// func (this *stringType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return this.subtypeOf(tdef)
// }
//
// func (this *stringType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return td2.subtypeOf(this)
//     }
//     return tdef.subtypeOf(this)
// }

/******************************************************************************
 * list type
 ******************************************************************************/

type listType struct {
    baseNode
    tdef typedef
}

func newListType(tdef typedef, line int) *listType {
    return &listType {
        baseNode: baseNode {
            line: line,
        },
        tdef: tdef,
    }
}

func (lt *listType) setFilename(filename string) {
    lt.baseNode.setFilename(filename)
    lt.tdef.setFilename(filename)
}

func (lt *listType) getType() typedef { return lt }
func (lt *listType) subtypeOf(tdef typedef) bool { return lt.subtypeOf_(tdef, nil) }

func (lt *listType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch at := tdef.getType().(type) {
        case *nothingType:
            return true
        case *listType:
            return lt.tdef.subtypeOf(at.tdef)
        default:
            return false
    }
}

func (lt *listType) join(tdef typedef) (typedef, bool) {
    switch at := tdef.getType().(type) {
        case *nothingType:
            return lt, true
        case *listType:
            td, ok := lt.tdef.join(at.tdef)
            if ok == true {
                return newListType(td, lt.line), true
            } else {
                return lt, false
            }
        default:
            return lt, false
    }
}

// func (lt *listType) meet(tdef typedef) (typedef, bool) {
//     switch at := tdef.getType().(type) {
//         case *nothingType:
//             return at, true
//         case *listType:
//             td, ok := lt.tdef.meet(at.tdef)
//             if ok == true {
//                 return newListType(td, lt.line), true
//             } else {
//                 return at, false
//             }
//         default:
//             return at, false
//     }
// }

func (lt *listType) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    lt.tdef.typeCheck(ctx, log)
}

func (lt *listType) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    lt.tdef.projectionCheck(ctx, elog, rlog)
}

func (lt *listType) prettyPrint(iw util.IndentedWriter) {
    iw.Print("[")
    lt.tdef.prettyPrint(iw)
    iw.Print("]")
}

func (lt *listType) String() string {
    return "[" + lt.tdef.String() + "]"
}

func (lt *listType) defaultValue() expression {
    return newListExpr(make([]expression, 0), lt.line)
}

// func (this *listType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *listType:
//             return this.tdef.submatch(td.tdef, ctx)
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return false
// }
//
// func (this *listType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *listType:
//             return this.tdef.supermatch(td.tdef, ctx)
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return td2.subtypeOf(this)
//     }
//     return false
// }

/******************************************************************************
 * label type
 ******************************************************************************/

type labelType struct {
    baseNode
    label string
}

func newLabelType(label string, line int) *labelType {
    return &labelType {
        baseNode: baseNode {
            line: line,
        },
        label: label,
    }
}

func (_ *labelType) typeCheck(_ *typeCheckContext, _ util.ErrorLog) {}
func (_ *labelType) projectionCheck(_ *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}
func (lt *labelType) subtypeOf(tdef typedef) bool { return lt.subtypeOf_(tdef, nil) }

func (lt *labelType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch ltype := tdef.getType().(type) {
        case *labelType:
            return lt.label == ltype.label
        default:
            return false
    }
}

func (_ *labelType) join(_ typedef) (typedef, bool)         { return nil, false }
// func (_ *labelType) meet(_ typedef) (typedef, bool)         { return nil, false }

func (lt *labelType) defaultValue() expression              { return newLabelExpr("", lt.line) }

func (lt *labelType) getType() typedef                      { return lt }
func (lt *labelType) prettyPrint(iw util.IndentedWriter)    { iw.Print(lt.String()) }
func (lt *labelType) String() string                        { return lt.label }

/******************************************************************************
 * record type
 ******************************************************************************/

type recordType struct {
    baseNode
    labels []string
    tdefs []typedef
    tdefMap map[string]typedef
}

func newRecordType(labels []string, tdefs []typedef, line int) *recordType {
    tdefMap := make(map[string]typedef)
    for i, label := range labels {
        tdefMap[label] = tdefs[i]
    }
    return &recordType {
        baseNode: baseNode {
            line: line,
        },
        labels: labels,
        tdefs: tdefs,
        tdefMap: tdefMap,
    }
}

func (rt *recordType) setFilename(filename string) {
    rt.baseNode.setFilename(filename)
    for _, tdef := range rt.tdefs {
        tdef.setFilename(filename)
    }
}

func (rt *recordType) subtypeOf(tdef typedef) bool { return rt.subtypeOf_(tdef, nil) }

//TODO: add more operators to recordType to have more typeclasses
func (rt *recordType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch rtype := tdef.getType().(type) {
        case *nothingType:
            return true
        case *recordType:
            for i, lab := range rtype.labels {
                td1 := rtype.tdefs[i]
                td2, ok := rt.tdefMap[lab]
                if !ok || !td2.subtypeOf(td1) {
                    return false
                }
            }
            return true
        default:
            return false
    }
}

func (rt *recordType) join(tdef typedef) (typedef, bool) {
    switch rtype := tdef.getType().(type) {
        case *nothingType:
            return rt, true
        case *recordType:
            labels := make([]string, 0, len(rtype.labels))
            tdefs  := make([]typedef, 0, len(rtype.labels))
            for i, lab := range rtype.labels {
                td1 := rtype.tdefs[i]
                td2, ok := rt.tdefMap[lab]
                if !ok {
                    return rt, false
                }
                max, ok := td2.join(td1)
                if !ok {
                    return rt, false
                }
                labels = append(labels, lab)
                tdefs = append(tdefs, max)
            }
            return newRecordType(labels, tdefs, rt.line), true
        default:
            return tdef, false
    }
}

func (rt *recordType) defaultValue() expression {
    expressions := make([]expression, len(rt.tdefs))
    for i, tdef := range rt.tdefs {
        expressions[i] = tdef.defaultValue()
    }
    return newRecordExpr(rt.labels, expressions, rt.line)
}

// func (rt *recordType) meet(tdef typedef) (typedef, bool) {
//     switch rtype := tdef.getType().(type) {
//         case *nothingType:
//             return rtype, true
//         case *recordType:
//             labels := make([]string, 0, len(rtype.labels))
//             tdefs  := make([]typedef, 0, len(rtype.labels))
//             for i, lab := range rtype.labels {
//                 td1 := rtype.tdefs[i]
//                 td2, ok := rt.tdefMap[lab]
//                 if !ok {
//                     return rt, false
//                 }
//                 min, ok := td2.meet(td1)
//                 if !ok {
//                     return rt, false
//                 }
//                 labels = append(labels, lab)
//                 tdefs = append(tdefs, min)
//             }
//             return newRecordType(labels, tdefs, rt.line), true
//         default:
//             return rt, false
//     }
// }

func (rt *recordType) getType() typedef { return rt }

func (rt *recordType) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    for i := range rt.labels {
        for j := i + 1; j < len(rt.labels); j++ {
            if rt.labels[i] == rt.labels[j] {
                rt.reportErrorf(log, "duplicate definition of record label: %q.", rt.labels[j])
            }
        }
    }
    for _, tdef := range rt.tdefs {
        tdef.typeCheck(ctx, log)
    }
}

func (rt *recordType) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    for _, tdef := range rt.tdefs {
        tdef.projectionCheck(ctx, elog, rlog)
    }
}

func (rt *recordType) prettyPrint(iw util.IndentedWriter) {
    iw.Print("record { ")
    for i, label := range rt.labels {
        iw.Print(label + " ")
        rt.tdefs[i].prettyPrint(iw)
        iw.Print("; ")
    }
    iw.Print(" }")
}

func (rt *recordType) String() string {
    s := "record { "
    for i, label := range rt.labels {
        s += label + " " + rt.tdefs[i].String() + "; "
    }
    return s + " }"
}

// func (this *tupleType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *tupleType:
//             for k := range td.tupleMap {
//                 if this.tupleMap[k] == nil || this.tupleMap[k].submatch(td.tupleMap[k], ctx) == false {
//                     return false
//                 }
//             }
//             return true
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return false
// }

// func (this *tupleType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *tupleType:
//             for k := range this.tupleMap {
//                 if td.tupleMap[k] == nil || this.tupleMap[k].supermatch(td.tupleMap[k], ctx) == false {
//                     return false
//                 }
//             }
//             return true
//       case *typeVar:
//           td2, ok := ctx[td.name]
//           if ok == false {
//               ctx[td.name] = this
//               return true
//           }
//           return td2.subtypeOf(this)
//     }
//     return false
// }
// 
// func (this *tupleType) substitute(ctx substitutionContext) typedef {
//     tdefs := make([]typedef, len(this.tdefs))
//     for i, tdef := range this.tdefs {
//         tdefs[i] = tdef.substitute(ctx)
//     }

//     return newTupleType(this.labels, tdefs, this.line)
// }

/******************************************************************************
 * broker type
 ******************************************************************************/

type brokerType struct {
    baseNode
    gconfig typedef
    requester *participantType
}

func newBrokerType(gconfig typedef, requester *participantType, line int) *brokerType {
    return &brokerType {
         baseNode: baseNode {
            line: line,
        },
        gconfig: gconfig,
        requester: requester,
    }
}

func (bt *brokerType) setFilename(filename string) {
    bt.baseNode.setFilename(filename)
    bt.gconfig.setFilename(filename)
    if bt.requester != nil {
        bt.requester.setFilename(filename)
    }
}

func (bt *brokerType) subtypeOf(tdef typedef) bool {
    visited := util.NewHashSet[typePair]() 
    return bt.subtypeOf_(tdef, visited)
}

func (bt *brokerType) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
    switch btype := tdef.getType().(type) {
        case *nothingType:
            return true
        case *brokerType:
            if (bt.requester == nil) != (btype.requester == nil) {
                return false
            }

            if !bt.gconfig.subtypeOf_(btype.gconfig, visited) {
                return false
            }

            if bt.requester == nil /*&& btype.requester == nil*/ {
                return true
            }
            return bt.requester.subtypeOf(btype.requester)
        default:
            return false
    }
}

func (bt *brokerType) join(tdef typedef) (typedef, bool) {
    switch btype := tdef.getType().(type) {
        case *nothingType:
            return bt, true
        case *brokerType:
            gconfig, ok := bt.gconfig.join(btype.gconfig)
            broker := newBrokerType(gconfig, bt.requester, bt.line)
            if !ok {
                return broker, false
            }

            if (bt.requester == nil) != (btype.requester == nil) {
                return broker, false
            }

            if bt.requester == nil {
                return broker, true
            }
            return broker, bt.requester.subtypeOf(btype.requester) 
        default:
            return tdef, false
    }
}

// func (bt *brokerType) meet(tdef typedef) (typedef, bool) {
//     switch btype := tdef.getType().(type) {
//         case *nothingType:
//             return btype, true
//         case *brokerType:
//             gconfig, ok := bt.gconfig.meet(btype.gconfig)
//             ok = ok && bt.req.subtypeOf(btype.req)
//             return newBrokerType(gconfig, bt.req, bt.line), ok
//         default:
//             return bt, false
//     }
// }

func (bt *brokerType) getType() typedef { return bt }
func (bt *brokerType) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    if bt.requester != nil {
        bt.requester.typeCheck(ctx, log)
    }
    bt.gconfig.typeCheck(ctx, log)
    
    gconfig := bt.gconfig.getType();

    if gconfig == nil { return }
    gc, ok := gconfig.(globalConfig)
    if !ok {
        bt.gconfig.reportErrorf(log, "expecting global configuration type; instead found %q.", gconfig.String())
        return
    }

	if bt.requester == nil {
		return
	}

    found := false
    for _, participant := range gc.participants() {
        if participant == bt.requester.participant {
            found = true
            break
        }
    }

    if !found {
        stream := util.NewStream().Inc().Inc()
        bt.gconfig.prettyPrint(stream)
        bt.requester.reportErrorf(
            log,
            "broker type; global configuration:\n%s\n\tdoes not define request participant %q",
            stream.String(), bt.requester.String(),
        )
    }
}

func (bt *brokerType) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    if bt.requester != nil {
        bt.requester.projectionCheck(ctx, elog, rlog)
    }
    bt.gconfig.projectionCheck(ctx, elog, rlog)
}

func (bt *brokerType) defaultValue() expression {
    if bt.requester != nil {
        return newReqBroker(bt.requester.defaultValue().(*participantExpr), bt.gconfig, bt.line)
    }
    return newRecBroker(bt.gconfig, bt.line)
}

func (bt *brokerType) prettyPrint(iw util.IndentedWriter) {
    iw.Println("broker { ")
    iw.Inc()
    bt.gconfig.prettyPrint(iw)
    iw.Inc()
    iw.Print(" chan")
    iw.Dec()
    iw.Println()
    if bt.requester != nil {
        bt.requester.prettyPrint(iw)
        iw.Inc()
        iw.Print(" requester")
    }
    //iw.Println()
    iw.Dec()
    iw.Dec()
    iw.Println("}")
}

func (bt *brokerType) String() string {
    s := "broker { " + bt.gconfig.String() + " chan; " 
    if bt.requester != nil {
        s += bt.requester.String() + " requester;"
    }
    s += " }"
    return s
}

/******************************************************************************
 * arrow type
 ******************************************************************************/

// type arrowType struct {
//     baseNode
//     parameter typedef
//     body typedef
//     line int
// }
//
// func newArrowType(parameter typedef, body typedef, line int) (this *arrowType) {
//     this = new(arrowType)
//     this.parameter = parameter
//     this.body = body
//     this.init(line)
//     return
// }
//
// func (this *arrowType) setFilename(filename string) {
//     this.baseNode.setFilename(filename)
//     this.parameter.setFilename(filename)
//     this.body.setFilename(filename)
// }
//
// func (this *arrowType) subtypeOf(tdef typedef) bool {
//     arrow, ok := tdef.getType().(*arrowType)
//     if ok == true {
//         return  arrow.parameter.subtypeOf(this.parameter)  &&  this.body.subtypeOf(arrow.body)
//     }
//     return false
// }
//
// func (this *arrowType) submatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *arrowType:
//             if this.parameter.supermatch(td.parameter, ctx) == false {
//                 return false
//             }
//             if this.body.submatch(td.body, ctx) == false {
//                 return false
//             }
//             return true
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return this.subtypeOf(td2)
//     }
//     return false
// }
//
// func (this *arrowType) supermatch(tdef typedef, ctx substitutionContext) bool {
//     switch td := tdef.(type) {
//         case *arrowType:
//             if this.parameter.submatch(td.parameter, ctx) == false {
//                 return false
//             }
//             if this.body.supermatch(td.body, ctx) == false {
//                 return false
//             }
//             return true
//         case *typeVar:
//             td2, ok := ctx[td.name]
//             if ok == false {
//                 ctx[td.name] = this
//                 return true
//             }
//             return td2.subtypeOf(this)
//     }
//     return false
// }
//
// func (this *arrowType) substitute(ctx substitutionContext) typedef {
//     return newArrowType(this.parameter.substitute(ctx), this.body.substitute(ctx), this.line)
// }
//
// func (this *arrowType) getType() typedef {
//     return this
// }
//
// func (this *arrowType) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
//     this.parameter.typeCheck(ctx, log)
//     this.body.typeCheck(ctx, log)
//     return
// }
//
// func (this *arrowType) expressionCheck(ctx *typeCheckContext, log util.ErrorLog) typedef {
//     this.parameter.expressionCheck(ctx, log)
//     this.body.expressionCheck(ctx, log)
//     return this
// }
//
// func (this *arrowType) prettyPrint(stream *util.Stream) {
//     this.parameter.prettyPrint(stream)
//     stream.Print(" -> ")
//     this.body.prettyPrint(stream)
// }
//
// func (this *arrowType) String() string {
//     return this.parameter.String() + " -> " + this.body.String()
// }

/***
* create arrow type
***/

// func createArrowType(tdefs []typedef, line int) *arrowType {
//     if len(tdefs) == 1 {
//         return newArrowType(newEmptyTupleType(line), tdefs[0], line)
//     }
//     return newArrowType(tdefs[0], createArrowType_(tdefs, 1, line), line)
// }
//
// func createArrowType_(tdefs []typedef, index int, line int) typedef {
//     if index == len(tdefs) - 1 {
//         return tdefs[index]
//     } else {
//         return newArrowType( tdefs[index], createArrowType_(tdefs, index + 1, line), line )
//     }
// }

/******************************************************************************
 * name type
 ******************************************************************************/
type nameType struct {
    baseNode
    name string
    tdef typedef
}

func newNameType(name string, line int) *nameType {
    return &nameType {
        baseNode: baseNode {
            line: line,
        },
        name: name,
        tdef: nil,
    }
}

func (nt *nameType) getType() typedef { return nt.tdef }

func (nt *nameType) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    if nt.tdef = ctx.getType(nt.name); nt.tdef == nil {
        nt.reportErrorf(log, "undefined type: %q.", nt.name)
    }
}

func (_ *nameType) projectionCheck(_ *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}

func (nt *nameType) subtypeOf(tdef typedef) bool {
    visited := util.NewHashSet[typePair]()
    return nt.subtypeOf_(tdef, visited)
}

func (nt *nameType) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
    if nt.tdef == nil { return false }
    return nt.tdef.subtypeOf_(tdef, visited)
}

func (nt *nameType) join(tdef typedef) (typedef, bool) {
    if nt.tdef == nil { return tdef, false }
    return nt.tdef.join(tdef)
}

func (nt *nameType) defaultValue() expression {
    if nt.tdef == nil { return newNothingf(nt.line, "unknown type: %q", nt.name) }
    return nt.tdef.defaultValue()
}

// func (nt *nameType) meet(tdef typedef) (typedef, bool) {
//     if nt.tdef == nil { return nt, false }
//     return nt.tdef.meet(tdef)
// }

func (nt *nameType) prettyPrint(iw util.IndentedWriter) { iw.Print(nt.String()) }
func (nt *nameType) String() string                     { return nt.name }

/******************************************************************************
 * participantType type
 ******************************************************************************/

type participantType struct {
    baseNode
    participant string
}

func newParticipantType(participant string, line int) *participantType {
    return &participantType {
        baseNode: baseNode {
            line: line,
        },
        participant: participant,
    }
}

func (pt *participantType) getType() typedef { return pt }
func (_ *participantType) typeCheck(_ *typeCheckContext, log util.ErrorLog) {}
func (_ *participantType) projectionCheck(_ *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}
func (pt *participantType) subtypeOf(tdef typedef) bool { return pt.subtypeOf_(tdef, nil) }

func (pt *participantType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch ptype := tdef.getType().(type) {
        case *nothingType:
            return true
        case *participantType:
            return pt.participant == ptype.participant
        default:
            return false
    }
}

func (pt *participantType) join(tdef typedef) (typedef, bool) {
    switch ptype := tdef.getType().(type) {
        case *nothingType:
            return pt, true
        case *participantType:
            return pt, pt.participant == ptype.participant
        default:
            return tdef, false
    }
}

func (pt *participantType) defaultValue() expression {
    return newParticipantExpr(pt.participant, pt.line)
}

// func (pt *participantType) meet(tdef typedef) (typedef, bool) {
//     switch ptype := tdef.getType().(type) {
//         case *nothingType:
//             return ptype, true
//         case *participantType:
//             return pt, pt.participant == ptype.participant
//         default:
//             return pt, false
//     }
// }

func (pt *participantType) prettyPrint(iw util.IndentedWriter) { iw.Print(pt.String()) }
func (pt *participantType) String() string { return pt.participant }

/******************************************************************************
 * type variables and polymorphism
 *******************************************************************************/

// type nameVar struct {
//     baseNode
//     name string
//     tvar *typeVar
// }

// func newNameVar(name string, line int) *nameVar {
//     return &nameVar {
//         baseNode: baseNode {
//             line: line,
//         },
//         name: name,
//         tvar: nil,
//     }
// }

// func (nv *nameVar) getType() typedef { 
//     if nv.tvar == nil { return nil }
//     return nv.tvar.getType()
// }

// func (nv *nameVar) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
//     nv.tvar, _ = ctx.addTypeVar(nv.name)
// }

// func (nv *nameVar) subtypeOf(tdef typedef) bool {
//     visited := util.NewHashSet[typePair]()
//     return nv.subtypeOf_(tdef, visited)
// }

// func (nv *nameVar) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
//     if nv.tvar == nil { return false }
//     return nv.tvar.subtypeOf_(tdef, visited)
// }

// func (nv *nameVar) join(tdef typedef) (typedef, bool) {
//     if nv.tvar == nil { return tdef, false }
//     return nv.tvar.join(tdef)
// }

// func (nv *nameVar) prettyPrint(iw util.IndentedWriter) { iw.Print(nv.String()) }
// func (nv *nameVar) String() string                     { return nv.name }

// /****
//  * type variables
//  ****/

// type typeVar struct {
//     id          int           // unique id (debugging)
//     name        string        // e.g. "α"
//     bound       typedef       // nil if unbound
//     constraints []typeClass   // e.g. {NumTC, EqTC}
// }

// // monotonic counter
// var typeVarCounter int = 0
// func freshTypeVarID() int {
//     c := typeVarCounter
//     typeVarCounter++ 
//     return c
// }

// func newTypeVar(name string, cs ...typeClass) *typeVar {
//     return &typeVar {
//         id:          freshTypeVarID(),
//         name:        name,
//         constraints: cs,
//     }
// }

// func (v *typeVar) setFilename(filename string) {
//     if v.bound != nil {
//         v.bound.setFilename(filename)
//     }
// }

// func (v *typeVar) getType() typedef {
//     if v.bound != nil {
//         return v.bound.getType()
//     }
//     return v
// }

// func (v *typeVar) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
//     if v.bound != nil {
//         v.bound.typeCheck(ctx, log)
//     }
// }

// func (v *typeVar) subtypeOf(tdef typedef) bool { return v.subtypeOf_(tdef, nil) }
// func (v *typeVar) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
//     if v.bound != nil {
//         return v.bound.subtypeOf_(tdef, visited)
//     }
//     _, ok := v.join(tdef)
// //  t := tdef.getType()
//     // if t != nil {
//     // if tv, ok := t.(*typeVar); ok {
//     //     return v.name == tv.name
//     // }
//     // }
//     return ok //false 
// }

// func (v *typeVar) join(tdef typedef) (typedef, bool) {
//     td := tdef.getType()
//     if td == nil {
//         return v, false
//     }
//     if v.bound != nil {
//         return v.bound.join(td)
//     }
//     // TODO: consider running an occursIn method to avoid a -> [a] or a -> record {... id a ...}
//     // TODO: how about constrains? 
//     v.bound = td
//     return td, true
// }

// func (v *typeVar) prettyPrint(iw util.IndentedWriter) { iw.Print(v.String()) }
// func (v *typeVar) String() string {
//     if v.bound != nil { return v.bound.String() }
//     if len(v.constraints) == 0 { return v.name }
//     // e.g. "α:Num&Eq"
//     s := v.name + ":"
//     for i, c := range v.constraints {
//         if i > 0 { s += "&" }
//         s += c.String()
//     }
//     return s
// }

// /****
//  * type variable context
//  ****/

// type typeVarContext struct {
//     typeVariables map[string]*typeVar
// }

// func newTypeVarContext() *typeVarContext {
//     return &typeVarContext {
//         typeVariables: make(map[string]*typeVar),
//     }
// }

// func (tvc *typeVarContext) addTypeVar(name string) (*typeVar, bool) {
//     tv, ok := tvc.typeVariables[name];
//     if !ok {
//         tv = newTypeVar(name)
//         tvc.typeVariables[name] = tv
//     }
//     return tv, ok
// }

// // type typeVar struct {
// //     baseNode
// //     name string
// //     tdef typedef
// // }
// //
// // func newTypeVar(name string, tdef typedef, line int) (this *typeVar) {
// //     this = new(typeVar)
// //     this.name = name
// //     this.tdef = tdef
// //     this.init(line)
// //     return
// // }
// //
// // func (this *typeVar) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
// //     return
// // }
// //
// // func (this *typeVar) getType() typedef {
// //     return this.tdef
// // }
// //
// // //bool int float string list name typeVar
// // //any eq num ord list
// // func (this *typeVar) subtypeOf(tdef typedef) bool {
// //     return this.tdef.subtypeOf(tdef)
// // }
// //
// // func (this *typeVar) join(tdef typedef) (typedef, bool) {
// //     return this.tdef.join(tdef)
// // }
// //
// // // func (this *typeVar) submatch(tclass typeclass) bool {
// // //       switch tclass.kind {
// // //           case Any:
// // //               return true
// // //           case Eq:
// // //               return true
// // //           case Num:
// // //               return this.kind == Num || this.kind == Ord || this.kind == Integer || this.kind == Floating || this.kind == String || this.kind == List
// // //           case Ord:
// // //               return this.kind == Ord || this.kind == Integer || this.kind == Floating
// // //           case Boolean:
// // //               if this.kind == Boolean {
// // //                   this.tdef = newBoolType(this.line)
// // //                   return true
// // //               }
// // //           case Integer:
// // //               if this.kind == Integer {
// // //                   this.tdef = newIntType(this.line)
// // //                   return true
// // //               }
// // //           case Floating:
// // //               if this.kind == Floating {
// // //                   this.tdef = newFloatType(this.line)
// // //                   return true
// // //               }
// // //           case String:
// // //               if this.kind == String {
// // //                   this.tdef = newStringType(this.line)
// // //                   return true
// // //               }
// // //           case List:
// // //               return //this.tdef.submatch(tclass.tclass)
// // //     }
// // //     return false
// // // }
// //
// // // func (this *typeVar) submatch(tdef typedef, ctx substitutionContext) bool {
// // //     switch td := tdef.(type) {
// // //         case *typeVar:
// // //             td2, ok := ctx[td.name]
// // //             if ok == false {
// // //                 ctx[td.name] = this
// // //                 return true
// // //             }
// // //             return this.subtypeOf(td2)
// // //     }
// // //     return this.subtypeOf(tdef)
// // // }
// // //
// // // func (this *typeVar) supermatch(tdef typedef, ctx substitutionContext) bool {
// // //     switch td := tdef.(type) {
// // //         case *typeVar:
// // //             td2, ok := ctx[td.name]
// // //             if ok == false {
// // //                 ctx[td.name] = this
// // //                 return true
// // //             }
// // //             return td2.subtypeOf(this)
// // //       }
// // //
// // //     //td, ok := ctx[this.name]
// // //     //if ok == false {
// // //     //    ctx[this.name] = tdef
// // //     //    return true
// // //     //}
// // //     return tdef.subtypeOf(this)
// // // }
// // //
// // // func (this *typeVar) substitute(ctx substitutionContext) typedef {
// // //     if ctx[this.name] == nil {
// // //         return this
// // //     }
// // //     return ctx[this.name]
// // //   //  return this
// // //   //  tdef, ok := ctx[this.name]
// // //   //  if ok == false {
// // //   //      return this
// // //   //  }
// // //     //if tdef.String() == this.String() {
// // //     //    return tdef
// // //     //}
// // //     //fmt.Println(tdef.String(), this.String())
// // //     //return tdef.substitute(ctx)
// // //     //var tmp typedef = nil //tdef.substitute(ctx)
// // //     //for {
// // //     //    if tmp == tdef {
// // //     //        break
// // //     //    }
// // //     //    tmp = tdef.substitute(ctx)
// // //     //    tdef = tmp
// // //     //    //fmt.Println(tmp.String(), tdef.String())
// // //     //}
// // // //    return tdef
// // // }
// //
// // func (this *typeVar) prettyPrint(stream *util.Stream) {
// //     stream.Print(this.String())
// // }
// //
// // func (this *typeVar) String() string {
// //     return this.name + ":" + this.tdef.String()
// // }

/******************************************************************************
 * types and typeclasses
 *****************************************************************************/

// TODO Possibly even add polymorphism
type typeClass interface {
    contains(typedef) bool
    fmt.Stringer
}

// Zero-sized types (no fields).
type anyTC struct{}
type eqTC  struct{}
type numTC struct{}
type ordTC struct{}

func (anyTC) contains(_ typedef) bool { return true }
func (anyTC) String() string          { return "Any" }

func (eqTC) contains(tdef typedef) bool {
    switch tdef.getType().(type) {
        case *boolType, *intType, *floatType, *stringType, *listType:
            return true
        default:
            return false
    }
}
func (eqTC) String() string { return "Eq" }

func (numTC) contains(tdef typedef) bool {
    switch tdef.getType().(type) {
        case *intType, *floatType, *stringType, *listType:
            return true
        default:
            return false
    }
}
func (numTC) String() string { return "Num" }

func (ordTC) contains(tdef typedef) bool {
    switch tdef.getType().(type) {
        case *intType, *floatType, *stringType:
            return true
        default:
            return false
    }
}
func (ordTC) String() string { return "Ord" }

// Singletons (no allocations; zero values).
var (
    anytc typeClass = anyTC{}
    eq  typeClass = eqTC{}
    num typeClass = numTC{}
    ord typeClass = ordTC{}
)

/******************************************************************************
* localAbstraction
*******************************************************************************/

type localAbstraction struct {
    baseNode
    parameters []*participantType
    loc local
}

func newLocalAbstraction(parameters []*participantType, loc local, line int) *localAbstraction {
    return &localAbstraction {
        baseNode: baseNode {
            line: line,
        },
        parameters: parameters,
        loc: loc,
    }
}

func (la *localAbstraction) setFilename(filename string) {
    la.baseNode.setFilename(filename)
    for _, p := range la.parameters {
        p.setFilename(filename)
    }
    la.loc.setFilename(filename)
}

func (la *localAbstraction) getType() typedef { return la }
func (la *localAbstraction) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    la.loc.typeCheck(ctx, log)
}

func (la *localAbstraction) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    ctx.push()
    defer ctx.pop()
    for _, p := range la.parameters {
        if ctx.addParticipant(p) == false {
            p.reportErrorf(elog, "duplicated participant definition: %q.", p.String())
        }
    }
    for _, p := range la.parameters {
        p.projectionCheck(ctx, elog, rlog)
    }
    la.loc.projectionCheck(ctx, elog, rlog)
}

func (la *localAbstraction) subtypeOf(tdef typedef) bool {
    visited := util.NewHashSet[typePair]()
    return la.subtypeOf_(tdef, visited)
}

func (la *localAbstraction) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
    tPair := typePair{a: la, b:tdef}

    if visited.Contains(tPair) == true {
        return true
    }

    visited.Add(tPair)

    switch locAbstr := tdef.getType().(type) {
        case *localAbstraction:
            if (len(la.parameters) != len(locAbstr.parameters)) {
                // error no subtype
                return false
            }
            labstr := locAbstr.substitute(la.parameters)

            return la.loc.subtypeOf(labstr.loc, visited)
        default:
            return false
    }
}

func (la *localAbstraction) join(tdef typedef) (typedef, bool) {
    r, ok := tdef.getType().(*localAbstraction)
    if !ok || r == nil {
        return la, false
    }
    if len(la.parameters) != len(r.parameters) {
        return la, false
    }
    for i, p := range la.parameters {
        if !p.subtypeOf(r.parameters[i]) {
            return la, false
        }
    }

    loc, result := la.loc.join(r.loc)
    return newLocalAbstraction(la.parameters, loc, la.line), result
}

// func (la *localAbstraction) meet(tdef typedef) (typedef, bool) {
//     r, ok := tdef.getType().(*localAbstraction)
//     if !ok || r == nil {
//         return la, false
//     }
//     if len(la.parameters) != len(r.parameters) {
//         return la, false
//     }
//     for i, p := range la.parameters {
//         if !p.subtypeOf(r.parameters[i]) {
//             return la, false
//         }
//     }

//     loc, result := la.loc.meet(r.loc)
//     return newLocalAbstraction(la.parameters, loc, la.line), result
// }

func (la *localAbstraction) substitute(ptypes []*participantType) *localAbstraction {
    substitution := make(map[string]string)
    parameters := make([]*participantType, len(la.parameters))
    for i, p := range la.parameters {
        if i < len(ptypes) {
            substitution[p.participant] = ptypes[i].participant
            parameters[i] = ptypes[i]
        } else {
            substitution[p.participant] = p.participant
            parameters[i] = p
        }
    }
    loc := la.loc.substitute(substitution)
    return newLocalAbstraction(parameters, loc, la.line)
}

func (la *localAbstraction) apply(arguments []*participantExpr) local {
    if (len(la.parameters) != len(arguments)) {
        // error no join
        return la.loc
    }
    substitution := make(map[string]string)
    for i := range la.parameters {
        substitution[la.parameters[i].participant] = arguments[i].id
    }
    return la.loc.substitute(substitution)
}

func (la *localAbstraction) defaultValue() expression {
    return newNothingf(la.line, "unexpected local abstraction type")
    // parameters := make([]*participantExpr, len(la.parameters))
    // for i, par := range la.parameters {
    //     parameters[i] = par.defaultValue()
    // }
    // body := la.loc.defaultValue()
    // return newAbstraction(body, parameters, la.line)
}

func (la *localAbstraction) prettyPrint(iw util.IndentedWriter) {
    iw.Print("local")
    for _, participant := range la.parameters {
        iw.Print(" ")
        participant.prettyPrint(iw)
    }
    iw.Print(". ")
    iw.Inc()
    iw.Println("")
    la.loc.prettyPrint(iw)
    iw.Dec()
    return
}

func (la *localAbstraction) String() string {
    s := "local"
    for _, participant := range la.parameters {
        s += " " + participant.String()
    }
    return s + ". " + la.loc.String()
}

/******************************************************************************
* projection
*******************************************************************************/

type projection struct {
    baseNode
    ptype       *participantType
    conf         typedef
    locAbstr    *localAbstraction
}

func newProjection(conf typedef, ptype *participantType, line int) *projection {
    return &projection {
        baseNode: baseNode {
            line: line,
        },
        ptype: ptype,
        conf: conf,
    }
}

func (p *projection) getType() typedef { return p.locAbstr }

func (p *projection) typeCheck(ctx *typeCheckContext, log util.ErrorLog) {
    p.conf.typeCheck(ctx, log)
    conf := p.conf.getType()
    if conf == nil {
        return
    }
    _, ok := conf.(globalConfig)
    if !ok {
        p.reportErrorf(log, "expecting global type; instead found %q", conf.String())
        return
    }
}

func (p *projection) projectionCheck(ctx *projectionCheckContext, elog util.ErrorLog, rlog util.ReportLog) {
    p.conf.projectionCheck(ctx, elog, rlog)
    // TODO: ensure p.ptype is one of gconfig.participants()
    gconfig := p.conf.getType().(globalConfig)
    p.locAbstr = gconfig.project(p.ptype)
}

func (p *projection) subtypeOf(tdef typedef) bool {
    visited := util.NewHashSet[typePair]()
    return p.subtypeOf_(tdef, visited)
}

func (p *projection) subtypeOf_(tdef typedef, visited *util.HashSet[typePair]) bool {
    if p.locAbstr == nil { return false }
    return p.locAbstr.subtypeOf_(tdef, visited)
}

func (p *projection) join(tdef typedef) (typedef, bool) {
    if p.locAbstr == nil { return p, false }
    return p.locAbstr.join(tdef)
}

// func (p *projection) meet(tdef typedef) (typedef, bool) {
//     if p.locAbstr == nil { return p, false }
//     return p.locAbstr.meet(tdef)
// }

func (p *projection) apply(arguments []*participantExpr) local {
    return p.locAbstr.apply(arguments)
}

func (p *projection) defaultValue() expression {
    if p.locAbstr == nil {
        return newNothingf(p.line, "unprojected type: %q", p.String())
    }
    return p.locAbstr.defaultValue()
}

func (p *projection) prettyPrint(iw util.IndentedWriter) {
    p.conf.prettyPrint(iw)
    iw.Print("(")
    p.ptype.prettyPrint(iw)
    iw.Print(")")
}

func (p *projection) String() string {
    s := p.conf.String() + "("
    s += p.ptype.String()
    return s + ")"
}

/**********************************************************************************
 * ioType
 **********************************************************************************/

 type ioType struct {
    baseNode
 }

 func newioType(line int) *ioType {
    return &ioType {
        baseNode{line: line},
    }
 }

func (iot *ioType) getType() typedef { return iot }
func (iot *ioType) subtypeOf(tdef typedef) bool { return iot.subtypeOf_(tdef, nil) }
func (iot *ioType) subtypeOf_(tdef typedef, _ *util.HashSet[typePair]) bool {
    switch tdef.(type) {
        case *nothingType, *ioType:
            return true
        default:
            return false
    }
}

func (iot *ioType) join(tdef typedef) (typedef, bool) {
    switch tdef.(type) {
        case *nothingType, *ioType:
            return iot, true
        default:
            return tdef, false
    }
}

// func (iot *ioType) meet(tdef typedef) (typedef, bool) {
//     switch tdef.(type) {
//         case *nothingType, *ioType:
//             return tdef, true
//         default:
//             return iot, false
//     }
// }

func (iot *ioType) typeCheck(_ *typeCheckContext, _ util.ErrorLog) {}
func (iot *ioType) projectionCheck(_ *projectionCheckContext, _ util.ErrorLog, _ util.ReportLog) {}


func (iot *ioType) defaultValue() expression {
    return newEmptyHandle(iot.line)
}

func (iot *ioType) prettyPrint(iw util.IndentedWriter) { iw.Print( iot.String() ) }
func (_ *ioType) String() string { return "io" }
