/*
    Storage tests
*/
package main

import (
    "testing"

    // "github.com/stretchr/testify/assert"
)

func Test_AddColumn(t *testing.T) {

    test_store := NewStore()

    test_store.NewCollection("test_coll")
    

    expectedNil := test_store.AddColumn("test_coll", "testcolumn", "int")
    if expectedNil != nil {
        t.Errorf("Column type error expected, %v returned", expectedNil.Error())
    }

    
    expectedErr := test_store.AddColumn("test_coll", "testcolumn_fail", "faketype")
    if expectedErr.Error() != "Column type does not exist" {
        t.Errorf("Column type error expected, %v returned", expectedErr.Error())
    }

    expectedErr = test_store.AddColumn("test_coll", "testcolumn", "string")
    if expectedErr.Error() != "Column name already exists" {
        t.Errorf("Column name already exists expected, %v returned", expectedErr.Error())
    }

}

//func Test_AddObject(t *testing.T) {

//}

//func Test_Select(t *testing.T) {

//}
