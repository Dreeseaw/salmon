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
func Select(
    coll *C.char,
    sel_payload *C.char,
    s_size C.int,
    filt_payload *C.char,
    f_size C.int,
) C.int {

    coll_str := C.GoString(coll)

    filt_list, err := r_object(
        C.GoBytes(unsafe.Pointer(filt_payload), f_size),
        nil,
    )
    if err != nil {
        return 0
    }

    filters := make(map[string]interface{})
    for _, filter := range filt_list {
        filt_tuple := filter.([]interface{})
        filters[filt_tuple[0].(string)] = filt_tuple[1]
    }

    sel_list, err := r_object(
        C.GoBytes(unsafe.Pointer(sel_payload), s_size),
        nil,
    )
    if err != nil {
        return 0 
    }

    selectors := make([]string, len(sel_list))
    for s_i, selector := range sel_list {
        selectors[s_i] = selector.(string)
    }

    results, err := StorePtr.Select(coll_str, selectors, filters)
	if err != nil {
        return 0
    }

    // print results out for now
    for res_row := range results {
        fmt.Println(res_row)
    }
    return 1

    // marshaled_rows, err := w_object(results)
    // ret_size = (*C.int)((unsafe.Pointer)(len(marshaled_rows)))
    // c_rows := C.CBytes(marshaled_rows)
    // defer C.free(c_rows)
    // return (*C.uchar)(c_rows)

}

func main() {}
