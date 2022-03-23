package main

import (
    "context"
    "encoding/binary"
    // "flag"
    "fmt"
    "os"
    // "os/signal"
    // "syscall"
    "time"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-core/host"
    "github.com/libp2p/go-libp2p-core/network"
    "github.com/libp2p/go-libp2p-core/peer"
    // "github.com/libp2p/go-libp2p-discovery"
    "github.com/multiformats/go-multiaddr"

    // "crypto/rand"
    "bufio"
    "io"
    "log"
    mrand "math/rand"

    "github.com/libp2p/go-libp2p-core/crypto"
    "github.com/libp2p/go-libp2p-core/peerstore"

)

//lotus 12D3KooWBPuBLDxznaw27fk9k9dt2Vp2MSGodUNBz3tP557iWmtQ
//caterham 12D3KooWLr6FwuGrFpQj6VhWjvFpehsrk1yZpud6AVU9QmQ2LdNV

func readData(rw *bufio.ReadWriter) {
    for {
        str, _ := rw.ReadString('\n')

        if str == "" {
            return
        }
        if str != "\n" {
            // Green console colour:    \x1b[32m
            // Reset console colour:    \x1b[0m
            fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
        }

    }
}

func writeData(rw *bufio.ReadWriter) {
    stdReader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("> ")
        sendData, err := stdReader.ReadString('\n')
        if err != nil {
            log.Println(err)
            return
        }

        rw.WriteString(fmt.Sprintf("%s\n", sendData))
        rw.Flush()
    }
}
// const protocolID = "/xrpl/1.0.0"
// const discoveryNamespace = "xrpl"

func startPeer(ctx context.Context, h host.Host, streamHandler network.StreamHandler) {
    // Set a function as stream handler.
    // This function is called when a peer connects, and starts a stream with this protocol.
    // Only applies on the receiving side.
    h.SetStreamHandler("/xrpl/1.0.0", streamHandler)

    // Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
    var port string
    for _, la := range h.Network().ListenAddresses() {
        if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
            port = p
            break
        }
    }

    if port == "" {
        log.Println("was not able to find actual local port")
        return
    }

    log.Println("Waiting for incoming connection")
    log.Println()
}

func handleStream(s network.Stream) {
    log.Println("Got a new stream!")

    // Create a buffer stream for non blocking read and write.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

    go readData(rw)
    go writeData(rw)

    // stream 's' will stay open until you close it (or the other side closes it).
}


func makeHost(port int, seed io.Reader) (host.Host, error) {
    // Creates a new ED25519 key pair for this host.
    // Using ed25519 instead of RSA because RSA implementation in go prevents deterministic behavior
    // We need deterministic behavior, since it is a test enviroment
    // REMEMBER TO USE RSA WITH A RANDOM SEED IN PRODUCTION
    prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, seed)
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


func startPeerAndConnect(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
    log.Println("This node's multiaddresses:")
    for _, la := range h.Addrs() {
        log.Printf(" - %v\n", la)
    }
    log.Println()

    // Turn the destination into a multiaddr.
    maddr, err := multiaddr.NewMultiaddr(destination)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // Extract the peer ID from the multiaddr.
    info, err := peer.AddrInfoFromP2pAddr(maddr)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // Add the destination's peer multiaddress in the peerstore.
    // This will be used during connection and stream creation by libp2p.
    h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

    // Start a stream with the destination.
    // Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
    s, err := h.NewStream(context.Background(), info.ID, "/xrpl/1.0.0")
    if err != nil {
        log.Println(err)
        return nil, err
    }
    log.Println("Established connection to destination")

    // Create a buffered stream so that read and writes are non blocking.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

    return rw, nil
}


func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Setting the TCP port as 0 makes libp2p choose an available port for us.
    // You could, of course, specify one if you like.
    sourcePort := 45511
    nodeIDMachine := 220

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

    peerAddr := "/ip4/191.168.20.19/tcp/45511/p2p/12D3KooWLr6FwuGrFpQj6VhWjvFpehsrk1yZpud6AVU9QmQ2LdNV"

    rw, err := startPeerAndConnect(ctx, host, peerAddr)
    if err != nil {
        startPeer(ctx, host, handleStream)
        // log.Println(err)
        // return
    } else {
        // Create a thread to read and write data.
        go writeData(rw)
        go readData(rw)
    }

    // Wait forever
    select {}

    // sigCh := make(chan os.Signal)
    // signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
    // <-sigCh
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

