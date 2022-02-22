/*
    tests for unmarshal (duh)
*/
package main

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestR_string(t *testing.T) {
    data := []byte{10,20,30,40,50}
    target := []byte{10,20}
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
