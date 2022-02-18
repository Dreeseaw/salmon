/*
    Uses kelindar/column for single-node, in-mem 
    storage & operations
*/
package main

import (
    "errors"

    "github.com/kelindar/column"
)

type ColFunc func() column.Column

var CollectionTypeMap = map[string]ColFunc{
    "string": column.ForEnum,
    "float": column.ForFloat64,
    "int": column.ForInt16,
    "bool": column.ForBool,
}

type Store struct {
    CollMap map[string]*column.Collection
}

func NewStore() *Store {
    cm := make(map[string]*column.Collection)
    ret := Store{
        CollMap: cm,
    }
    return &ret
}

func (s *Store) NewCollection(name string) error {
    if _, exists := s.CollMap[name]; exists {
        return errors.New("Collection with that name already exists")
    }
    newColl := column.NewCollection()
    newColl.CreateColumn("serial", column.ForKey())
    s.CollMap[name] = newColl
    return nil
}

func (s *Store) AddColumn(coll, cn, ct string) error {
    if collection, exists := s.CollMap[coll]; exists {
        if colfunc, exists := CollectionTypeMap[ct]; !exists {
            collection.CreateColumn(cn, colfunc())
            return nil
        }
        return errors.New("Column type does not exist")
    } 
    return errors.New("Collection does not exist")
}
