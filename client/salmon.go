/*
salmon.go

go interface to cache
*/

package main

import (
    // "fmt"
    // "errors"
)

type blank struct {}

type ColumnMetadata struct {
    Type string `json:"type"`
}
type TableMetadata map[string]ColumnMetadata

type Salmon struct {
    ManagerThread  *Manager
    CommsThread    *CommsManager
    ManagerChannel chan Command
    CommsChannel   chan Command
    FinishChannel  chan blank
}

func NewSalmon() *Salmon {
    mc  := make(chan Command)
    cmc := make(chan Command)
    fc  := make(chan blank)
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
        FinishChannel: fc,
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

    go sal.CommsThread.Start(sal.FinishChannel)
    go sal.ManagerThread.Start(sal.FinishChannel)

    return nil
}

// Close client (can be deferred)
func (sal *Salmon) Close() {
    sal.FinishChannel <- blank{}
}

// Insert an object into the system
func (sal *Salmon) Insert(table string, object Object) error {
    
    // create result channel
    rc := make(chan CommandResult)

    // validate object
    cmd := InsertCommand{
        TableName: table,
        Obj: object,
        ResultChan: rc,
    }

    // send to manager channel
    sal.ManagerChannel <- cmd

    // wait for result
    results, _ := <- rc

    return results.Error
}

// Select queries a table in a SQL-ish fashion
func (sal *Salmon) Select(table string, selectors []string, filters []filter) ([]Object, error) {
   
    // result channel
    rc := make(chan CommandResult)

    // create command
    cmd := SelectCommand{
        TableName: table,
        Selectors: selectors,
        Filters: filters,
        ResultChan: rc,
    }

    // send to manchan
    sal.ManagerChannel <- cmd

    // wait for results
    results, _ := <- rc

    if results.Error != nil {
        return results.Objects, nil
    }
    return nil, results.Error

}
