package main

import (
    "sync"
    "errors"
)

type Client struct {
    Id       string
    ReplChan InsertCommChan
}

type ClientMap struct {
    mu   sync.RWMutex
    clis map[string]*Client
}

func NewClientMap() *ClientMap {
    return &ClientMap{
        clis: make(map[string]*Client),
    }
}

func (cm *ClientMap) Add(cli *Client) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if _, exists := cm.clis[cli.Id]; exists {
        return errors.New("client map id collision")
    }
    cm.clis[cli.Id] = cli
    return nil
}

func (cm *ClientMap) Remove(cliID string) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if _, exists := cm.clis[cliID]; !exists {
        return errors.New("client map can't delete non-existent id")
    }
    delete(cm.clis, cliID)
    return nil
}

func (cm *ClientMap) Get(cliID string) (*Client, error) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    cli, exists := cm.clis[cliID]
    if !exists {
        return nil, errors.New("client map can't find id")
    }
    return cli, nil
}
