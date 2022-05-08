/*
    Partition code for router
*/
package main

import (
    "fmt"
    "math/rand"
    "sync"
)

type PartitionID []byte

func mergeIds(inp []PartitionID) PartitionID {
    bs := make([]byte, 0)
    for _, pid := range inp {
        for _, b := range []byte(pid) {
            bs = append(bs, b)
        }
    }
    return PartitionID(bs)
}

// PartitionSet stores all relevant partition data
// for a particular table & holds access locks
type PartitionSet struct {
    cmux       sync.RWMutex
    RF         int //replication factor, constant
    PKeys      []PKey
    ClientIds  []string // roughly reflects current ids
    Partitions map[string]*Partition
    version    int
}

func NewPartitionSet(pks []PKey, ids []string) *PartitionSet {
    return &PartitionSet{
        RF: 2, // hardcode to start dev
        PKeys: pks,
        ClientIds: ids,
        Partitions: make(map[string]*Partition),
        version: 0,
    }
}

// UpdateClients syncs a pset's client list to the current cmap
// and issues commands for adding & removing a node
// TODO: run as background thread that reads changes from cmap
func (ps *PartitionSet) UpdateClients(clis *ClientMap) {
    ps.cmux.Lock()
    defer ps.cmux.Unlock()

    ces := clis.GetEvents(ps.version)

    // TODO: collapse new events
    // ie: a client connects & disconnects quickly (errant startup)
    // no need to rebalance twice
    // or if a client fails & recovers between system inserts

    for _, cev := range ces {
        // apply each event in order
        if cev.Event {
            ps.ClientIds = append(ps.ClientIds, cev.CId)
            // TODO: rebalance partitions
        } else {
            for i, cId := range ps.ClientIds {
                if cId == cev.CId {
                    ps.ClientIds = append(ps.ClientIds[:i], ps.ClientIds[i+1:]...)
                    break
                }
            }
            // TODO: rebalance partitions
        }
        ps.version = ps.version + 1
    }

    return
}

func (ps *PartitionSet) Process(obj map[string]interface{}, origin string) []string {

    pIds := make([]PartitionID, 0)
    for i, pk := range ps.PKeys {
        val, _ := obj[pk.Col()]
        hash := pk.Hash(val)
        fmt.Printf("insert from: %v, pkey %v hashed as %v", origin, i, string(hash))
        pIds = append(pIds, hash)
    }
    pId := mergeIds(pIds)

    part, exists := ps.Partitions[string(pId)]
    if !exists {
        // create new partition
        fmt.Printf("[server] creating new partition %v", string(pId))
        part = NewPartition()
        ps.Partitions[string(pId)] = part
        part.Owners = append(part.Owners, origin)

        // for now, add N-1 random other clients
        // in the future, this is where I can start maximizing
        // data locality (store client stats, etc)
        otherOwners := ps.GetRandomOwners(ps.RF-1, origin)
        part.Owners = append(part.Owners, otherOwners)
    }

    return part.Owners
}

type Partition struct {
    Owners []string // client ids of current partition owners
}

func NewPartition() *Partition {
    return &Partition{
        Owners: make([]string, 0),
    }
}

func (ps *PartitionSet) GetRandomOwners(amt int, exclude string) string {
    avail := make([]string, 0)
    for _, id := range ps.ClientIds {
        if id != exclude {
            avail = append(avail, id)
        }
    }
    randN := rand.Intn(len(avail))
    return avail[randN]
}
