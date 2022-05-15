package commands

type filter interface {
    Process(interface{}) bool
    ColName() string
}

// StringFilter processes filters on strings
type StringFilter struct {
    Col, Op  string
    Val string
}

// BoolFilter processes filters on booleans
type BoolFilter struct {
    Col, Op  string
    Val bool
}

// IntFilter processes filters on ints
type IntFilter struct {
    Col, Op  string
    Val int32
}

// FloatFilter processes filters on floats
type FloatFilter struct {
    Col, Op  string
    Val float64
}

// TODO extend string comparing
func (f StringFilter) Process(v interface{}) bool {
    return v.(string) == f.Val
}

func (f BoolFilter) Process(v interface{}) bool {
    return (v == f.Val)
}

func (f IntFilter) Process(v interface{}) bool {
    comp := v.(int32)
    switch f.Op {
    case "=":
        return comp == f.Val
    case ">":
        return comp > f.Val
    case "<":
        return comp < f.Val
    case "<=":
        return comp <= f.Val
    case ">=":
        return comp >= f.Val
    }
    return false
}

func (f FloatFilter) Process(v interface{}) bool {
    comp := v.(float64)
    switch f.Op {
    case "=":
        return comp == f.Val
    case ">":
        return comp > f.Val
    case "<":
        return comp < f.Val
    case "<=":
        return comp <= f.Val
    case ">=":
        return comp >= f.Val
    }
    return false
}

func (f BoolFilter) ColName() string {
    return f.Col
}
func (f StringFilter) ColName() string {
    return f.Col
}
func (f IntFilter) ColName() string {
    return f.Col
}
func (f FloatFilter) ColName() string {
    return f.Col
}
