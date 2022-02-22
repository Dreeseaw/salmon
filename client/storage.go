/*
    Uses kelindar/column for single-node, in-mem 
    storage & operations
*/
package main

// #cgo pkg-config: python3
// #include <stdlib.h>
// #include <Python.h>
import "C"
import (
    "fmt"
    "errors"
    "unsafe"

    "github.com/kelindar/column"
)

type ColFunc func() column.Column
type CollTypeMap map[string]string
type GoTypeFunc func(pyobj *C.PyObject) interface{} 

var CollectionTypeMap = map[string]ColFunc{
    "string": column.ForString,
    "float": column.ForFloat64,
    "int": column.ForInt64,
    "bool": column.ForBool,
}

func PyToGoString(pyobj *C.PyObject) interface{} {
    return C.PyByteArray_AsString(pyobj)
}

func PyToGoFloat(pyobj *C.PyObject) interface{} {
    return C.PyFloat_AsDouble(pyobj)
}

func PyToGoInt(pyobj *C.PyObject) interface{} {
    return C.PyLong_AsLong(pyobj)
}

func PyToGoBool(pyobj *C.PyObject) interface{} {
    return !(C.PyLong_AsLong(pyobj) == 0)
}

var GoTypeFuncMap = map[string]GoTypeFunc{
    "string": PyToGoString,
    "float": PyToGoFloat,
    "int": PyToGoInt,
    "bool": PyToGoBool,
}

type Store struct {
    CollMap         map[string]*column.Collection
    CollMetadataMap map[string]CollTypeMap
}

func NewStore() *Store {
    cm := make(map[string]*column.Collection)
    mm := make(map[string]CollTypeMap)
    ret := Store{
        CollMap: cm,
        CollMetadataMap: mm,
    }
    return &ret
}

func (s *Store) NewCollection(name string) error {
    if _, exists := s.CollMap[name]; exists {
        return errors.New("Collection with that name already exists")
    }
    newColl := column.NewCollection()
    // newColl.CreateColumn("serial", column.ForKey())
    s.CollMap[name] = newColl
    ctp := make(CollTypeMap)
    s.CollMetadataMap[name] = ctp
    return nil
}

func (s *Store) AddColumn(coll, cn, ct string) error {
    if collection, exists := s.CollMap[coll]; exists {
        if colfunc, exists := CollectionTypeMap[ct]; exists {
            collection.CreateColumn(cn, colfunc())
            ctp, _ := s.CollMetadataMap[coll]
            ctp[cn] = ct
            return nil
        }
        fmt.Printf("Column type does not exist")
        return errors.New("Column type does not exist")
    } 
    fmt.Printf("Collection does not exist")
    return errors.New("Collection does not exist")
}

func (s *Store) AddObject(coll string, obj []byte) error {
    if collection, exists := s.CollMap[coll]; exists {
        ctp, _ := s.CollMetadataMap[coll]
        tmpobj := make(map[string]interface{})
        odir := C.PyObject_Dir(obj)
        ostr := C.PyByteArray_AsString(odir)
        gostr := C.GoString(ostr)
        fmt.Printf(gostr)

        for colname, coltype := range ctp {
            colname_c := C.CString(colname)
            defer C.free(unsafe.Pointer(colname_c))
            type_obj := C.PyObject_GetAttrString(obj, colname_c)
            tmpobj[colname] = GoTypeFuncMap[coltype](type_obj)
        }
        collection.InsertObject(tmpobj)
        return nil
    }
    return errors.New("Collection does not exist")
}

