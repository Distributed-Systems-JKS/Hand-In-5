package main;

import(
	proto "Hand-In-4/grpc"
	"fmt"
	"strconv"
	"log"
	"context"
	"math/rand/v2"
	"encoding/csv"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


	var id int64
	var nodes [3]string // Addresses of other nodes
	var leaderNodeId int // Last known leader
	var clients[3] proto.AuctionClient
	var database *os.File

func makeClient(nodeId int64, databaseImport *os.File) {
	id = nodeId
	leaderNodeId = 2
	database = databaseImport

	for i := 0; i < 3; i++ {
		nodes[i] = "localhost:505" + strconv.Itoa(i);
		conn, err := grpc.NewClient(nodes[i], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		clients[i] = proto.NewAuctionClient(conn)
	}

	sendBid()
	result()
	
}

func sendBid(){
	bid := &proto.Bid{Id: id, Amount: int64(rand.IntN(300))} // Sends a random bid from 0-300 (for now)
	fmt.Println("Client " + strconv.Itoa(int(id)) + " has sent: ")
	fmt.Println(bid.GetAmount())


	// Start out by trying known leaderNode first
	reply, err := clients[leaderNodeId].Bid(context.Background(), bid, grpc.WaitForReady(true))
	if err == nil {
		fmt.Println(reply.Outcome)
		return
	}

	log.Printf("not leader, leader is: %v", err)
	
	// Tries other nodes instead
	for i := range clients{
		if i == leaderNodeId {
			continue // We already tried this node
		}
		reply, err := clients[i].Bid(context.Background(), bid, grpc.WaitForReady(true))
		if(err != nil){
			log.Printf("not leader, leader is: %v", err)
			continue; // Tries next leader instead
		} 
		leaderNodeId = i // Set new leader
		fmt.Println(reply.Outcome)
	}
}

func result(){
	reply, err := clients[leaderNodeId].Result(context.Background(), &proto.Empty{}, grpc.WaitForReady(true))
	if err == nil {
		fmt.Println("Amount: " + strconv.Itoa(int(reply.Amount)), ", status: " + reply.Status)
		return
	}

	log.Printf("not leader, leader is: %v", err)

	for i := range clients{
		if i == leaderNodeId {
			continue // We already tried this node
		}

		reply, err := clients[i].Result(context.Background(), &proto.Empty{}, grpc.WaitForReady(true))
		if(err != nil){
			log.Println(err)
			continue;
		}
		leaderNodeId = i // Set new leader
		fmt.Println("Amount: " + strconv.Itoa(int(reply.Amount)), ", status: " + reply.Status)
	}
}

func replicateMessage(amount int64){
	reply, err := clients[leaderNodeId].Replicate(context.Background(), &proto.Replicate{Id: id, Amount: amount}, grpc.WaitForReady(true))
	if err == nil {
		w := csv.NewWriter(database)
		writeInDatabase(*w,reply.GetId(), reply.GetAmount())
		fmt.Println("Replicate: " + reply.Outcome + " in file: " + database.Name())
		return
	}
	
	log.Printf("not leader, leader is: %v", err)

	for i := range clients{
		if i == leaderNodeId {
			continue // We already tried this node
		}

		reply, err := clients[i].Replicate(context.Background(), &proto.Replicate{}, grpc.WaitForReady(true))
		if(err != nil){
			log.Println(err)
			continue;
		}
		
		leaderNodeId = i // Set new leader

		w := csv.NewWriter(database)
		writeInDatabase(*w,reply.GetId(), reply.GetAmount())

		fmt.Println("Replicate: " + reply.Outcome)
	}


	
}

