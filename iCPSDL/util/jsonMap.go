package util

//import "sync"
//import "context"
//import "time"

/******************************************************************************
 * json mini parser
 ******************************************************************************/

 type JsonMap map[string]interface{}
 type jsonArray []interface{}

func (this JsonMap) GetString(name string) string {
    return this[name].(string)
}

func (this JsonMap) GetInt(name string) int {
    return this[name].(int)
}

func (this JsonMap) GetFloat(name string) float64 {
    return this[name].(float64)
}

func (this JsonMap) GetMap(name string) JsonMap {
    return this[name].(JsonMap)
}

func (this JsonMap) GetArray(name string) jsonArray {
    return this[name].(jsonArray)
}

// func (this jsonArray) getString(index int) string {
//     return this[index].(string)
// }
//
// func (this jsonArray) getInt(index int) int {
//     return this[index].(int)
// }
//
// func (this jsonArray) getFloat(index int) float64 {
//     return this[index].(float64)
// }
//
// func (this jsonArray) getMap(index int) JsonMap {
//     return this[index].(JsonMap)
// }
//
// func (this jsonArray) getArray(index int) jsonArray {
//     return this[index].(jsonArray)
//}
