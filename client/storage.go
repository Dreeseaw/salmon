/*
   Uses kelindar/column for single-node, in-mem
   storage & operations
*/
package main

import (
	"errors"
	"fmt"

	"github.com/kelindar/column"
)

type ColFunc func() column.Column
type CollTypeList []ColMetadata

var CollectionTypeMap = map[string]ColFunc{
	"string": column.ForString,
	"float":  column.ForFloat64,
	"int":    column.ForInt32,
	"bool":   column.ForBool,
}

type ColMetadata struct {
    Name   string
    Type   string
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
        ctp, _ := s.CollMetadataMap[coll]
        for _, colMeta := range ctp {
            if cn == colMeta.Name {
                return errors.New("Column name already exists")
            }
        }
		if colfunc, exists := CollectionTypeMap[ct]; exists {
			collection.CreateColumn(cn, colfunc())
			s.CollMetadataMap[coll] = append(ctp, ColMetadata{
                Name: cn,
                Type: ct,
            })
			return nil
		}
		fmt.Printf("Column type does not exist")
		return errors.New("Column type does not exist")
	}
	fmt.Printf("Collection does not exist")
	return errors.New("Collection does not exist")
}

// AddObject adds an object to a collection, does not mock SQl insert
func (s *Store) AddObject(coll string, obj []interface{}) error {
	if collection, exists := s.CollMap[coll]; exists {
		ctp, _ := s.CollMetadataMap[coll]
		obj_map := make(map[string]interface{})

		for col_i, col_meta := range ctp {
			obj_map[col_meta.Name] = obj[col_i]
		}

		collection.InsertObject(obj_map)
		return nil
	}
	return errors.New("Collection does not exist")
}

// TODO buffer result rows back to user (cursor)
// Select mocks a SQL-flavored select key word
func (s *Store) Select(coll string, selectors []string, filters []filter) ([]column.Object, error) {
    
    collection, exists := s.CollMap[coll]
    if !exists {
        return nil, errors.New("Collection does not exist")
    }

    result_rows := make([]column.Object, 0)

    collection.Query(func(txn *column.Txn) error {

        // filter rows, account for bool
        for _, f := range filters {
            bf, ok := f.(BoolFilter)
            if ok {
                if bf.Val {
                    txn = txn.With(f.ColName())
                } else {
                    txn = txn.Without(f.ColName())
                }
            } else {
                txn = txn.WithValue(f.ColName(), func(v interface{}) bool {
                    return f.Process(v)
                })
            }
        }

        // range and return selected data
        return txn.Range(func (i uint32) {
            row_obj := make(column.Object)
            for _, sel := range selectors {
                value, _ := txn.Any(sel).Get()
                row_obj[sel] = value
            }
            result_rows = append(result_rows, row_obj)
        })
        return nil
    })

    return result_rows, nil
}
