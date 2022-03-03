/*
    tests for unmarshal (duh)
*/
package main

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestR_string(t *testing.T) {
    data := []byte{114,115,30,40,50}
    target := []byte{114,115}
    testIndex := 0

    func_ret := r_string(data, 2, &testIndex)
    assert.Equal(t, target, func_ret)
    assert.Equal(t, testIndex, 2)
}

func TestR_long(t *testing.T) {
    data := []byte{240,248,0,0,5,6,7}
    target := (int32)(63728)
    testIndex := 0

    func_ret := r_long(data, &testIndex)
    assert.Equal(t, target, func_ret)
    assert.Equal(t, testIndex, 4)
}

func TestR_float(t *testing.T) {
    data := []byte{0,0,0,0,0,0,4,64}
    target := (float64)(2.5)
    testIndex := 0

    func_ret := r_float(data, &testIndex)
    assert.Equal(t, target, func_ret)
    assert.Equal(t, testIndex, 8)

    data_bigger := []byte{35,190,19,51,40,22,234,64}
    target = (float64)(53425.256235)
    testIndex = 0

    func_ret = r_float(data_bigger, &testIndex)
    assert.Equal(t, target, func_ret)
    assert.Equal(t, testIndex, 8)
}

func Test_simple_tuple(t *testing.T) {
    data := []byte{169,3,84,231,51,51,51,51,51,51,243,63,218,4,116,101,115,116}
    target := []interface{}{true, (float64)(1.2), "test"}

    func_ret, _ := r_object(data, nil)
    for ret_pos, _ := range func_ret {
        assert.Equal(t, func_ret[ret_pos], target[ret_pos])
    }
}

func Test_simple_list(t *testing.T) {
    data := []byte{91,3,0,0,0,84,231,51,51,51,51,51,51,243,63,218,4,116,101,115,116}
    target := []interface{}{true, (float64)(1.2), "test"}

    func_ret, _ := r_object(data, nil)
    for ret_pos, _ := range func_ret {
        assert.Equal(t, func_ret[ret_pos], target[ret_pos])
    }
}

func Test_simple_composite(t *testing.T) {
    data := []byte{91,2,0,0,0,169,3,84,231,51,51,51,51,51,51,243,63,218,4,116,101,115,116,91,3,0,0,0,70,233,34,0,0,0,218,3,112,108,122}

    target := []interface{}{
        []interface{}{true, (float64)(1.2), "test"},
        []interface{}{false, (int32)(34), "plz"},
    }

    func_ret, _ := r_object(data, nil)
    assert.Equal(t, 2, len(func_ret))
    for ret_pos, _ := range func_ret {
        assert.Equal(t, func_ret[ret_pos], target[ret_pos])
    }
}
