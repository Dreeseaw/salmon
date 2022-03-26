package main

import (
)

type Command interface {}

type Object map[string]interface{}

type InsertCommand struct {
    TableName  string
    Obj        Object
    ResultChan chan CommandResult
}

type SelectCommand struct {
    TableName  string
    Selectors  []string
    Filters    []filter
    ResultChan chan CommandResult
}

type CommandResult struct {
    Error   error
    Objects []Object
}
