package main

import (
  "context"
  "flag"
  "log"
  "net"
  "strings"

  "github.com/yourorg/kvstore/internal/api"
  pb "github.com/yourorg/kvstore/proto"
  "google.golang.org/grpc"
)

func main() {
  addr := flag.String("addr", "127.0.0.1:50051", "listen address")
  peers := flag.String("peer", "", "comma-separated initial peers")
  flag.Parse()

  srv := api.NewServer(*addr)
  if *peers != "" {
    for _, p := range strings.Split(*peers, ",") {
      conn, err := grpc.Dial(p, grpc.WithInsecure())
      if err != nil {
        log.Printf("join dial %s error: %v", p, err)
        continue
      }
      client := pb.NewKVStoreClient(conn)
      client.Join(context.Background(), &pb.JoinRequest{Address: *addr})
    }
  }

  lis, err := net.Listen("tcp", *addr)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  grpcServer := grpc.NewServer()
  pb.RegisterKVStoreServer(grpcServer, srv)
  log.Printf("gRPC server listening on %s", *addr)
  grpcServer.Serve(lis)
}
