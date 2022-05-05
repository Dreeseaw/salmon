package salmon

import (
    "os"
    "testing"

    "github.com/Dreeseaw/salmon/shared/config"

	"github.com/kelindar/column"
    "github.com/stretchr/testify/assert"
)

const (
    TMP_CONFIG string = "/tmp/salmon.yaml"
)

func write_test_config() {

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

func Test_Insert_NoRouter(t *testing.T) {

    testSal := NewSalmon("mock")
    write_test_config()
    testSal.Init(TMP_CONFIG)

    testSal.Start()

    testObj := map[string]interface{}{
        "testcolumnint": (int32)(16),
        "testcolumnstr": "tester",
        "testcolumnfloat": (float64)(73.8),
        "testcolumnbool": false,
    }

    err := testSal.Insert("testtable", testObj)
    assert.Nil(t, err)

    testTable, exists := testSal.ManagerThread.Tables["testtable"]
    assert.Equal(t, exists, true)

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

func Test_Select_NoRouter(t *testing.T) {
    assert.Equal(t, 1, 1)
}
