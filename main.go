package main

import (
    // "context"
    "encoding/binary"
    // "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-core/host"
    "github.com/libp2p/go-libp2p-core/network"
    // "github.com/libp2p/go-libp2p-core/peer"
    // "github.com/libp2p/go-libp2p-discovery"
    "github.com/multiformats/go-multiaddr"

    // "crypto/rand"
    // "bufio"
    "io"
    "log"
    mrand "math/rand"

    "github.com/libp2p/go-libp2p-core/crypto"
    // "github.com/libp2p/go-libp2p-core/peerstore"

)

const protocolID = "/xrpl/1.0.0"
const discoveryNamespace = "xrpl"

func makeHost(port int, randomness io.Reader) (host.Host, error) {
    // Creates a new RSA key pair for this host.
    prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // 0.0.0.0 will listen on any interface device.
    sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

    // libp2p.New constructs a new libp2p Host.
    // Other options can be added here.
    return libp2p.New(
        libp2p.ListenAddrs(sourceMultiAddr),
        libp2p.Identity(prvKey),
    )
}

func main() {
    // Add -peer-address flag
    // peerAddr := flag.String("peer-address", "", "peer address")
    // flag.Parse()

    // Create the libp2p host.
    //
    // Note that we are explicitly passing the listen address and restricting it to IPv4 over the
    // loopback interface (127.0.0.1).
    //
    // Setting the TCP port as 0 makes libp2p choose an available port for us.
    // You could, of course, specify one if you like.
    sourcePort := 0
    nodeIDMachine := 220
    // host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"), libp2p.Identity())
    // if err != nil {
    //     panic(err)
    // }
    // defer host.Close()
    var r io.Reader


    //Use the name of the node to derive the key. Always derive the same key
    r = mrand.New(mrand.NewSource(int64(nodeIDMachine)))

    //Creates the new libp2p
    host, err := makeHost(sourcePort, r)
    if err != nil {
        log.Println(err)
        return
    }

    // Print this node's addresses and ID
    fmt.Println("Addresses:", host.Addrs())
    fmt.Println("ID:", host.ID())

    // // Setup a stream handler.
    // //
    // // This gets called every time a peer connects and opens a stream to this node.
    // host.SetStreamHandler(protocolID, func(s network.Stream) {
    //     go writeCounter(s)
    //     go readCounter(s)
    // })

    // // Setup peer discovery.
    // // discoveryService, err := discovery.NewMdnsService(
    // //     context.Background(),
    // //     host,
    // //     time.Second,
    // //     discoveryNamespace,
    // // )
    // // if err != nil {
    // //     panic(err)
    // // }
    // // defer discoveryService.Close()

    // // discoveryService.RegisterNotifee(&discoveryNotifee{h: host})
    // peerAddr := "/ip4/192.168.20.19/tcp/45511 "

    // // If we received a peer address, we should connect to it.
    // if *peerAddr != "" {
    //     // Parse the multiaddr string.
    //     peerMA, err := multiaddr.NewMultiaddr(*peerAddr)
    //     if err != nil {
    //         panic(err)
    //     }
    //     peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
    //     if err != nil {
    //         panic(err)
    //     }

    //     // Connect to the node at the given address.
    //     if err := host.Connect(context.Background(), *peerAddrInfo); err != nil {
    //         panic(err)
    //     }
    //     fmt.Println("Connected to", peerAddrInfo.String())

    //     // Open a stream with the given peer.
    //     s, err := host.NewStream(context.Background(), peerAddrInfo.ID, protocolID)
    //     if err != nil {
    //         panic(err)
    //     }

    //     // Start the write and read threads.
    //     go writeCounter(s)
    //     go readCounter(s)
    // }

    sigCh := make(chan os.Signal)
    signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
    <-sigCh
}

func writeCounter(s network.Stream) {
    var counter uint64

    for {
        <-time.After(time.Second)
        counter++

        err := binary.Write(s, binary.BigEndian, counter)
        if err != nil {
            panic(err)
        }
    }
}

func readCounter(s network.Stream) {
    for {
        var counter uint64

        err := binary.Read(s, binary.BigEndian, &counter)
        if err != nil {
            panic(err)
        }

        fmt.Printf("Received %d from %s\n", counter, s.ID())
    }
}

// type discoveryNotifee struct {
//     h host.Host
// }

// func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
//     fmt.Println("found peer", peerInfo.String())
// }

