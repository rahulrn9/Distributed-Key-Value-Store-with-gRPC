package dht_test

import (
  "testing"
  "github.com/yourorg/kvstore/internal/dht"
)

func TestBucketIndex(t *testing.T) {
  a, b := "a1b2", "c3d4"
  dist := dht.Distance(a, b)
  idx := dht.bucketIndex(dist)
  if idx < 0 {
    t.Errorf("invalid bucket index %d", idx)
  }
}