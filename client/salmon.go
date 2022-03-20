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
    ManagerThread *Manager
    CommsThread   *CommsManager
}

func NewSalmon() *Salmon {
    man  := NewManager(ManagerOptions{})
    comm := NewCommsManager(CommsManagerOptions{})
    return &Salmon{
        ManagerThread: man,
        CommsThread: comm,
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
func (sal *Salmon) Insert(object map[string]interface{}) error {
    // validate object

    // cache

    return nil
}

func (sal *Salmon) Select(selectors []string, filters []filter) error {
    return nil    
}
