package main

import (
    "fmt"
    "log"
    "net"
    "context"
    // "time"

    "google.golang.org/grpc"
    // "google.golang.org/grpc/credentials/insecure"
    pb "github.com/libp2p/gossipGoSnt/proto"

    // proto "github.com/golang/protobuf/proto"
    // protoreflect "google.golang.org/protobuf/reflect/protoreflect"
    // protoimpl "google.golang.org/protobuf/runtime/protoimpl"
    // reflect "reflect"
    // sync "sync"
)


// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedGossipMessageServer
}


func (s *server) ToLibP2P(ctx context.Context, in *pb.Gossip) (*pb.Control, error) {
    // Log format is "time | handler | received/sent | orign/destination | data"
    log.Printf("| gRPC-Server | Received | Rippled | %v | %v \n", in.GetHash(), in.GetValidator_Key())

    //Send message to gossipsub
     for i := 0; i<len(publishingTopics); i++ {
        publishingTopics[i].Publish(in.GetMessage(), in.GetValidator_Key(), in.GetHash())   
        log.Println("Message published on topic", publishingTopics[i].name)
    } 

    return &pb.Control{Stream: true}, nil
}


func gRPCserver() {
    log.Println("Starting gRPC server")

    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCportServer))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
        log.Println("------------------------------------------------------------------")
    }
    s := grpc.NewServer()
    pb.RegisterGossipMessageServer(s, &server{})
    log.Printf("server listening at %v", lis.Addr())
    log.Println("------------------------------------------------------------------")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
        log.Println("------------------------------------------------------------------")
    }
}
