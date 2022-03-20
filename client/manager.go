package main

import (
    // "fmt"
)

type ManagerOptions struct {
}

type Manager struct {
    Tables map[string]*Table
}

func NewManager(mo ManagerOptions) *Manager {
    return &Manager{
        Tables: make(map[string]*Table),
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
    return
}
