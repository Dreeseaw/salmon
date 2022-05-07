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

func (cm *ClientMap) GetMany(ids []string) ([]*Client, error) {
    clis := make([]*Client, 0)
    for _, idn := range ids {
        cli, err := cm.Get(idn)
        if err != nil {
            return nil, err 
        }
        clis = append(clis, cli)
    }
    return clis, nil
}

func (cm *ClientMap) GetAll() []*Client {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    retter := make([]*Client, 0)
    for _, cli := range cm.clis {
        retter = append(retter, cli)
    }

    return retter
}
