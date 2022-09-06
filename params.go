// From github.com/libp2p/gossipsub-hardening

package main

import (
	// "encoding/json"
	// "fmt"
	// "strconv"
	// "strings"
	// "time"
)


type HeartbeatParams struct {
	InitialDelay time.Duration
	Interval     time.Duration
}


type OverlayParams struct {
	d            int
	dlo          int
	dhi          int
	dscore       int
	dlazy        int
	dout         int
	gossipFactor float64
}
