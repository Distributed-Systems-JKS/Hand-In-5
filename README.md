# Hand-In-5

## Run the program

Our program works for three nodes, no more or less.

The port should be specified in the command line when starting up a Node.

When starting the program open up 3 terminals.

Write in the first terminal:
```bash
go run ./main  -port=5050
```

Write in the second terminal:
```bash
go run ./main  -port=5051
```

Write in the third terminal:
```bash
go run ./main  -port=5052
```
## Own notes (DELETE WHEN SUBMITTING)

### Update proto file

#### 1:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```
#### 2:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/proto.proto
```
