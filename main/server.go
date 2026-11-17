package main;

import(
	proto "Hand-In-4/grpc"
	"context"
	"fmt"
	"sync"
	"net"
	"log"
	"time"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
	"os"
	"strconv"
	"encoding/csv"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServer
	mu              sync.Mutex //mutex to fix concurrency issues
	id 				int64
	isLeader 		bool
	leader			string
	peers 			[]string // Addresses of other nodes
	highestBid 		int64
	highestBidder 	string
	allBidders 		map[string]bool
	endTime 		time.Time
	database		*os.File
}

func (s *AuctionServer) Bid(ctx context.Context, bid *proto.Bid) (*proto.Ack, error) {
	s.mu.Lock() //Lock now
	defer s.mu.Unlock() //Unlock when it returns: To take turns during the auction, so not all bids run at the same time

	fmt.Println("Bidder " + bid.GetName() + " send a bid");

	if bid.GetAmount() <= 0 || time.Now().After(s.endTime) { 
		return &proto.Ack{Outcome: "Auction is closed"}, nil //or/and invalid amount
	}
	if bid.GetAmount() <= s.highestBid{
		return &proto.Ack{Outcome: "Unsuccesfull"}, nil //Bid too low
	}

	s.highestBid = bid.GetAmount()
	s.highestBidder = bid.GetName()

	//s.allBidders[strconv.Itoa(int(bid.GetId()))] = true

	w := csv.NewWriter(s.database);
	writeInDatabase(*w,bid.GetName(),bid.GetAmount())

	// Opdater backup servere før den sender til klienten
	sendReplication(bid.GetName(), bid.GetAmount())
		
	return &proto.Ack{Outcome: "Success"}, nil
	
	// Update allBidders
}

func (s *AuctionServer) Result(ctx context.Context, empty *proto.Empty) (*proto.HighestBid, error){
	s.mu.Lock() //lock now
	defer s.mu.Unlock() //Unlock when it returns: To take turns during the auction, so not all bids run at the same time
	

	//Replaced Status: "Auction over" with this instead, so it will only print when the auction is actually over
	status := "The auction is in process"
	if time.Now().After(s.endTime){
		status = "The auction is now over" 
	}

	return &proto.HighestBid{Name: s.highestBidder,  Amount: s.highestBid, Status: status}, nil
}

func (s *AuctionServer) Replicate(ctx context.Context, replicate *proto.Replicate)(*proto.ReplicateAck, error){
	//write in the database
	w := csv.NewWriter(database)
	writeInDatabase(*w,replicate.GetName(), replicate.GetAmount())

	//Update highest bidder here??
	//Make sure if it fails it sends "fail", but don't know how to see if it fails?
	return &proto.ReplicateAck{Id: id, Amount: replicate.GetAmount(), Outcome: "Success"}, nil
}


func writeInDatabase(w csv.Writer, name string, amount int64){
	
	bid := [][]string{{name, strconv.Itoa(int(amount))}}

	w.Write(bid[0])
	
	//Ensures that the written data is in the csv file
	w.Flush()
}
func makeServer(port string, NodeId int64, database *os.File) {
	// Listen on TCP port
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}

	server := &AuctionServer{
		id: 		NodeId,
		isLeader: 	NodeId == 2,	// Midlertidigt hardcoded til at være højeste Id ud af portene, kan benytte heartbeat gRPC Ping?
		leader:		"localhost:5052",	// Midlertidigt hardcoded
		highestBid: 0,
		endTime: 	time.Now().Add(100 * time.Second), // Skal dette være lamport clock i stedet? Ved ikke med synkronisering osv.
		database:   database,
	}
	// Creating a new gRPC server instance
	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServer(grpcServer, server) 


	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}