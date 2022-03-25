/*
   main code for marshalling
   mocks python3's marshal.loads()
*/
package main

import (
	"fmt"
	// "log"
	// "encoding/binary"
	// "math"
)


// w_long returns []byte representing int32 
func w_long(value int32) []byte {
	long_bytes := make([]byte, 4)
    long_bytes[0] = (byte)(value     & 0xff)
    long_bytes[1] = (byte)(value>> 8 & 0xff)
    long_bytes[2] = (byte)(value>>16 & 0xff)
    long_bytes[3] = (byte)(value>>24 & 0xff)
    return long_bytes
}


// w_float returns []byte representing float64
func w_float(value float64) []byte {
	float_bytes := make([]byte, 8)
    
    return float_bytes
	// bits := binary.LittleEndian.Uint64(float_bytes)
	// return math.Float64frombits(bits)
}


// w_object returns payload ([]byte) given mock python object
func w_object(object []interface{}) ([]byte, error) {

	// for now, just focus on a list of tuples
    inner_obj_ct := len(object)
    inner_obj_len := len(object[0].([]interface{}))

    // make payload to start writing to
    payload := make([]byte, 0)

    // set object type byte
    payload = append(payload, (byte)(LIST_TYPE_BYTE))
    payload = append(payload, w_long((int32)(inner_obj_ct))...)

    for _, inner_obj := range object {
        tmp_obj := inner_obj.([]interface{})
        payload = append(payload, (byte)(TUPLE_TYPE_BYTE))
        payload = append(payload, (byte)(inner_obj_len))

        for _, inner_field := range tmp_obj {
            switch val := inner_field.(type) {
            case int32:
                payload = append(payload, (byte)(INT_TYPE_BYTE))
                payload = append(payload, w_long(val)...)
                fmt.Println("int: ", val)
            case float64:
                fmt.Println("float: ", val)
            case bool:
                if val {
                    payload = append(payload, (byte)(TRUE_TYPE_BYTE))
                } else {
                    payload = append(payload, (byte)(FALSE_TYPE_BYTE))
                }  
                fmt.Println("bool: ", val)
            case string:
                payload = append(payload, (byte)(STR_TYPE_BYTE))
                payload = append(payload, ([]byte)(val)...)
                fmt.Println("string: ", val)
            default:
                fmt.Println("unknown")
            }
        }
    }

    return payload, nil
}
