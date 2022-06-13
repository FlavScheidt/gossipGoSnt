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
    unlName string
    unlPublishing []string
    publishSubscribed bool //indicates if the node is subscribed to the same topic in which it publishes
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func newNode(experiment string)(nodeInfo) {
    var thisNode nodeInfo;

    //get node name from /etc/hostname
    nodeName, err := ioutil.ReadFile("/etc/hostname")
    if err != nil {
        log.Fatal(err)
    }
    thisNode.name =  strings.TrimSpace(fmt.Sprintf("%s",nodeName))

    //get node id from clusterConfig.csv
    clusterConfig, err := os.Open("./clusterConfig.csv")
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
            thisNode.unlName = dataExtract[3]
            thisNode.publishSubscribed = false
        }
    }

    //Creates list of peers
    if experiment == "validator" {
        unlFile := thisNode.name 

        unl, err := os.Open("/root/gossipGoSnt/clusterConfig/validator/"+unlFile+".txt")
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
            clusterConfig, err := os.Open("./clusterConfig.csv")
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
    } else {
        //List of topics to publish in
        //List files on the directory
        files, err := ioutil.ReadDir("/root/gossipGoSnt/clusterConfig/unl/")
        if err != nil {
            log.Fatal(err)
        }

        //iterate over files until this node is found or EOF
        for _, file := range files {
            fileName := file.Name()
            unl, err := os.Open("/root/gossipGoSnt/clusterConfig/unl/"+fileName)
            if err != nil {
                log.Fatal(err)
            }
            defer unl.Close()
            scanner = bufio.NewScanner(unl)
                if err := scanner.Err(); err != nil {
                log.Fatal(err)
            }

            //if this node is on the file, add the filename to the list
            for scanner.Scan() {
                if strings.TrimSpace(scanner.Text()) == thisNode.name {
                    thisNode.unlPublishing = append(thisNode.unlPublishing, fileName[:len(fileName)-4])
                    //If it is the UNL the node is subscribed for, we sinalize for future use
                    if fileName[:len(fileName)-4] == thisNode.unlName {
                        thisNode.publishSubscribed = true
                        log.Println("INFO: Listens to the same topic it publishes")
                    }
                }
            }

        }
        log.Println("Publishing list size: ", len(thisNode.unlPublishing))
       
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