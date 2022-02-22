/*
    CPython interface to cache
*/
package main

// #cgo pkg-config: python3
// #include <stdlib.h>
// #include <Python.h>
import "C"

import (
    "fmt"
    "unsafe"
    
//    py3 "github.com/go-python/cpy3" 
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
    if err := StorePtr.NewCollection(nameStr); err != nil {
        return false
    }
    return true
}

//export AddCol
func AddCol(collection, colname, coltype *C.char) bool {
    if StorePtr == nil {
        return false
    }
    cl := C.GoString(collection)
    cn := C.GoString(colname)
    ct := C.GoString(coltype)
    if err := StorePtr.AddColumn(cl, cn, ct); err != nil {
        return false
    }
    return true
}

//export Insert
func Insert(coll *C.char, payload *C.char, p_size C.int) bool {
    // cn := C.GoString(coll)
    // p_len := len(C.GoString(payload))

    p_gb := C.GoBytes(unsafe.Pointer(payload), p_size)
    for one_byte := range p_gb {
        fmt.Printf("b")
        fmt.Printf(string(one_byte))
    }

    // b := make([]byte, 24)
    // b = (*[1<<30]byte)(unsafe.Pointer(&payload))[0:24]
    // for one_byte := range b {
    //     fmt.Printf(string(one_byte))
    // }
    // obj := (*C.PyObject)(pyobj)
    // if err := StorePtr.AddObject(cn, obj); err != nil {
    //     return false
    // }
    return true
}

//export Get
func Get(coll string) bool {
    return false
}

func main() {}
