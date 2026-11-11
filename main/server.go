package main;

import(
	proto "Hand-In-4/grpc"
	"context"
	"fmt"
	"net"
	"log"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServer
	id 				int64
	isLeader 		bool
	leader			string
	peers 			[]string // Addresses of other nodes
	highestBid 		int64
	highestBidder 	int64
	allBidders 		map[string]bool
	endTime 		time.Time
}

func (s *AuctionServer) Bid(ctx context.Context, bid *proto.Bid) (*proto.Ack, error) {
	if !s.isLeader { // All non-leaders reject client requests. (For now): Uses a gRPC status error and includes the leader's address...
		return nil, status.Errorf(codes.FailedPrecondition, "not leader; leader=%s", s.leader)
	} 

	if (bid.Amount > s.highestBid) {
		s.highestBid = bid.GetAmount()
		s.highestBidder = bid.GetId()
		// Opdater backup servere før den sender til klienten
		return &proto.Ack{Outcome: "Succes"}, nil
	}
	// Update allBidders

	return &proto.Ack{Outcome: "Unsuccesfull"}, nil
}
func (s *AuctionServer) Result(ctx context.Context, empty *proto.Empty) (*proto.HighestBid, error){
	if !s.isLeader { // All non-leaders reject client requests. (For now): Uses a gRPC status error and includes the leader's address...
		return nil, status.Errorf(codes.FailedPrecondition, "not leader; leader=%s", s.leader)
	} 
	return &proto.HighestBid{Id: s.highestBidder,  Amount: s.highestBid, Status: "Auction over"}, nil
}
func makeServer(port string, NodeId int64) {
	
	// Listen on TCP port
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}
	server := &AuctionServer{
		id: 		NodeId,
		isLeader: 	NodeId == 2,	// Midlertidigt hardcoded til at være højeste Id ud af portene
		leader:		"localhost:5052",	// Midlertidigt hardcoded
		peers: 		[]string{"localhost:5050", "localhost:5051", "localhost:5052"}, // Midlertidigt hardcoded til at være på de 3 porte 
		highestBid: 0,
		allBidders: make(map[string]bool),	// Skal initialiseres for at kunne bruges senere
		endTime: 	time.Now().Add(100 * time.Second), // Skal dette være lamport clock i stedet? Ved ikke med synkronisering osv.
	}
	

	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServer(grpcServer, server) 

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}