/*
   CPython interface to cache
*/
package main

// #include <stdlib.h>
import "C"

import (
	// "fmt"
	"unsafe"
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
	cn := C.GoString(coll)
	go_bytes := C.GoBytes(unsafe.Pointer(payload), p_size)
	obj, err := r_object(go_bytes, nil)
	if err != nil {
		return false
	}

	if err = StorePtr.AddObject(cn, obj); err != nil {
		return false
	}

	return true
}

//export Select
func Select(coll *C.char, sel_payload *C.char, s_size C.int) bool {

    coll_str := C.GoString(coll)

    sel_list, err := r_object(
        C.GoBytes(unsafe.Pointer(sel_payload), s_size),
        nil,
    )
    if err != nil {
        return false
    }

    selectors := make(map[string]interface{})
    for _, selector := range sel_list {
        sel_tuple := selector.([]interface{})
        selectors[sel_tuple[0].(string)] = sel_tuple[1]
    }
   
    err = StorePtr.Select(coll_str, selectors)
	if err != nil {
        return false
    }
    return true
}

func main() {}
