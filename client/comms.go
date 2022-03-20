/*
comms.go

This thread sends & recieves commands to/from the router
No other thread should open ports or reqs
*/

package main

import (
    "io/ioutil"
    "net/http"
    "encoding/json"
)

type CommsManagerOptions struct {
    RouterAddr string
}

type CommsManager struct {
    RouterAddr string
}

func NewCommsManager(cmo CommsManagerOptions) *CommsManager {
    return &CommsManager{
        RouterAddr: cmo.RouterAddr,
    }
}

// Init comms manager, decode tables
func (cm *CommsManager) Init() (map[string]TableMetadata, error) {
    
    acceptRoute := cm.RouterAddr+"/accept"
    resp, err := http.Get(acceptRoute)
    if err != nil {
        return nil, err
    }

    tableJson, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    tableData := make(map[string]TableMetadata)
    if err = json.Unmarshal(tableJson, &tableData); err != nil {
        return nil, err
    }

    return tableData, nil
}

// Start the communication manager
func (cm *CommsManager) Start() {
    return
}
