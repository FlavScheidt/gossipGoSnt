package main

import (
    "context"
    // "fmt"
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
    Validator_Key   []byte
    Hash            string
    // SenderID        string
    // SenderName      string
}

// Subscribe tries to subscribe to the PubSub topic for the room name, returning
// a ChatRoom on success.
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

// Publish sends a message to the pubsub topic.
func (cr *Topic) Publish(message []byte, validatorKey []byte, hash string) error {
    m := Message{
        Message:        message,
        Validator_Key:  validatorKey,
        Hash:           hash,
        // SenderID:       cr.self.Pretty(),
        // SenderName:     cr.name,
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

        //Send to rippled
        response, err := gRPCclient.ToRippled(cr.ctx, &pb.Gossip{Message: m.Message, Validator_Key: m.Validator_Key, Hash: m.Hash})
        if err != nil {
            log.Fatalf("Error when calling ToRippled: %s", err)
        }
        // Log format is "time | handler | received/sent | orign/destination | data"
        log.Println("___________________________________________")
        log.Printf("| gRPC-Client | Sent | Rippled | %v | %v \n", m.Hash, m.Validator_Key)
        log.Printf("Response from server: %s", response.Stream)
    }
}

func topicName(peerName string) string {
    return "validator:" + peerName
}
