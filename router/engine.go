/*
    Routing Engine
    The main thread for processing distributed inserts & selects.
    Manages partition placement & creation.
*/

package main

import (
)

type RoutingEngine struct {
}

func NewRoutingEngine() *RoutingEngine{
    return &RoutingEngine{}
}
