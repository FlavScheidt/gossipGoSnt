package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "github.com/FlavScheidt/gossipGoSnt/proto"



    // proto "github.com/golang/protobuf/proto"
    // protoreflect "google.golang.org/protobuf/reflect/protoreflect"
    // protoimpl "google.golang.org/protobuf/runtime/protoimpl"
    // reflect "reflect"
    // sync "sync"
)

func gRPCclientConnection() (pb.GossipMessageClient){
    addr := "localhost:20052"

    log.Println("Starting gRPC client")
    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
        log.Println("___________________________________________")
    }
    defer conn.Close()
    return pb.NewGossipMessageClient(conn)
}


func gRPCclientSend(message []byte, validatorKey string, hash string, clientStub pb.GossipMessageClient) {
   
    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := clientStub.ToRippled(ctx, &pb.Gossip{Message: message, Validator_Key: validatorKey, Hash: hash})
    log.Println("gRPC-Client | Message from GSub node ID: " + string(validatorKey) + " sent to rippled server");
    if err != nil {
        log.Fatalf("could not send message: %v", err)
        log.Println("___________________________________________")
    }
    log.Printf("Message sent to rippled | %t", r.GetStream())
}