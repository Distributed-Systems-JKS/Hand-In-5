package main;

import(
	proto "Hand-In-4/grpc"
	"fmt"
	"strconv"
	"log"
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

	var id int64
	var nodes [3]string // Addresses of other nodes
	var leaderNodeId int // Last known leader
	var clients[3] proto.AuctionClient
	var database *os.File

func makeClient(port string, nodeId int64, databaseImport *os.File) {
	id = nodeId
	leaderNodeId = 2
	database = databaseImport

	//making the other client stubs, so we can communicate to the other servers
	for i := 0; i < 3; i++ {
		var basePort string = "localhost:505" + strconv.Itoa(i);
		conn, err := grpc.NewClient(basePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		clients[i] = proto.NewAuctionClient(conn)
		fmt.Println(i, ": ", clients[i]);
		}

	
}



func sendReplication(bidderName string, bidderAmount int64){

		for i := 0; i < 3; i++ {
			if(i != leaderNodeId){
				reply, err := clients[i].Replicate(context.Background(), &proto.Replicate{Name: bidderName, Amount: bidderAmount}, grpc.WaitForReady(true))
				if(err != nil){
					log.Println(err)
				
					continue;
				}

				fmt.Println(reply.Outcome);

		}
		
}
}

