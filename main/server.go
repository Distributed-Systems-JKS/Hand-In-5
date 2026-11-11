package main;

import(
	proto "Hand-In-4/grpc"
	"context"
	"fmt"
	"net"
	"log"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedAuctionServer

}

func (s *server) Bid(ctx context.Context, bid *proto.Bid) (*proto.Ack, error) {
	return &proto.Ack{Outcome: "Succes"}, nil
}
func (s *server) Result(ctx context.Context, empty *proto.Empty) (*proto.HighestBid, error){
	return &proto.HighestBid{Id: 1,  Amount: 100, Status: "Auction over"}, nil
}
func makeServer(port string) {
	
	// Listen on TCP port
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}


	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServer(grpcServer, &server{}) 

	

	//algorithmLogger.Printf("Server listening on %s", port)

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}