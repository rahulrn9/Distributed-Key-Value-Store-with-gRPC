package dht

import (
  "context"
  pb "github.com/yourorg/kvstore/proto"
  "google.golang.org/grpc"
  "sync"
  "time"
)

const timeout = 5 * time.Second

// Store holds key-value pairs with a replication factor.
type Store struct {
  data map[string][]byte
  mux  sync.RWMutex
}

// NewStore initializes an in-memory store.
func NewStore() *Store {
  return &Store{data: make(map[string][]byte)}
}

// Put stores the value for key and replicates to peers.
func (s *Store) Put(key string, value []byte, peers []string) {
  s.mux.Lock()
  s.data[key] = value
  s.mux.Unlock()
  for _, peer := range peers {
    go func(p string) {
      conn, err := grpc.Dial(p, grpc.WithInsecure())
      if err != nil {
        return
      }
      client := pb.NewKVStoreClient(conn)
      ctx, cancel := context.WithTimeout(context.Background(), timeout)
      defer cancel()
      client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
    }(peer)
  }
}

// Get returns the value if present.
func (s *Store) Get(key string) ([]byte, bool) {
  s.mux.RLock()
  v, ok := s.data[key]
  s.mux.RUnlock()
  return v, ok
}
