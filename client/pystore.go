/*
    CPython interface to cache
*/
package main

import (
    "C"
    "fmt"
)

var StorePtr *Store = nil

//export Initstore
func Initstore() bool {
    if StorePtr != nil {
        return false
    }
    StorePtr = NewStore()
    return true
}

//export Newcoll
func Newcoll(name *C.char) bool {
    if StorePtr == nil {
        return false
    }
    nameStr := C.GoString(name)
    fmt.Printf(nameStr)
    if err := StorePtr.NewCollection(nameStr); err != nil {
        return false
    }
    return true
}

//export AddCol
func AddCol(collection *C.char, colname *C.char, coltype *C.char) bool {
    if StorePtr == nil {
        return false
    }
    cl := C.GoString(collection)
    cn := C.GoString(colname)
    ct := C.GoString(coltype)
    fmt.Printf(cl+cn+ct)
    if err := StorePtr.AddColumn(cl, cn, ct); err != nil {
        return false
    }
    return true
}

//export Set
func Set(coll string) bool {
    return false
}

//export Get
func Get(coll string) bool {
    return false
}

func main() {}
