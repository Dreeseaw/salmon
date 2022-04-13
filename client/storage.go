/*
   Uses kelindar/column for single-node, in-mem
   storage & operations
*/
package salmon

import (
	"errors"
//    "fmt"

	"github.com/kelindar/column"
)

type ColFunc func() column.Column

var CollectionTypeMap = map[string]ColFunc{
	"string": column.ForString,
	"float":  column.ForFloat64,
	"int":    column.ForInt32,
	"bool":   column.ForBool,
}

type AggFunc func(a, b interface{}) interface{}

var AggFuncs = map[string]AggFunc{
    "sum": func(a, b interface{}) interface{} {
        switch v := a.(type) {
        case int32:
            return v + b.(int32)
        case float64:
            return v + b.(float64)
        default:
            panic(a) //TODO
        }
    },
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

func (ta *Table) Select(selectors []Selector, filters []filter) ([]Object, error) {
    
    result_rows := make([]Object, 0)

    isAggregated := false
    for _, sel := range selectors {
        if sel.Aggregator != "" {
            isAggregated = true
        }
    }

    if isAggregated {
        for _, sel := range selectors {
            if sel.Aggregator == "" {
                return nil, errors.New("Mix of aggregated & non-aggregated columns")
            }
        }
    }

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
                    return f.Process(v)
                })
            }
        }

        if isAggregated {
            aggrObj := make(Object)
            return txn.Range(func (i uint32) {
                for _, sel := range selectors {
                    value, _ := txn.Any(sel.ColName).Get()
                    aggFunc, _ := AggFuncs[sel.Aggregator]
                    ta.AggrCollector(aggrObj, sel, value, aggFunc)
                }
                result_rows = append(result_rows, aggrObj)
            })
            return nil
        }

        // range and return selected cols
        return txn.Range(func (i uint32) {
            row_obj := make(Object)
            for _, sel := range selectors {
                value, _ := txn.Any(sel.ColName).Get()
                row_obj[sel.ColName] = value
            }
            result_rows = append(result_rows, row_obj)
        })
        return nil
    })

    if isAggregated {
        return result_rows[len(result_rows)-1:], nil
    }
    return result_rows, nil
}

func (ta *Table) AggrCollector(obj Object, sel Selector, val interface{}, af AggFunc) {
    if prev, exists := obj[sel.ColName]; exists {
        obj[sel.ColName] = af(prev, val)
    } else {
        obj[sel.ColName] = val
    }
}
