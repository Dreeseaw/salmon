package main

type Command interface {}

type InsertCommand struct {
    TableName string
    Obj       map[string]interface{}
}

type SelectCommand struct {
    TableName string
    Selectors []string
    Filters   []filter
}
