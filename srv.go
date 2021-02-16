package protobunt

import (
	"context"
	"github.com/tidwall/buntdb"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "protobunt/proto"
)

type BuntServer struct {
	db *buntdb.DB
	protobuntVersion string
}

type server struct {
	pb.UnimplementedProtoBuntServer
	buntServer BuntServer
}


func (s *server) VersionCheck(ctx context.Context, in *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{ServerVersion: VERSION}, nil
}


func (s *server) View(ctx context.Context, in *pb.ViewRequest) (*pb.ViewResponse, error) {
	response := new(pb.ViewResponse)
	_ = s.buntServer.db.View(func(tx *buntdb.Tx) error {
		if in.Action == GET{
			val, err := tx.Get(in.GetKey())
			response.Val = val
			if err != nil {
				response.Error = err.Error()
			}
		}
		return nil
	})
	return response, nil
}


func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	response := new(pb.UpdateResponse)
	_  = s.buntServer.db.Update(func(tx *buntdb.Tx) error {
		if in.Action == SET {
			prev, repl, err := tx.Set(in.GetKey(), in.GetValue(), nil)
			response.PreviousValue = prev
			response.Replaced = repl
			if err != nil {
				response.Error = err.Error()
			}
		}
		if in.Action == DELETE {
			val, err := tx.Delete(in.GetKey())
			response.PreviousValue = val
			if err != nil {
				response.Error = err.Error()
			}
		}
		return nil
	})
	return response, nil
}


func StartBuntServer(host string, port string, db *buntdb.DB) {

	buntServer := BuntServer{}
	buntServer.protobuntVersion = VERSION

	buntServer.db = db
	defer buntServer.db.Close()

	link := host + ":" + port
	listener, _ := net.Listen("tcp", link)
	log.Println("Starting BuntDB Server...")

	defer listener.Close()

	s := grpc.NewServer()
	pb.RegisterProtoBuntServer(s, &server{buntServer: buntServer})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}