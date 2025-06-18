package api

import (
  "context"
  "github.com/yourorg/kvstore/internal/dht"
  pb "github.com/yourorg/kvstore/proto"
)

// KVServer implements KVStore gRPC service.
type KVServer struct {
  pb.UnimplementedKVStoreServer
  table *dht.RoutingTable
  store *dht.Store
  self  string
}

// NewServer creates a new server.
func NewServer(selfAddr string) *KVServer {
  return &KVServer{
    table: dht.NewRoutingTable(selfAddr),
    store: dht.NewStore(),
    self:  selfAddr,
  }
}

func (s *KVServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
  peers := s.table.ClosestPeers(req.Key, 3)
  s.store.Put(req.Key, req.Value, peers)
  return &pb.PutResponse{Success: true}, nil
}

func (s *KVServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
  v, ok := s.store.Get(req.Key)
  return &pb.GetResponse{Value: v, Found: ok}, nil
}

func (s *KVServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
  s.table.AddPeer(req.Address)
  peers := s.table.ClosestPeers("", 10)
  return &pb.JoinResponse{Peers: peers}, nil
}

func (s *KVServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
  return &pb.PingResponse{Alive: true}, nil
}

func (s *KVServer) Heartbeat(stream pb.KVStore_HeartbeatServer) error {
  for {
    req, err := stream.Recv()
    if err != nil {
      return err
    }
    s.table.AddPeer(req.Sender)
    if err := stream.Send(&pb.PingResponse{Alive: true}); err != nil {
      return err
    }
  }
}
