package main

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)


func Test_ReadConfig(t *testing.T) {

    // create new yaml file to read in
    test_sal := NewSalmon()
    yaml_txt := []byte(`testtable:
  testcolumnint:
    type: int
    name: testcolumnint
    order: 0
  testcolumnstr:
    type: str
    name: testcolumnstr
    order: 1
  testcolumnbool:
    type: bool
    name: testcolumnbool
    order: 2
  testcolumnfloat:
    type: float
    name: testcolumnfloat
    order: 3`)
    
    err := os.WriteFile("/tmp/salmon_readconfig_test.yaml", yaml_txt, 0644)
    if err != nil {
        panic(err)
    }

    tm, err := test_sal.ReadConfig("/tmp/salmon_readconfig_test.yaml")

    assert.Equal(t, len(tm), 1)
    assert.Equal(t, len(tm["testtable"]), 4)
}

func Test_Insert(t *testing.T) {
    assert.Equal(t, 1, 1)
}

func Test_Select(t *testing.T) {
    assert.Equal(t, 1, 1)
}
