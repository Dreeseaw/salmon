/*
    Partition code for router
*/
package main

type Partitioner interface {
}

type BasePartitioner struct {
}

func NewBasePartitioner() *BasePartitioner {
    return &BasePartitioner{}
}

func (bp *BasePartitioner) GetReplicas() {
    
}
