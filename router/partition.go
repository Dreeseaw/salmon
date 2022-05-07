/*
    Partition code for router
*/
package main

import (
    "math/rand"
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
    RF         int //replication factor, constant
    PKeys      []PKey
    ClientIds  []string // roughly reflects current ids
    Partitions map[string]*Partition
}

func NewPartitionSet(pks []PKey, ids []string) *PartitionSet {
    return &PartitionSet{
        RF: 2, // hardcode to start dev
        PKeys: pks,
        ClientIds: ids,
        Partitions: make(map[string]*Partition),
    }
}

func (ps *PartitionSet) Process(obj map[string]interface{}, origin string) []string {

    pIds := make([]PartitionID, 0)
    for _, pk := range ps.PKeys {
        val, _ := obj[pk.Col()]
        pIds = append(pIds, pk.Hash(val))
    }
    pId := mergeIds(pIds)

    part, exists := ps.Partitions[string(pId)]
    if !exists {
        // create new partition
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
