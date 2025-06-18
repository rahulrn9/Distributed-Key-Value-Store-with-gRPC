package api_test

import (
  "context"
  "net"
  "testing"

  "github.com/yourorg/kvstore/internal/api"
  pb "github.com/yourorg/kvstore/proto"
  "google.golang.org/grpc"
)

func TestHeartbeat(t *testing.T) {
  lis, _ := net.Listen("tcp", ":0")
  srv := grpc.NewServer()
  s := api.NewServer(lis.Addr().String())
  pb.RegisterKVStoreServer(srv, s)
  go srv.Serve(lis)

  conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
  if err != nil {
    t.Fatal(err)
  }
  client := pb.NewKVStoreClient(conn)
  stream, err := client.Heartbeat(context.Background())
  if err != nil {
    t.Fatal(err)
  }

  go stream.Send(&pb.PingRequest{Sender: "node1"})
  resp, err := stream.Recv()
  if err != nil || !resp.Alive {
    t.Errorf("expected alive response")
  }
}