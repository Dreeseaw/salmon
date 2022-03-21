package main

import (
    // "fmt"
)

type ManagerOptions struct {
    ManChan chan Command
    CommsChan chan Command
}

type Manager struct {
    Tables    map[string]*Table
    ManChan   chan Command
    CommsChan chan Command
}

func NewManager(mo ManagerOptions) *Manager {
    return &Manager{
        Tables: make(map[string]*Table),
        ManChan: mo.ManChan,
        CommsChan: mo.CommsChan,
    }
}

// Init manager tables
func (m *Manager) Init(tableData map[string]TableMetadata) {
    for tName, tMeta := range tableData {
        table := NewTable(tMeta)
        m.Tables[tName] = table
    }
}

// Start manager
func (m *Manager) Start() {
    for {
        select {
            case cmd := <- m.ManChan:
                m.Process(cmd)
        }
    }

    return
}

// Process insert or select
func (m *Manager) Process(cmd Command) {
    return
}
