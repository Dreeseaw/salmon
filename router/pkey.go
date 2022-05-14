/*
    PKey code for each type of column

    contains info on a table's current ditribution wrt the pkey
*/
package main

import (
    "math"
    "encoding/binary"
)

type PKey interface {
    Col() string
    Hash(interface{}) PartitionID
}

func NewPKey(name, ctype string) PKey {
    switch ctype {
    case "int":
        return &IntKey{
            col: name,
            Seen: make([]int32, 0),
        }
    case "string":
        return &StringKey{
            col: name,
            Seen: make([]string, 0),
        }
    case "bool":
        return &BoolKey{
            col: name,
        }
    case "float":
        return &FloatKey{
            col: name,
            Seen: make([]float64, 0),
        }
    }
    return nil
}

type IntKey struct {
    col  string
    Seen []int32
}

func (ik *IntKey) Hash(inp interface{}) PartitionID {
    bs := make([]byte, 4)
    binary.BigEndian.PutUint32(bs, uint32(inp.(int32)))
    return PartitionID(bs)
}

func (ik *IntKey) Col() string {
    return ik.col
}

type StringKey struct {
    col  string
    Seen []string
}

func (sk *StringKey) Hash(inp interface{}) PartitionID {
    bs := []byte(inp.(string))
    return PartitionID(bs)
}

func (sk *StringKey) Col() string {
    return sk.col
}

type BoolKey struct {
    col  string
}

func (bk *BoolKey) Hash(inp interface{}) PartitionID {
    bs := make([]byte, 1)
    if inp.(bool) {
        bs[0] = byte(int8(1))
    } else {
        bs[0] = byte(int8(0))
    }
    return PartitionID(bs)
}

func (bk *BoolKey) Col() string {
    return bk.col
}

type FloatKey struct {
    col  string
    Seen []float64
}

func (fk *FloatKey) Hash(inp interface{}) PartitionID {
    bits := math.Float64bits(inp.(float64))
    var bs []byte
    binary.BigEndian.PutUint64(bs, bits)
    return PartitionID(bs)
}

func (fk *FloatKey) Col() string {
    return fk.col
}
