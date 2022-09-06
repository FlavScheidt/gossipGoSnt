package main

import (
    "context"
    "io"
    "log"
    "flag"
    "fmt"
    "io/ioutil"
    "strings"
    mrand "math/rand"
    // "log/syslog"
    "os"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "github.com/FlavScheidt/gossipGoSnt/proto"
    pubsub "github.com/FlavScheidt/go-libp2p-pubsub"
)


// -----------------------------------------
//      Define ports
// -----------------------------------------
const gRPCportServer = "50051"
const sourcePort = 45511 //for libp2p

//Global, because I cannot modify the toRippled function
var nodeTopic *Topic
var publishingTopics []*Topic

func main() {

    // -----------------------------------------
    //      Set log file
    // -----------------------------------------
    LOG_FILE := "./log.out"
    // open log file
    logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Panic(err)
    }
    defer logFile.Close()

    mw := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(mw)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    // -----------------------------------------
    //      Create Context
    // -----------------------------------------
    ctx := context.Background()
    log.Println("Starting...")

    // -----------------------------------------
    //      Get the command line arguments
    //          We need to know which kind of experiment to run
    //          Here we just verify if the arguments are correct
    //
    //          We also get parameters for the gossipsub tests
    // -----------------------------------------
    experimentType := flag.String("type", "", "Type of experiment. Default is empty, shuts down")
    flag.Parse()

    log.Println(strings.ToLower(*experimentType))
    switch strings.ToLower(*experimentType) {
        case "general":
            log.Println("Experiment: One topic for everyone")
        case "validator":
            log.Println("Experiment: One topic per validator")
        case "unl":
            log.Println("Experiment: One topic per UNL")
        default:
            log.Println("Experiment type not recognized. Shutting down")
            return
    }

    d := flag.Int("d", 8, "Target peers in the mesh. Default 8")
    dlo := flag.Int("dlo", 6, "Finds more peers bellow this value. Default 6")
    dhi := flag.Int("dhi", 12, "Remove peers above this value. Default 12")
    dscore := flag.Int("dscore", 4, "When prunning the mesh for oversubscription, keep this many highest-scoring peers. Default 4")
    dlazy := flag.Int("dlazy", 8, "Minimum number of peers to gossip to. Default 8")
    dout := flag.Int("dout", 2, "When pruning the mesh for oversubscription, keep this many outbound connected peers. Default 2")
    gossipFactor := flag.Float("gossipFactor", 0.25, "The factor of peers to gossip to during a round. With d_lazy as a minimum. Default 0.25")

    InitialDelay := flag.Float("InitialDelay", 100 * time.Millisecond, "Heatbeat Initial delay. Default 0,1s")
    Interval := flag.Float("Interval", 1 * time.Second, "Heartbeat interval. Default 1s")

    //GS parameters
    op := OverlayParams{
        d:            *d,
        dlo:          *dlo,
        dhi:          *dhi,
        dscore:       *dscore,
        dlazy:        *dlazy,
        dout:         *dout,
        gossipFactor: *gossipFactor,
    }

    hb := HeartbeatParams{
            InitialDelay: *InitialDelay,
            Interval:     *Interval,
    }

    flag.Parse()

    log.Println("EXECUTION INFO")
    log.Println("Experiment type:", strings.ToLower(*experimentType))
    log.Println("----")
      
    log.Println("d = ", d)
    log.Println("dlo = ", dlo)
    log.Println("dhi = ", dhi)
    log.Println("dscore = ", dscore)
    log.Println("dlazy = ", dlazy)
    log.Println("dout = ", dout)
    log.Println("gossipFactor = ", gossipFactor)
    log.Println("InitialDelay = ", InitialDelay)
    log.Println("Interval = ", Interval)


    // -----------------------------------------
    //      Create LibP2P node
    //          Need to do that for every type of experiment anyway
    // -----------------------------------------
    log.Println("------------------------------------------------------------------")
    log.Println("Libp2p Node")
    //Get node data
    thisNode := newNode(strings.ToLower(*experimentType))

    var r io.Reader

    //Use the name of the node to derive the key. Always derive the same key
    r = mrand.New(mrand.NewSource(int64(thisNode.id)))

    //Creates the new libp2p
    host, err := makeHost(sourcePort, r)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("ID: %s", host.ID().Pretty())
    log.Println("Multiaddresses:")
    for _, la := range host.Addrs() {
        log.Printf(" - %v\n", la)
    }

    // -----------------------------------------
    //      Libp2p Connections
    //          general connects statically with everyone in the peers list
    //          the other two make a dynamic discovery in the network using mDNS
    // -----------------------------------------
    // if strings.ToLower(*experimentType) == "general" { //|| strings.ToLower(*experimentType) == "unl" {
    //     //First we try to connect to everyone on the list
    //     for i := 0; i<len(thisNode.peersList); i++ {
    //         log.Println("Calling ", thisNode.peersList[i].name)
    //         peerAddr := getPeerMultAddr(thisNode.peersList[i])

    //         go startPeerAndConnect(ctx, host, peerAddr)
    //         // if err != nil {
    //         //    log.Println("Peer is not online... Next one.")
    //         // } //else {
    //         //     go writeData(rw)
    //         //     go readData(rw)
    //         // }
    //     }

    //     //Now we wait for incoming connections
    //     go startPeer(ctx, host, handleStream)
    // } else {
        log.Println("Finding peers...")
        // setup local mDNS discovery
        if err := setupDiscovery(host); err != nil {
            panic(err)
        }
    // }
    log.Println("------------------------------------------------------------------")

    //Create new GossipSub instance
    tracer, err := pubsub.NewJSONTracer("./trace.json")
    if err != nil {
      panic(err)
    }

    //GossipSub Parameters
    cfg := NodeConfig{
        OverlayParams:           op,
        // Tracer:                  tracer,
        Heartbeat:               hb,
    }

    //GossipSub parameters
    // opts, err := pubsubOptions(cfg)
    // if err != nil {
    //     return nil, err
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

    ps, err := pubsub.NewGossipSub(ctx, host, pubsub.WithEventTracer(tracer))//, opts...)
    if err != nil {
        panic(err)
    }
    log.Println("GossipSub service created")

    // p := PubsubNode{
    //     cfg:      cfg,
    //     ctx:      ctx,
    //     h:        h,
    //     ps:       ps,
    // }
    // log.Println("GossipSub node created")

    // -----------------------------------------
    //      gRPC Client
    // -----------------------------------------
    //Get node ephemeral key generated by rippled
    ephKeyBytes, err := ioutil.ReadFile("/root/sntrippled/my_build/key.out")
    if err != nil {
        log.Fatal(err)
    }
    ephKey := strings.TrimSpace(fmt.Sprintf("%s",ephKeyBytes))
    log.Println("------------------------------------------------------------------")
    log.Println("GRPC Node")
    log.Println("Rippled ephemeral key: ", ephKey)

    var conn *grpc.ClientConn
    conn, err = grpc.Dial("localhost:20052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %s", err)
    }
    defer conn.Close()

    c := pb.NewGossipMessageClient(conn)

    log.Println("Stub Created")
    log.Println("------------------------------------------------------------------")

    // -----------------------------------------
    //      GossipSub
    // -----------------------------------------
    log.Println("------------------------------------------------------------------")
    log.Println("GossipSub")

    //Joining topics
    switch strings.ToLower(*experimentType) {
    case "general":

        //Create peerInfo for everything
        validationsTopic    := peerInfo{name:"validations", id:0, ip: ""}
        // proposalsTopic      = peerInfo{name: "proposals", id: 0, ip: ""}

        nodeTopic, err = Subscribe(ctx, ps, c, host.ID(), validationsTopic)
        publishingTopics = append(publishingTopics, nodeTopic)
        if err != nil {
            panic(err)
        }
        log.Println("Joined topic for ", nodeTopic.name)

    case "validator":
        var topicsList []*Topic
        var topicAux *Topic

        //First, subscrive to nodes own topic
        nodeTopic, err = Subscribe(ctx, ps, c, host.ID(), peerInfo{name: thisNode.name})
        publishingTopics = append(publishingTopics, nodeTopic)
        if err != nil {
            panic(err)
        }
        log.Println("Joined own topic")

        //Subscribe to each topic in the peers list
        for i := 0; i<len(thisNode.peersList); i++ {
            topicAux, err = Subscribe(ctx, ps, c, host.ID(), thisNode.peersList[i])
            topicsList = append(topicsList, topicAux)
            if err != nil {
                panic(err)
            }
            log.Println("Joined topic for ", topicsList[i].name)
        }
    case "unl":
        var topicsList []*Topic
        var topicAux *Topic

        //Subscribe on the topic to listen
        topicAux, err = Subscribe(ctx, ps, c, host.ID(), peerInfo{name: thisNode.unlName})
        topicsList = append(topicsList, topicAux)
        if err != nil {
            panic(err)
        }
        log.Println("Joined topic for ", topicsList[0].name)

        //If this node also publishes to this topic, we add it to the publishing list
        j := 0
        if thisNode.publishSubscribed == true {
            publishingTopics = append(publishingTopics, topicAux)
            log.Println("Also publishes on this topic")
            j++
        }

        //Subscribe to the topic to publish 
        for i := 0; i<len(thisNode.unlPublishing); i++ {
            //first we need to know if we are already subscribed to the topic
            if thisNode.unlPublishing[i] == thisNode.unlName {
                log.Println("Already subscribed to", thisNode.unlPublishing[i])
            } else {
                topicAux, err = SubscribeWithoutReceiving(ctx, ps, c, host.ID(), peerInfo{name:thisNode.unlPublishing[i]})
                publishingTopics = append(publishingTopics, topicAux)
                if err != nil {
                    panic(err)
                }
                log.Println("Joined publishing topic for ", publishingTopics[j].name)
                j++
            }
        }

    }
    log.Println("------------------------------------------------------------------")

    
    
    // -----------------------------------------
    //      gRPC Server
    // -----------------------------------------
    log.Println("------------------------------------------------------------------")
    log.Println("GRPC Server")

    go gRPCserver()

    // Wait forever
    select {}

}
