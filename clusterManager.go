package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "strconv"
    "github.com/libp2p/go-libp2p-core/crypto"
    "github.com/libp2p/go-libp2p-core/peer"
    "io"
    mrand "math/rand"
)

type peerInfo struct {
    name string
    ip string
    id int
}

type nodeInfo struct {
    name string
    id int
    peersList []peerInfo
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func newNode()(nodeInfo) {
    var thisNode nodeInfo;

    //get node name from /etc/hostname
    nodeName, err := ioutil.ReadFile("/etc/hostname")
    if err != nil {
        log.Fatal(err)
    }
    thisNode.name =  strings.TrimSpace(fmt.Sprintf("%s",nodeName))

    //get node id from clusterConfig.csv
    clusterConfig, err := os.Open("clusterConfig.csv")
    if err != nil {
        log.Fatal(err)
    }
    // defer clusterConfig.Close()

    //Read line by line to get the info from the csv
    scanner := bufio.NewScanner(clusterConfig)
        if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    searchString := ","+thisNode.name+","
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), searchString) {
            dataExtract := strings.Split(scanner.Text(), ",")
            thisNode.id, err = strconv.Atoi(dataExtract[2])
        }
    }

    //Creates list of peers
    unl, err := os.Open("./clusterConfig/"+thisNode.name+".txt")
    if err != nil {
        log.Fatal(err)
    }
    defer unl.Close()
    scanner = bufio.NewScanner(unl)
        if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    for scanner.Scan() {
        thisNode.peersList = append(thisNode.peersList, peerInfo{strings.TrimSpace(scanner.Text()), "", 0})
    }

    //Inserts info about every peer in the list
    for j := 0; j<len(thisNode.peersList); j++ {
        searchString := ","+thisNode.peersList[j].name+","
        clusterConfig, err := os.Open("clusterConfig.csv")
        if err != nil {
            log.Fatal(err)
        }
        scanner = bufio.NewScanner(clusterConfig)
        if err := scanner.Err(); err != nil {
            log.Fatal(err)
        }

        for scanner.Scan() {
            if strings.Contains(scanner.Text(), searchString) {
                dataExtract := strings.Split(scanner.Text(), ",")
                thisNode.peersList[j].id, err = strconv.Atoi(dataExtract[2])
                thisNode.peersList[j].ip = dataExtract[0]
            }
        }

    }

    return thisNode
}

func getPeerMultAddr(peerElement peerInfo)(string) {
     var seed io.Reader

    //Use the name of the node to derive the key. Always derive the same key
    seed = mrand.New(mrand.NewSource(int64(peerElement.id)))
    prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, seed)
    if err != nil {
        log.Println(err)
        return ""
    }

    peerID, err := peer.IDFromPublicKey(prvKey.GetPublic())
    pID := fmt.Sprintf("%s", peerID)

    mAddr := "/ip4/"+peerElement.ip+"/tcp/45511/p2p/"+pID
    return mAddr
}