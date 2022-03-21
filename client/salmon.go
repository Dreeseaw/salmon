/*
salmon.go

go interface to cache
*/

package main

import (
    // "fmt"
)

type ColumnMetadata struct {
    Type string `json:"type"`
}
type TableMetadata map[string]ColumnMetadata

type Salmon struct {
    ManagerThread  *Manager
    CommsThread    *CommsManager
    ManagerChannel chan Command
    CommsChannel   chan Command
}

func NewSalmon() *Salmon {
    mc  := make(chan Command)
    cmc := make(chan Command)
    man  := NewManager(ManagerOptions{
        ManChan: mc,
        CommsChan: cmc,
    })
    comm := NewCommsManager(CommsManagerOptions{
        RouterAddr: "http://localhost:1323",
        ManChan: mc,
        CommsChan: cmc,
    })

    return &Salmon{
        ManagerThread: man,
        CommsThread: comm,
        ManagerChannel: mc,
        CommsChannel: cmc,
    }
}

// Init the salmon client
func (sal *Salmon) Init() error {
    tables, err := sal.CommsThread.Init()
    if err != nil {
        return err
    }
    sal.ManagerThread.Init(tables)
    return nil
}

// Start the salmon client
func (sal *Salmon) Start() error {

    go sal.CommsThread.Start()
    go sal.ManagerThread.Start()

    return nil
}

// Insert an object into the system
func (sal *Salmon) Insert(table string, object map[string]interface{}) error {
    // validate object
    cmd := InsertCommand{
        TableName: table,
        Obj: object,
    }

    // send to manager channel
    sal.ManagerChannel <- cmd

    return nil
}

// Select queries a table in a SQL-ish fashion
func (sal *Salmon) Select(table string, selectors []string, filters []filter) error {
   
    // create command
    cmd := SelectCommand{
        TableName: table,
        Selectors: selectors,
        Filters: filters,
    }

    // send to manchan
    sal.ManagerChannel <- cmd

    return nil
}
