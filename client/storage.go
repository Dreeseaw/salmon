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
//     "reflect"

	"github.com/kelindar/column"
)

type ColFunc func() column.Column
type ColReaderFunc func(*column.Txn, string) column.Reader
type CollTypeList []ColMetadata

var CollectionTypeMap = map[string]ColFunc{
	"string": column.ForString,
	"float":  column.ForFloat64,
	"int":    column.ForInt32,
	"bool":   column.ForBool,
}

/*
var ColumnReaderMap = map[string]ColReaderFunc{
    "string": (*column.Txn).String,
    "float":  (*column.Txn).Float64,
    "int":    (*column.Txn).Int32,
    "bool":   (*column.Txn).Bool,
}
*/

// type ColReader = column.anyWriter

type ColMetadata struct {
    Name   string
    Type   string
    // Reader column.Reader
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
            //reader_func, _ := ColumnReaderMap[cn]
			s.CollMetadataMap[coll] = append(ctp, ColMetadata{
                Name: cn,
                Type: ct,
            //    Reader: reader_func,
            })
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

        // list -> map w/ column names
		for col_i, col_meta := range ctp {
			obj_map[col_meta.Name] = obj[col_i]
		}

		collection.InsertObject(obj_map)
		return nil
	}
	return errors.New("Collection does not exist")
}

func (s *Store) Select(coll string, selectors []string, filters map[string]interface{}) error {
    
    collection, exists := s.CollMap[coll]
    if !exists {
        return errors.New("Collection does not exist")
    }

    // collMeta, _ := s.CollMetadataMap[coll]
    // selector_funcs := make([]column.Reader, len(selectors))

    /*
    for sel_i, sel := range selectors {
        for _, col_meta := range collMeta {
            if sel == col_meta.Name {
                rd_func, _ := ColumnReaderMap[col_meta.Type]
                selector_funcs[sel_i] = rd_func
            }
        }
    }
    */

    // cursor := make([]interface{})

    collection.Query(func(txn *column.Txn) error {
        // for testing
        // strs := txn.String("testcolstr")

        // create needed column readers
        /*
        colReaders := make(map[string]ColReader)
        for _, sel := range selectors {
            tmp_rd := txn.Any(sel)
            colReaders[sel] = tmp_rd
            fmt.Println(reflect.TypeOf(tmp_rd))
            tp, _ := tmp_rd.Get()
            fmt.Println("tp:")
            fmt.Println(tp)
        }
        */

        // filter rows
        for colname, colval := range filters {
            txn = txn.WithValue(colname, func(v interface{}) bool {
                return v == colval
            })
        }

        // TODO add objects to cursor list
        return txn.Range(func (i uint32) {
            for _, sel := range selectors {
                //reader, _ := colReaders[sel]
                reader := txn.Any(sel)
                to_print, _ := reader.Get()
                fmt.Println(to_print)
            }
        })
        return nil
    })

    return nil
}
