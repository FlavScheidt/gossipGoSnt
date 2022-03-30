package main

import (
    // "bufio"
    // "fmt"
    // "io/ioutil"
    "log"
    // "os"
    // "strings"
    // "io"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "gossip_message"

    proto "github.com/golang/protobuf/proto"
    protoreflect "google.golang.org/protobuf/reflect/protoreflect"
    protoimpl "google.golang.org/protobuf/runtime/protoimpl"
    reflect "reflect"
    sync "sync"
)


// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedGossipMessageServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) ToLibP2P(ctx context.Context, in *pb.Gossip) (*pb.Control, error) {
    //NEED TO KNOW HOW TO GET THJE DATE USING GO
    //I STOPPED HERE
    log.Println("___________________________________________")
    log.Printf(Date.now(), " | gRPC-Server | Received from rippled %v \n", in.Getmessage())
    log.Printf(Date.now(), " | gRPC-Server | Msg Validation key: %v \n", in.Getvalidator_key())
    // log.Printf("Received: %v", in.GetMessage())
    return &pb.Control{stream: true}, nil
}


func gRPCserver() {
    log.Println("___________________________________________")
    log.Println("Starting gRPC server")

    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCport))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
        log.Println("___________________________________________")
    }
    s := grpc.NewServer()
    pb.RegisterGossipMessageServer(s, &server{})
    log.Printf("server listening at %v", lis.Addr())
    log.Println("___________________________________________")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
        log.Println("___________________________________________")
    }
}
