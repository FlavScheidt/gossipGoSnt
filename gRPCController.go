package main

import (
    // "bufio"
    // "fmt"
    // "io/ioutil"
    // "log"
    // "os"
    // "strings"
    // "io"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "google.golang.org/grpc/examples/helloworld/helloworld"

    proto "github.com/golang/protobuf/proto"
    protoreflect "google.golang.org/protobuf/reflect/protoreflect"
    protoimpl "google.golang.org/protobuf/runtime/protoimpl"
    reflect "reflect"
    sync "sync"
)

const port = 20052

func client() {
    addr = "localhost:50051"

    log.Println("Starting gRPC client")
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)

    // Contact the server and print out its response.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.toRippled(ctx, &pb.Gossip{message: "Hello World", validator_key: "Hello"})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.GetMessage())
}