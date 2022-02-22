/*
    main code for unmarshalling
    python3's marshal.dumps()
*/
package main

import (
    "fmt"
    "strconv"
)

var (
    TUPLE_TYPE_BYTE uint8 = 169

)

// cut and return 1 byte from p
func r_byte(p []byte, index *int) byte {
    ret := p[*index]
    *index += 1
    return ret
}

// cut and return str_len bytes from p
func r_string(p []byte, str_len int, index *int) []byte {
    ret := p[*index:*index+str_len]
    *index += str_len
    return ret
}

// cut and return int32 from p
func r_long(p []byte, index *int) int32 {
    long_bytes := r_string(p, 4, index)
    var ret int32 = -1
    ret = (int32)(long_bytes[0])
    ret = ret | ((int32)(long_bytes[1]) << 8)
    ret = ret | ((int32)(long_bytes[2]) << 16)
    ret = ret | ((int32)(long_bytes[3]) << 24)
    return ret
}

// cut and return float64 from p
func r_float(p []byte, index *int) float64 {
    float_bytes := r_string(p, 8, index)
    var ret float64
    ret, _ = strconv.ParseFloat(string(float_bytes), 64)
    return ret
}

// given payload, return a mock python tuple
func unmarshal_tuple(payload []byte) ([]interface{}, error) {

    curIndex := 0

    // read 'tuple' type to validate
    tuple_type_byte := r_byte(payload, &curIndex)
    if (uint8)(tuple_type_byte) != TUPLE_TYPE_BYTE {
        return nil, fmt.Errorf("tuple type incorrect: %v", (uint8)(tuple_type_byte))
    }

    // read tuple_len
    tuple_len_byte := r_byte(payload, &curIndex)
    tuple_len := (uint8)(tuple_len_byte)
    if tuple_len == 0 {
        return nil, fmt.Errorf("tuple len is 0")
    }

    tmp_tuple := make([]interface{}, tuple_len)

    for tuple_pos, _ := range tmp_tuple {
        
        obj_type_byte := r_byte(payload, &curIndex)
        obj_type := (uint8)(obj_type_byte)

        if obj_type == 233 {
            // int type
            val := r_long(payload, &curIndex)
            tmp_tuple[tuple_pos] = (int32)(val)
        } else if obj_type == 231 {
            // float type
            val := r_float(payload, &curIndex)
            tmp_tuple[tuple_pos] = (float64)(val)
        } else if obj_type == 218 {
            // string type
            str_len := r_byte(payload, &curIndex)
            val := r_string(payload, (int)(str_len), &curIndex)
            tmp_tuple[tuple_pos] = val
        } else if obj_type == 70 {
            // false bool 
            tmp_tuple[tuple_pos] = false
        } else if obj_type == 84 {
            // true bool
            tmp_tuple[tuple_pos] = true
        } else {
            return nil, fmt.Errorf("type not recognized: %v", (uint8)(obj_type))
        }
    }

    return tmp_tuple, nil
}
