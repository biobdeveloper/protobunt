# ProtoBunt

![Go](https://github.com/biobdeveloper/protobunt/workflows/Go/badge.svg)

Server and client for BuntDB using simple text protocol over TCP. 

[BuntDB] is awesome simple key/value storage written on Go. 
It can be used locally, but there is not out-of-box way to use remote database connection.
ProtoBunt trying to be a solution.

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
	// Firstly, open some database...
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


func main()  {
	
	// When the client is created, it makes the first test connection 
	// to negotiate the library version with the remote server
	cli := CreateBuntClient("127.0.0.1", "8080")

	// Create some record by pass key, value args
	cli.Update(SET, "Alice", "Bob")

	// And look at this
	get := cli.View(GET, "Alice")
	fmt.Printf("Alice loves %s", get)
         // Alice loves Bob
}
```

## Warning
ProtoBundDB is currenty alpha version!

## Why not 3rd-side modern protocol like gRPC?
Now I'm choosing a protocol to use in the future. I will be glad to receive advices and issues.

## Protocol Specification
Currently using simple text protocol
* `\t `symbol as the request parts separator
* `\n `symbol as the end of request string

Client sending requests with text data 
* `View\tGet\t{Alice}\n`
* `Update\tSet\t{Alice:Bob}\n`
* ...

[BuntDB]: https://github.com/tidwall/buntdb