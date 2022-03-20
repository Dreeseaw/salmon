/*
    Storage tests
*/
package main

import (
    "fmt"
    "testing"

    "github.com/kelindar/column"
    "github.com/stretchr/testify/assert"
)

func Test_InsertObject(t *testing.T) {
    tm := TableMetadata{
        "testcolumnint": ColumnMetadata{Type: "int"},
        "testcolumnstr": ColumnMetadata{Type: "string"},
        "testcolumnfloat": ColumnMetadata{Type: "float"},
        "testcolumnbool": ColumnMetadata{Type: "bool"},
    }
    testTable := NewTable(tm)
    
    testObj := map[string]interface{}{
        "testcolumnint": (int32)(16),
        "testcolumnstr": "tester",
        "testcolumnfloat": (float64)(73.8),
        "testcolumnbool": false,
    }
    testTable.InsertObject(testObj)

    var count int
    testTable.Coll.Query(func (txn *column.Txn) error {
        count = txn.Count()
        return nil
    })
    assert.Equal(t, count, 1)

    var str_res string
    var int_res int32
    var flo_res float64
    var bool_res bool
    testTable.Coll.Query(func (txn *column.Txn) error {
        s_rd := txn.String("testcolumnstr")
        i_rd := txn.Int32("testcolumnint")
        f_rd := txn.Float64("testcolumnfloat")
        b_rd := txn.Bool("testcolumnbool")
    
        return txn.Range(func (i uint32) {
            str_res, _ = s_rd.Get()
            int_res, _ = i_rd.Get()
            flo_res, _ = f_rd.Get()
            bool_res = b_rd.Get()
        })
    })
    assert.Equal(t, str_res, "tester")
    assert.Equal(t, int_res, (int32)(16))
    assert.Equal(t, flo_res, (float64)(73.8))
    assert.Equal(t, bool_res, false)
}

func Test_Filtering(t *testing.T) {

    tm := TableMetadata{
        "testcolumnint": ColumnMetadata{Type: "int"},
        "testcolumnstr": ColumnMetadata{Type: "string"},
        "testcolumnfloat": ColumnMetadata{Type: "float"},
        "testcolumnbool": ColumnMetadata{Type: "bool"},
    }
    testTable := NewTable(tm)

    for i := 0; i < 20; i++ {
        for j := 0; j < 20; j++ {
            testObj := map[string]interface{}{
                "testcolumnint": (int32)(i),
                "testcolumnstr": "tester",
                "testcolumnfloat": (float64)(j),
                "testcolumnbool": false,
            }
            testTable.InsertObject(testObj)
        }
    }

    filters := []filter{
        IntFilter{
            Col: "testcolumnint",
            Op: "<=",
            Val: (int32)(6),
        },
        FloatFilter{
            Col: "testcolumnfloat",
            Op: ">",
            Val: (float64)(10.5),
        },
        StringFilter{
            Col: "testcolumnstr",
            Op: "=",
            Val: "tester",
        },
        BoolFilter{
            Col: "testcolumnbool",
            Op: "=",
            Val: false,
        },
    }

    result, _ := testTable.Select(
        []string{"testcolumnint", "testcolumnfloat"},
        filters,
    )

    fmt.Println(result) // print on fail
    assert.Equal(t, 63, len(result))
}

func Test_Select(t *testing.T) {

    tm := TableMetadata{
        "testcolumnint": ColumnMetadata{Type: "int"},
        "testcolumnstr": ColumnMetadata{Type: "string"},
        "testcolumnfloat": ColumnMetadata{Type: "float"},
        "testcolumnbool": ColumnMetadata{Type: "bool"},
    }
    testTable := NewTable(tm)
    
    for i := 0; i < 20; i++ {
        for j := 0; j < 20; j++ {
            testObj := map[string]interface{}{
                "testcolumnint": (int32)(i),
                "testcolumnstr": "tester",
                "testcolumnfloat": (float64)(j),
                "testcolumnbool": false,
            }
            testTable.InsertObject(testObj)
        }
    }
    
    filters := []filter{
        IntFilter{
            Col: "testcolumnint",
            Op: "=",
            Val: (int32)(13),
        },
        StringFilter{
            Col: "testcolumnstr",
            Op: "=",
            Val: "tester",
        },
    }
    result, _ := testTable.Select(
        []string{"testcolumnint", "testcolumnstr", "testcolumnfloat"},
        filters,
    )
    
    for _, resObj := range result {
        val, _ := resObj["testcolumnint"]
        assert.Equal(t, (int32)(13), val.(int32))
    }

    assert.Equal(t, 20, len(result))
}
