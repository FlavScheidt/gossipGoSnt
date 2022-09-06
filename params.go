// From github.com/libp2p/gossipsub-hardening

package main

import (
	// "encoding/json"
	// "fmt"
	// "strconv"
	// "strings"
	"time"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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


func pubsubOptions(cfg NodeConfig) ([]pubsub.Option, error) {
	opts := []pubsub.Option{
		// pubsub.WithEventTracer(cfg.Tracer),
		// pubsub.WithFloodPublish(cfg.FloodPublishing),
		// scoreParamsOption(cfg.PeerScoreParams),
	}

	// if cfg.PeerScoreInspect.Inspect != nil && cfg.PeerScoreInspect.Period != 0 {
	// 	opts = append(opts, pubsub.WithPeerScoreInspect(cfg.PeerScoreInspect.Inspect, cfg.PeerScoreInspect.Period))
	// }

	// if cfg.ValidateQueueSize > 0 {
	// 	opts = append(opts, pubsub.WithValidateQueueSize(cfg.ValidateQueueSize))
	// }

	// if cfg.OutboundQueueSize > 0 {
	// 	opts = append(opts, pubsub.WithPeerOutboundQueueSize(cfg.OutboundQueueSize))
	// }

	// Set the overlay parameters
	if cfg.OverlayParams.d >= 0 {
		pubsub.GossipSubD = cfg.OverlayParams.d
	}
	if cfg.OverlayParams.dlo >= 0 {
		pubsub.GossipSubDlo = cfg.OverlayParams.dlo
	}
	if cfg.OverlayParams.dhi >= 0 {
		pubsub.GossipSubDhi = cfg.OverlayParams.dhi
	}
	if cfg.OverlayParams.dscore >= 0 {
		pubsub.GossipSubDscore = cfg.OverlayParams.dscore
	}
	if cfg.OverlayParams.dlazy >= 0 {
		pubsub.GossipSubDlazy = cfg.OverlayParams.dlazy
	}
	if cfg.OverlayParams.dout >= 0 {
		pubsub.GossipSubDout = cfg.OverlayParams.dout
	}
	if cfg.OverlayParams.gossipFactor > 0 {
		pubsub.GossipSubGossipFactor = cfg.OverlayParams.gossipFactor
	}

	// // set opportunistic graft params
	// if cfg.OpportunisticGraftTicks > 0 {
	// 	pubsub.GossipSubOpportunisticGraftTicks = uint64(cfg.OpportunisticGraftTicks)
	// }

	return opts, nil
}
