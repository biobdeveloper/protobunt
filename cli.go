package protobunt

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	pb "protobunt/proto"
)

const (
	GET = "Get"
	SET = "Set"
	DELETE = "Delete"
)


type BuntClient struct {
	conn net.Conn
	dbName string
	protobuntVersion string
}


func CreateBuntClient(host string, port string) (pb.ProtoBuntClient, context.Context, context.CancelFunc) {
	var link = strings.Join([]string{host, port}, ":")
	conn, err := grpc.Dial(link, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewProtoBuntClient(conn)

	clientDeadline := time.Now().Add(time.Duration(1) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	r, err := client.VersionCheck(ctx, &pb.TestRequest{ClientVersion: "0.0.2"})
	if err != nil {
		log.Fatalf("could not exec: %v", err)
	}

	if r.GetServerVersion() != VERSION {
		log.Printf("Warning: different versions server and client: %s vs %s", r.GetServerVersion(), VERSION)
	}

	return client, ctx, cancel
}