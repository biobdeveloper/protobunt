# ProtoBunt

![Go](https://github.com/biobdeveloper/protobunt/workflows/Go/badge.svg)

Server and client for [BuntDB] using [gRPC] over binary [Protocol Buffers] on top of simple TCP. 

[BuntDB] is awesome simple key/value storage written on Go. 
It can be used locally, but there is not out-of-box way to use remote database connection.
ProtoBunt claims to be the solution.

## Install
```sh
$ go get -u github.com/biobdeveloper/protobunt
```

## Usage

### Start BundDb Server
```go
package main

import "github.com/tidwall/buntdb"
import "github.com/biobdeveloper/protobunt"


func main() {
	// First, open some database...
	db, _ := buntdb.Open("dbname.db")
	
	// ...Then start the server
	host := "127.0.0.1"
	port := "8080"
	StartBuntServer(host, port, db)
	
}
```

### Setup BundDb Client
```go
package main

import "fmt"
import "github.com/biobdeveloper/protobunt"
import pb "protobunt/proto"


func main()  {
	
	// When the client is created, it makes the first test connection 
	// to negotiate the library version with the remote server
	cli, ctx, cancel := CreateBuntClient("127.0.0.1", "8080")
	defer cancel()

	// Create some record by pass key, value args
	req1 := pb.UpdateRequest{Key: "Alice", Value: "Bob", Action: SET}
	cli.Update(ctx, &req1)

	req2 := pb.ViewRequest{Key: "Alice", Action: GET}
	initValue, _ := cli.View(ctx, &req2)
	res2 := initValue.GetVal()

	fmt.Printf("Alice loves %s", res2)
	// Alice loves Bob
}
```

## Warning
ProtoBundDB is currenty alpha version!


[BuntDB]: https://github.com/tidwall/buntdb
[Protocol Buffers]: https://developers.google.com/protocol-buffers/
[gRPC]: https://github.com/grpc/grpc