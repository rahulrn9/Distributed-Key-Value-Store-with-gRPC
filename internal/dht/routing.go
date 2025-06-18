package dht

import (
  "math/big"
  "sync"
)

// RoutingTable maintains peer buckets.
type RoutingTable struct {
  localID   string
  buckets   map[int][]string
  bucketMux sync.Mutex
}

// NewRoutingTable initializes a table for the given node ID.
func NewRoutingTable(nodeID string) *RoutingTable {
  return &RoutingTable{localID: nodeID, buckets: make(map[int][]string)}
}

// Distance computes XOR-based distance between two node IDs (as hex strings)
func Distance(a, b string) *big.Int {
  ai := new(big.Int)
  bi := new(big.Int)
  ai.SetString(a, 16)
  bi.SetString(b, 16)
  return new(big.Int).Xor(ai, bi)
}

// bucketIndex returns the bucket for a given distance (first set bit index)
func bucketIndex(dist *big.Int) int {
  return dist.BitLen() - 1
}

// AddPeer inserts a peer into the appropriate bucket.
func (rt *RoutingTable) AddPeer(peerID string) {
  dist := Distance(rt.localID, peerID)
  idx := bucketIndex(dist)
  rt.bucketMux.Lock()
  defer rt.bucketMux.Unlock()
  bucket := rt.buckets[idx]
  if len(bucket) < 20 {
    rt.buckets[idx] = append(bucket, peerID)
  }
}

// ClosestPeers returns up to count peers closest to target ID.
func (rt *RoutingTable) ClosestPeers(target string, count int) []string {
  dist := Distance(rt.localID, target)
  idx := bucketIndex(dist)
  rt.bucketMux.Lock()
  defer rt.bucketMux.Unlock()
  var result []string
  for i := idx; i >= 0 && len(result) < count; i-- {
    result = append(result, rt.buckets[i]...)
  }
  for i := idx+1; i < len(rt.buckets) && len(result) < count; i++ {
    result = append(result, rt.buckets[i]...)
  }
  if len(result) > count {
    result = result[:count]
  }
  return result
}
