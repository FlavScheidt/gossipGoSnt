// From github.com/libp2p/gossipsub-hardening

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/testground/sdk-go/ptypes"
	"github.com/testground/sdk-go/runtime"
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
