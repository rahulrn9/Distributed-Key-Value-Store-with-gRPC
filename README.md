# Distributed Key-Value Store

A peer-to-peer, distributed key-value store built with Go, gRPC, and Kademlia-inspired routing. This system provides a fault-tolerant, scalable data storage solution that automatically replicates data across nodes and supports health checks via streaming heartbeats.

---

## Table of Contents

* [Features](#features)
* [Architecture](#architecture)
* [Project Structure](#project-structure)
* [Prerequisites](#prerequisites)
* [Installation](#installation)
* [Building the Application](#building-the-application)
* [Running Locally](#running-locally)
* [Docker Image](#docker-image)
* [Kubernetes Deployment](#kubernetes-deployment)
* [Testing](#testing)
* [Contributing](#contributing)
* [License](#license)

---

## Features

* **Peer Discovery**: Kademlia-style routing table based on exclusive OR distance metric.
* **Data Replication**: Asynchronous replication of key-value pairs to nearest peers.
* **gRPC API**: High-performance remote procedure calls for `Put`, `Get`, `Join`, `Ping`, and streaming `Heartbeat`.
* **Health Checks**: Bidirectional streaming heartbeats to monitor liveness.
* **Configurable Replication Factor**: Specify number of replicas for each key.
* **Scalable Architecture**: Stateless server processes with state managed in an in-memory store.

---

## Architecture

1. **Routing Table**: Maintains buckets of peer nodes indexed by distance from local node identifier.
2. **Data Store**: In-memory map of string keys to byte arrays, protected by read-write locks.
3. **gRPC Server**: Implements service methods defined in Protocol Buffers:

   * `Put`: Store data and replicate to closest peers.
   * `Get`: Retrieve data if available locally.
   * `Join`: Add a new peer and return current peer list.
   * `Ping`: Simple liveness check.
   * `Heartbeat`: Bidirectional streaming for continuous health monitoring.
4. **Client Join Logic**: On startup, each node optionally dials initial peers and exchanges join requests.
5. **Replication Logic**: `Put` requests trigger asynchronous gRPC calls to nearest neighbors.

---

## Project Structure

```text
kvstore/
├── cmd/
│   └── server/
│       └── main.go           # Application entrypoint and peer join logic
├── internal/
│   ├── dht/
│   │   ├── routing.go        # Routing table and distance calculations
│   │   └── store.go          # Key-value store with replication logic
│   └── api/
│       └── server.go         # gRPC service implementation and heartbeat handling
├── proto/
│   └── kv.proto              # Protocol Buffer definitions for gRPC services
├── deployments/
│   └── helm/
│       └── kvstore/          # Helm chart for Kubernetes StatefulSet deployment
├── Dockerfile                # Multi-stage build for container image
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── README.md                 # Project documentation
```

---

## Prerequisites

* Go **version 1.20** or later
* Docker (for containerization)
* Kubernetes cluster (for Helm deployment)
* Helm **version 3** or later

---

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourorg/kvstore.git
   cd kvstore
   ```
2. **Generate gRPC code**:

   ```bash
   protoc --go_out=internal/api --go-grpc_out=internal/api proto/kv.proto
   ```

---

## Building the Application

Compile the Go server:

```bash
go build -o kvstore cmd/server/main.go
```

---

## Running Locally

Start a single node with no initial peers:

```bash
./kvstore --addr=127.0.0.1:50051
```

Start a second node and join the first:

```bash
./kvstore --addr=127.0.0.1:50052 --peer=127.0.0.1:50051
```

Use `grpcurl` or a custom client to invoke `Put` and `Get` methods.

---

## Docker Image

Build and run the container:

```bash
docker build -t gcr.io/your-project/kvstore:latest .
docker run -p 50051:50051 gcr.io/your-project/kvstore:latest --addr=0.0.0.0:50051
```

---

## Kubernetes Deployment

Using the provided Helm chart:

```bash
helm repo add kvstore https://yourorg.github.io/kvstore
helm install kvstore deployments/helm/kvstore --set replicaCount=3,image.repository=gcr.io/your-project/kvstore,image.tag=latest
```

This creates a StatefulSet with headless Service to enable peer-to-peer DNS discovery.

---

## Testing

Run unit tests:

```bash
go test ./internal/dht
go test ./internal/api
```

---

## Contributing

Contributions are welcome. Please follow these steps:

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/my-feature`.
3. Make your changes and add tests.
4. Commit and push: `git push origin feature/my-feature`.
5. Open a pull request for review.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
