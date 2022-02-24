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
