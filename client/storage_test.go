/*
    Storage tests
*/
package main

import (
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
    
    test_obj := []interface{}{16,"tester",73.8,false}
    test_store.AddObject("test_coll", test_obj)

    coll, _ := test_store.CollMap["test_coll"]
    var count int
    coll.Query(func (txn *column.Txn) error {
        count = txn.Count()
        return nil
    })
    assert.Equal(t, count, 1)
}

//func Test_Select(t *testing.T) {

//}
