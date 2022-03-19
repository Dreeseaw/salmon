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

func Test_AddColumn(t *testing.T) {

    test_store := NewStore()
    test_store.NewCollection("test_coll")
    
    expectedNil := test_store.AddColumn("test_coll", "testcolumn", "int")
    if expectedNil != nil {
        t.Errorf("Nil expected, %v returned", expectedNil.Error())
    }

    expectedErr := test_store.AddColumn("test_coll", "testcolumn_fail", "faketype")
    if expectedErr.Error() != "Column type does not exist" {
        t.Errorf("Column type error expected, %v returned", expectedErr.Error())
    }

    expectedErr = test_store.AddColumn("test_coll", "testcolumn", "string")
    if expectedErr.Error() != "Column name already exists" {
        t.Errorf("Column name already exists expected, %v returned", expectedErr.Error())
    }

    cmm, _ := test_store.CollMetadataMap["test_coll"]
    assert.Equal(t, cmm[0].Name, "testcolumn")
    assert.Equal(t, cmm[0].Type, "int")
}

func Test_AddObject(t *testing.T) {
    test_store := NewStore()
    test_store.NewCollection("test_coll")
    
    _ = test_store.AddColumn("test_coll", "testcolumnint", "int")
    _ = test_store.AddColumn("test_coll", "testcolumnstr", "string")
    _ = test_store.AddColumn("test_coll", "testcolumnfloat", "float")
    _ = test_store.AddColumn("test_coll", "testcolumnbool", "bool")
    
    test_obj := []interface{}{(int32)(16),"tester",(float64)(73.8),false}
    test_store.AddObject("test_coll", test_obj)

    coll, _ := test_store.CollMap["test_coll"]
    var count int
    coll.Query(func (txn *column.Txn) error {
        count = txn.Count()
        return nil
    })
    assert.Equal(t, count, 1)

    var str_res string
    var int_res int32
    var flo_res float64
    var bool_res bool
    coll.Query(func (txn *column.Txn) error {
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

    test_store := NewStore()
    test_store.NewCollection("test_coll")
    
    _ = test_store.AddColumn("test_coll", "testcolumnint", "int")
    _ = test_store.AddColumn("test_coll", "testcolumnstr", "string")
    _ = test_store.AddColumn("test_coll", "testcolumnfloat", "float")
    err := test_store.AddColumn("test_coll", "testcolumnbool", "bool")
    fmt.Println(err)

    for i := 0; i < 20; i++ {
        for j := 0; j < 20; j++ {
            // obj_bool := (0 == (i*j % 2))
            test_obj := []interface{}{(int32)(i),"tester",(float64)(j),false}
            test_store.AddObject("test_coll", test_obj)
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

    result, _ := test_store.Select(
        "test_coll", 
        []string{"testcolumnint", "testcolumnfloat"},
        filters,
    )
 
    /*
    for _, res_obj := range result {
        val, _ := res_obj["testcolumnint"]
        assert.Equal(t, (int32)(13), val.(int32))
    }
    */

    fmt.Println(result) // print on fail
    assert.Equal(t, 63, len(result))
}

func Test_Select(t *testing.T) {

    test_store := NewStore()
    test_store.NewCollection("test_coll")
    
    _ = test_store.AddColumn("test_coll", "testcolumnint", "int")
    _ = test_store.AddColumn("test_coll", "testcolumnstr", "string")
    _ = test_store.AddColumn("test_coll", "testcolumnfloat", "float")
    _ = test_store.AddColumn("test_coll", "testcolumnbool", "bool")
    
    for i := 0; i < 20; i++ {
        for j := 0; j < 20; j++ {
            test_obj := []interface{}{(int32)(i),"tester",(float64)(j),false}
            test_store.AddObject("test_coll", test_obj)
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
    result, _ := test_store.Select(
        "test_coll", 
        []string{"testcolumnint", "testcolumnstr", "testcolumnfloat"},
        filters,
    )
    
    for _, res_obj := range result {
        val, _ := res_obj["testcolumnint"]
        assert.Equal(t, (int32)(13), val.(int32))
    }

    assert.Equal(t, 20, len(result))
}
