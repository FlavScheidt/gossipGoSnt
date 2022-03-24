package main

import (
    "context"
    // "encoding/binary"
    // "flag"
    // "fmt"
    // "os"
    // "os/signal"
    // "syscall"
    // "time"

    // "github.com/libp2p/go-libp2p"
    // "github.com/libp2p/go-libp2p-core/host"
    // "github.com/libp2p/go-libp2p-core/network"
    // "github.com/libp2p/go-libp2p-core/peer"
    // // "github.com/libp2p/go-libp2p-discovery"
    // "github.com/multiformats/go-multiaddr"

    // "crypto/rand"
    // "bufio"
    "io"
    "log"
    mrand "math/rand"

    // "github.com/libp2p/go-libp2p-core/crypto"
    // "github.com/libp2p/go-libp2p-core/peerstore"
    // pubsub "github.com/libp2p/go-libp2p-pubsub"

)


func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Setting the TCP port as 0 makes libp2p choose an available port for us.
    // You could, of course, specify one if you like.
    sourcePort := 45511

    //Get node data
    thisNode := newNode()

    var r io.Reader

    //Use the name of the node to derive the key. Always derive the same key
    r = mrand.New(mrand.NewSource(int64(thisNode.id)))

    //Creates the new libp2p
    host, err := makeHost(sourcePort, r)
    if err != nil {
        log.Println(err)
        return
    }

    //First we try to connect to everyone on the list
    for i := 0; i<len(thisNode.peersList); i++ {
        peerAddr := getPeerMultAddr(thisNode.peersList[i])

        rw, err := startPeerAndConnect(ctx, host, peerAddr)
        if err != nil {
           log.Println("Peer is not online... Next one.")
        } else {
            go writeData(rw)
            go readData(rw)
        }
    }

    //Now we wait for incoming connections
    startPeer(ctx, host, handleStream)


    // create a new PubSub service using the GossipSub router
    // ps, err := pubsub.NewGossipSub(ctx, host)
    // if err != nil {
    //     panic(err)
    // }

    // Wait forever
    select {}

}