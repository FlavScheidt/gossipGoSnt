package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "strings"
    // "os"
    "encoding/json"
    "log"

    pubsub "github.com/libp2p/go-libp2p-pubsub"
    "github.com/libp2p/go-libp2p-core/peer"

    pb "github.com/FlavScheidt/gossipGoSnt/proto"
    // "google.golang.org/grpc"
)

const BufSize = 128

// Topic represents a subscription to a single PubSub topic. Messages
// can be published to the topic with validator.Publish, and received
// messages are pushed to the Messages channel.
type Topic struct {
    // Messages is a channel of messages received from other peers in the chat room
    Messages chan *Message

    ctx   context.Context
    ps    *pubsub.PubSub
    topic *pubsub.Topic
    sub   *pubsub.Subscription

    // validatorID     peer.ID
    self            peer.ID
    // validatorKey    string
    // ip              string
    name            string 
}

type Message struct {
    Message         []byte
    Validator_Key   string
    Hash            string
    SenderID        string
    SenderName      string
}

//Extracted from gossipsub-hardening
// type PubsubNode struct {
//     cfg      NodeConfig
//     ctx      context.Context
//     // shutdown func()
//     // runenv   *runtime.RunEnv
//     h        host.Host
//     ps       *pubsub.PubSub

//     lk     sync.RWMutex
//     // topics map[string]*topicState

//     pubwg sync.WaitGroup
// }

type NodeConfig struct {
    // topics to join when node starts
    // Topics []TopicConfig

    // whether we're a publisher or a lurker
    // Publisher bool

    // pubsub event tracer
    Tracer pubsub.EventTracer

    // Test instance identifier
    // Seq int64

    // How long to wait after connecting to bootstrap peers before publishing
    // Warmup time.Duration

    // How long to wait for cooldown
    // Cooldown time.Duration

    // Gossipsub heartbeat params
    // Heartbeat HeartbeatParams

    // whether to flood the network when publishing our own messages.
    // Ignored unless hardening_api build tag is present.
    // FloodPublishing bool

    // Params for peer scoring function. Ignored unless hardening_api build tag is present.
    // PeerScoreParams ScoreParams

    OverlayParams OverlayParams

    // Params for inspecting the scoring values.
    // PeerScoreInspect InspectParams

    // Size of the pubsub validation queue.
    // ValidateQueueSize int

    // Size of the pubsub outbound queue.
    // OutboundQueueSize int

    // Heartbeat tics for opportunistic grafting
    // OpportunisticGraftTicks int
}

func Subscribe(ctx context.Context, ps *pubsub.PubSub, gRPCclient pb.GossipMessageClient, selfID peer.ID, peerTopic peerInfo) (*Topic, error) {
    // join the pubsub topic
    topic, err := ps.Join(topicName(peerTopic.name))
    if err != nil {
        return nil, err
    }

    // and subscribe to it
    sub, err := topic.Subscribe()
    if err != nil {
        return nil, err
    }

    cr := &Topic{
        ctx:            ctx,
        ps:             ps,
        topic:          topic,
        sub:            sub,
        self:           selfID,
        // validatorID:    peerTopic.id,
        // validatorKey:   peer.
        // ip:             peer.ip,
        name:           peerTopic.name, 
        Messages: make(chan *Message, BufSize),
    }

    // start reading messages from the subscription in a loop
    go cr.readLoop(gRPCclient)
    return cr, nil
}

// Subscribe to the topic only for publishing
//Doenst really subscribes
func SubscribeWithoutReceiving(ctx context.Context, ps *pubsub.PubSub, gRPCclient pb.GossipMessageClient, selfID peer.ID, peerTopic peerInfo) (*Topic, error) {
    // join the pubsub topic
    topic, err := ps.Join(topicName(peerTopic.name))
    if err != nil {
        return nil, err
    }

    // and subscribe to it
    sub, err := topic.Subscribe()
    if err != nil {
        return nil, err
    }

    cr := &Topic{
        ctx:            ctx,
        ps:             ps,
        topic:          topic,
        sub:            sub,
        self:           selfID,
        // validatorID:    peerTopic.id,
        // validatorKey:   peer.
        // ip:             peer.ip,
        name:           peerTopic.name, 
        Messages: make(chan *Message, BufSize),
    }

    // start reading messages from the subscription in a loop
    // go cr.readLoop(gRPCclient)
    return cr, nil
}

// Publish sends a message to the pubsub topic.
func (cr *Topic) Publish(message []byte, validatorKey string, hash string) error {
    m := Message{
        Message:        message,
        Validator_Key:  validatorKey,
        Hash:           hash,
        SenderID:       cr.self.Pretty(),
        SenderName:     cr.name,
    }
    msgBytes, err := json.Marshal(m)
    if err != nil {
        return err
    }

    return cr.topic.Publish(cr.ctx, msgBytes)
}

func (cr *Topic) ListPeers() []peer.ID {
    return cr.ps.ListPeers(topicName(cr.name))
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *Topic) readLoop(gRPCclient pb.GossipMessageClient) {
    
    nodeName, err := ioutil.ReadFile("/etc/hostname")
    if err != nil {
        log.Fatal(err)
    }
    node :=  strings.TrimSpace(fmt.Sprintf("%s",nodeName))

    for {
        msg, err := cr.sub.Next(cr.ctx)
        if err != nil {
            close(cr.Messages)
            return
        }
        // only forward messages delivered by others
        if msg.ReceivedFrom == cr.self {
            continue
        }
        cm := new(Message)
        err = json.Unmarshal(msg.Data, cm)
        if err != nil {
            continue
        }
        // send valid messages onto the Messages channel
        cr.Messages <- cm
        m := <-cr.Messages
        // Log format is "time | node name| handler | received/sent | orign/destination | data"
        log.Printf("| %s | GossipSub | Recieved | GossipSub | %v | %v | %v|  %v | %v \n", node, cr.name, msg.ReceivedFrom, m.SenderName, m.Hash, m.Validator_Key)

        //Send to rippled
        _, err = gRPCclient.ToRippled(cr.ctx, &pb.Gossip{Message: m.Message, Validator_Key: m.Validator_Key, Hash: m.Hash})
        if err != nil {
            log.Fatalf("%s Error when calling ToRippled: %s", node, err)
        }
        // Log format is "time | node name | handler | received/sent | orign/destination | data"
        log.Printf(" | %s | gRPC-Client | Sent | Rippled | %v | %v \n", node, m.Hash, m.Validator_Key)
    }
}

func topicName(peerName string) string {
    return "validator:" + peerName
}
