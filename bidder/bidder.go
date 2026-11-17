package main

import (
	proto "Hand-In-4/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"flag"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clients[3] proto.AuctionClient
var id int64
var nodes [3]string // Addresses of other nodes
var leaderNodeId int // Last known leader
var name string;

func main(){
	nameFlag := flag.String("name", "bidder", "The name of the bidder")
	flag.Parse()
	name = *nameFlag


	for i := 0; i < 3; i++ {
		var basePort string = "localhost:505" + strconv.Itoa(i);
		conn, err := grpc.NewClient(basePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		clients[i] = proto.NewAuctionClient(conn)
		fmt.Println(i, ": ", clients[i]);
		}

		reader := bufio.NewReader(os.Stdin)

		for {
        text, _ := reader.ReadString('\n')

		if(strings.HasPrefix(text, "/bid")){
			bidInput := strings.Fields(text)
			amountText := bidInput[len(bidInput)-1]
			amount, err := strconv.ParseInt(amountText,10,64)
			
			if(err!=nil){
				log.Println(err)
			}

			bid(clients, amount)

		} else if (text == "/result\n"){
			result(clients)
		} else {
			fmt.Println("Not a known command, write /bid to make a bid and /result to see the highest bid")
		}

	

		
	}

}

func bid(clients[3] proto.AuctionClient, amount int64){
	bid := &proto.Bid{Name: name, Amount: amount}

	reply, err := clients[leaderNodeId].Bid(context.Background(), bid, grpc.WaitForReady(true))
	
	if(err != nil){
			log.Println(err)
			return;
		}
	fmt.Println(reply.Outcome)

}

func result(clients[3] proto.AuctionClient){
	
		reply, err := clients[leaderNodeId].Result(context.Background(), &proto.Empty{}, grpc.WaitForReady(true))
		if(err != nil){
			log.Println(err)
			return;
		}
		fmt.Println("Amount: " + strconv.Itoa(int(reply.Amount)), ", status: " + reply.Status)

}