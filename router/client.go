package main

import (
    "sync"
    "errors"
)

type Client struct {
    Id       string
    ReplChan InsertCommChan
}

type ClientEvent struct {
    CId   string
    Event bool // true = added, false = deleted
}

type ClientMap struct {
    mu   sync.RWMutex
    clis map[string]*Client
    ces  []ClientEvent
}

func NewClientMap() *ClientMap {
    return &ClientMap{
        clis: make(map[string]*Client),
        ces: make([]ClientEvent, 0),
    }
}

func (cm *ClientMap) Add(cli *Client) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if _, exists := cm.clis[cli.Id]; exists {
        return errors.New("client map id collision")
    }
    cm.clis[cli.Id] = cli
    cm.ces = append(cm.ces, ClientEvent{cli.Id, true})
    return nil
}

func (cm *ClientMap) Remove(cliID string) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if _, exists := cm.clis[cliID]; !exists {
        return errors.New("client map can't delete non-existent id")
    }
    delete(cm.clis, cliID)
    cm.ces = append(cm.ces, ClientEvent{cliID, false})
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

func (cm *ClientMap) GetEvents(ev int) []ClientEvent {
    if ev > len(cm.ces) {
        panic("pset counter > cmap counter")
    }
    if ev == len(cm.ces) {
        return nil
    }
    return cm.ces[ev:]
}

func (cm *ClientMap) Has(cid string) bool {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    _, exists := cm.clis[cid]
    return exists
}
