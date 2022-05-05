package config

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

const (
    TMP_CONFIG string = "/tmp/salmon.yaml"
)

func write_config() {

    if _, err := os.Stat(TMP_CONFIG); err == nil {
        os.Remove(TMP_CONFIG)
    }

    yaml_txt := []byte(`testtable:
  testcolumnint:
    type: int
    name: testcolumnint
    order: 0
  testcolumnstr:
    type: string
    name: testcolumnstr
    order: 1
    pkey: true
  testcolumnbool:
    type: bool
    name: testcolumnbool
    order: 2
  testcolumnfloat:
    type: float
    name: testcolumnfloat
    order: 3`)
    
    err := os.WriteFile(TMP_CONFIG, yaml_txt, 0644)
    if err != nil {
        panic(err)
    }
}

func Test_ReadConfig(t *testing.T) {

    expected := TableMetadata{
        "testcolumnint": ColumnMetadata{
            Type: "int",
            Name: "testcolumnint",
            Order: 0,
            PKey: false,
        },
        "testcolumnstr": ColumnMetadata{
            Type: "string",
            Name: "testcolumnstr",
            Order: 1,
            PKey: true,
        },
        "testcolumnbool": ColumnMetadata{
            Type: "bool",
            Name: "testcolumnbool",
            Order: 2,
            PKey: false,
        },
        "testcolumnfloat": ColumnMetadata{
            Type: "float",
            Name: "testcolumnfloat",
            Order: 3,
            PKey: false,
        },
    }

    // create test config (see above)
    write_config()

    res_tm, err := ReadConfig(TMP_CONFIG)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, len(res_tm))
    assert.Equal(t, expected, res_tm["testtable"])
}
