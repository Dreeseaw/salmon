/*
   Uses kelindar/column for single-node, in-mem
   storage & operations
*/
package main

// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"

	"github.com/kelindar/column"
)

type ColFunc func() column.Column
type CollTypeList []string

var CollectionTypeMap = map[string]ColFunc{
	"string": column.ForString,
	"float":  column.ForFloat64,
	"int":    column.ForInt32,
	"bool":   column.ForBool,
}

type Store struct {
	CollMap         map[string]*column.Collection
	CollMetadataMap map[string]CollTypeList
}

func NewStore() *Store {
	cm := make(map[string]*column.Collection)
	mm := make(map[string]CollTypeList)
	ret := Store{
		CollMap:         cm,
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
	ctp := make(CollTypeList, 0)
	s.CollMetadataMap[name] = ctp
	return nil
}

func (s *Store) AddColumn(coll, cn, ct string) error {
	if collection, exists := s.CollMap[coll]; exists {
		if colfunc, exists := CollectionTypeMap[ct]; exists {
			collection.CreateColumn(cn, colfunc())
			ctp, _ := s.CollMetadataMap[coll]
			s.CollMetadataMap[coll] = append(ctp, ct)
			return nil
		}
		fmt.Printf("Column type does not exist")
		return errors.New("Column type does not exist")
	}
	fmt.Printf("Collection does not exist")
	return errors.New("Collection does not exist")
}

// yep need generics
func (s *Store) AddObject(coll string, obj []interface{}) error {
	if collection, exists := s.CollMap[coll]; exists {
		ctp, _ := s.CollMetadataMap[coll]
		obj_map := make(map[string]interface{})
		// obj_map_index := 0

		for col_i, colname := range ctp {
			// fmt.Printf("colname: %v", colname)
			obj_map[colname] = obj[col_i]
			// obj_map_index += 1
		}

		collection.InsertObject(obj_map)
		return nil
	}
	return errors.New("Collection does not exist")
}
