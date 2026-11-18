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
	"io"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
	"os"
	"strconv"
	"encoding/csv"
	"google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type AuctionServer struct {
	proto.UnimplementedAuctionServer
	mu              sync.Mutex //mutex to fix concurrency issues
	id 				int64
	isLeader 		bool
	leader			string
	peers 			[]int64 // Ids of other nodes
	highestBid 		int64
	highestBidder 	string
	allBidders 		map[string]bool
	endTime 		time.Time
	database		*os.File
}

func (s *AuctionServer) Bid(ctx context.Context, bid *proto.Bid) (*proto.Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock() //Unlock when it returns: To take turns during the auction, so not all bids run at the same time

	fmt.Println("Bidder " + bid.GetName() + " send a bid");

	if bid.GetAmount() <= s.highestBid{
		return &proto.Ack{Outcome: "Unsuccesfull"}, nil //Bid too low
	}

	s.highestBid = bid.GetAmount()
	s.highestBidder = bid.GetName()

	w := csv.NewWriter(s.database);
	writeInDatabase(*w,bid.GetName(),bid.GetAmount())

	//Updating the backup serves before replying to client
	sendReplication(bid.GetName(), bid.GetAmount())
		
	return &proto.Ack{Outcome: "Success"}, nil
	
}

func (s *AuctionServer) Result(ctx context.Context, empty *proto.Empty) (*proto.HighestBid, error){
	s.mu.Lock() 
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

	database.Seek(0, io.SeekStart) //make sure we are at the top of the file
	r := csv.NewReader(database)
	data, err := r.ReadAll()

	if( err != nil){log.Println(err)}

	lastRecord := data[len(data)-1]
	s.highestBidder = lastRecord[0]
	s.highestBid, err = strconv.ParseInt(lastRecord[1], 10, 64)
	

	//Make sure if it fails it sends "fail", but don't know how to see if it fails?
	return &proto.ReplicateAck{Id: id, Amount: replicate.GetAmount(), Outcome: "Success"}, nil
}

func (s *AuctionServer) Election(ctx context.Context, election *proto.Empty)(*proto.Election, error){
	//just return because node only sends out to high-id nodes
	return &proto.Election{Id: id, IsLeader: s.isLeader}, nil;
}

func (s *AuctionServer) ServerNotResponding(ctx context.Context, message *proto.Empty)(*proto.Election, error){
	
	for i := 0; i < 3; i++ {
				ctxx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				if(int64(i) > id){
					fmt.Println(id, ": sending election to ", i)
					reply, err := clients[i].Election(ctxx, &proto.Empty{}, grpc.WaitForReady(true))
					if(err != nil){
						if status.Code(err) == codes.DeadlineExceeded {
						fmt.Println("Exceeded")
						// server did not reply in time, notify the other servers
        				s.isLeader = true;
						//Send messsage to all clients that a leader has been renewed (selected)
						return &proto.Election{Id: id, IsLeader: s.isLeader}, nil;
    					}
						continue;
						
					}
					log.Println("Got reply from: ", reply.Id)	
				
    	}
	}
	return &proto.Election{Id: id, IsLeader: s.isLeader}, nil;
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
		isLeader: 	NodeId == 2,	//Hardcoded: Leader is the highest node
		leader:		"localhost:5052",	//Hardcoded host
		highestBid: 0,
		endTime: 	time.Now().Add(10000 * time.Second), 
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

/*func ElectionTimer(s *AuctionServer ) *proto.Election{
	start := time.Now()
	for{
		if time.Since(start) >= 10*time.Second {
		//if server uses more than 10 seconds to answer tell the other servers
		s.isLeader = true;
			return &proto.Election{Id: id, IsLeader: s.isLeader};
		}	
	}
}*/

