/*
   Uses kelindar/column for single-node, in-mem
   storage & operations
*/
package salmon

import (
//	"errors"
    "fmt"

	"github.com/kelindar/column"
)

type ColFunc func() column.Column

var CollectionTypeMap = map[string]ColFunc{
	"string": column.ForString,
	"float":  column.ForFloat64,
	"int":    column.ForInt32,
	"bool":   column.ForBool,
}

type Table struct {
    Coll *column.Collection
    Meta TableMetadata
}

func orderColList(tm TableMetadata) []ColumnMetadata {
    ret := make([]ColumnMetadata, len(tm))
    for _, colMeta := range tm {
        ret[colMeta.Order] = colMeta
    }
    return ret
} 

func NewTable(tm TableMetadata) *Table {
    coll := column.NewCollection()

    // create columns in correct order
    for _, colMeta := range orderColList(tm) {
        colFunc, _ := CollectionTypeMap[colMeta.Type]
        coll.CreateColumn(colMeta.Name, colFunc())
    }

    return &Table{
        Coll: coll,
        Meta: tm,
    }
}

func (ta *Table) InsertObject(obj map[string]interface{}) error {
    ta.Coll.InsertObject(obj)
    return nil
}

func (ta *Table) Select(selectors []string, filters []filter) ([]Object, error) {
    
    result_rows := make([]Object, 0)

    ta.Coll.Query(func(txn *column.Txn) error {

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
                    fmt.Println(v)
                    return f.Process(v)
                })
            }
        }

        // range and return selected data
        return txn.Range(func (i uint32) {
            row_obj := make(Object)
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
