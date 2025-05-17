# 🟥 o7-is: A Redis Clone in Go

A lightweight Redis clone implemented in Go, capable of handling basic Redis commands, RESP parsing, master-replica replication, and RDB file serving.

---

## 🚀 Features

- Minimal Redis-compatible server built in Go
- RESP protocol parser
- Basic command handling (GET, SET, etc.)
- Master-Slave replication with `PSYNC` support
- Embedded RDB snapshot serving
- Easy command registration and extension

---

## 📁 Project Structure

o7-is/
│
├── app/ # Entry point for the Redis server (main.go)
├── internal/
│ ├── commands/ # Command handlers (GET, SET, etc.)
│ ├── serverHelpers/ # Replication helpers and handshake logic
│ ├── tools/ # RESP parser, config utilities, state management
├── go.mod / go.sum
├── README.md


---

## 🛠️ How It Works

- The server listens on a specified port using the RESP protocol.
- It parses client requests and dispatches them to appropriate handlers.
- If run with `--replicaof`, it will:
  - Connect to the master,
  - Perform handshake,
  - Receive and process RDB file,
  - Sync data and keep connection alive for future updates.

---

## 🧪 Supported Commands

- `PING`
- `ECHO`
- `SET`
- `GET`
- `INFO`
- `REPLCONF`
- `PSYNC`

_(More commands can be added by editing the `internal/commands` package.)_

---

## ⚙️ Usage

### 🖥 Start as Master

```bash
go run app/main.go --port 6379

```


### 🖥 Start as Replica

The replica connects to the master, performs PSYNC, and syncs the RDB dump.

```bash
go run app/main.go --port 6380 --replicaof 127.0.0.1 6379

```

## 🧹 TODO / Improvements
- 🧠**Add support for more data types currently only supporting strings**
- ✅ **Basic master-replica setup**
- ⚙️ **Implement `SYNC` and real-time data propagation**
- ➕ **Add support for more Redis commands**  _(e.g., `DEL`, `EXPIRE`, `INCR`, etc.)_
- 🫡 **Implement eviction policies**  _(like LRU / TTL-based cleanup)_
- 💾 **Add persistence support**  _(support for RDB and AOF)_
- ⏱ **Handle data expiration and TTL**
- 🔒 **Add concurrency-safe data structures**  _(and improve locking mechanisms)_
- 📊 **Add benchmarking**  _(and compare performance with real Redis)_


## 🔨 Build

To build the Redis clone binary:

```bash
go build -o o7-redis ./app

./o7-redis --port 6379

```

## 🧪 Testing with redis-cli

You can use the official `redis-cli` to connect and test your Redis clone:

```bash
redis-cli -p 6379

```