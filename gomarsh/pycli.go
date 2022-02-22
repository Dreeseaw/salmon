/*
    CPython interface to cache
*/
package main

// #include <stdlib.h>
import "C"

import (
    "fmt"
    "unsafe"
)

//export TestFunc
func TestFunc(payload *C.char, p_size C.int) bool {
    // cn := C.GoString(coll)
    // p_len := len(C.GoString(payload))

    p_gb := C.GoBytes(unsafe.Pointer(payload), p_size)
    tmp_obj, err := unmarshal_tuple(p_gb)
    if err != nil {
        panic(err)
    }
    fmt.Println(len(tmp_obj))
    for _, obj := range tmp_obj {
        fmt.Println(obj)
    }

    return true
}

func main() {}
