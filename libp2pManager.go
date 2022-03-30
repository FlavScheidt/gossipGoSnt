package main

import (
    "context"
    "fmt"
    "os"
    "bufio"
    "io"
    "log"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-core/host"
    "github.com/libp2p/go-libp2p-core/network"
    "github.com/libp2p/go-libp2p-core/peer"
    "github.com/multiformats/go-multiaddr"
    "github.com/libp2p/go-libp2p-core/crypto"
    "github.com/libp2p/go-libp2p-core/peerstore"

)

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

//Start peer and wait for icoming connections
//We are not veirfying if the node is on the list or not because we are using a controled test environment
//  This environment has programmed tasks to guarantee that everything will go accordinly
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

    log.Println("Waiting for incoming connections")
    log.Println()
}

//Habndle streams after receiving an incoming transmission
func handleStream(s network.Stream) {
    log.Println("Got a new stream!")

    // Create a buffer stream for non blocking read and write.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

    go readData(rw)
    go writeData(rw)

    // stream 's' will stay open until you close it (or the other side closes it).
}

//Creates the libp2p layer
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

//Start peert and connect to the list =)
func startPeerAndConnect(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
    // log.Println("This node's multiaddresses:")
    // for _, la := range h.Addrs() {
    //     log.Printf(" - %v\n", la)
    // }
    // log.Println()

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