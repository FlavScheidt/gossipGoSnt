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



func gRPCclient() {
    addr = "localhost:"+gRPCport

    log.Println("___________________________________________")
    log.Println("Starting gRPC client")
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
        log.Println("___________________________________________")
    }
    defer conn.Close()
    c := pb.NewGossipMessageClient(conn)

    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.toRippled(ctx, &pb.Gossip{message: "Hello World", validator_key: "Hello"})
    // log.Println(Date.now(), ' | gRPC-Client | Message from GSub node ID: ' + validator_key + ' sent to rippled server ');
    if err != nil {
        log.Fatalf("could not greet: %v", err)
        log.Println("___________________________________________")
    }
    log.Printf("Greeting: %s", r.Getmessage())
    log.Println("___________________________________________")
}