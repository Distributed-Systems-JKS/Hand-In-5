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
	"time"


	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

var clients[3] proto.AuctionClient
var id int64
var nodes [3]string // Addresses of other nodes
var leaderNodeId int = 2 // Last known leader
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
	
		}

		fmt.Println("Welcome to the Auction. \n To make a bid write /bid <amount> \n To check the current state of the auction write /result")
		reader := bufio.NewReader(os.Stdin)

		for {
        text, _ := reader.ReadString('\n')

		if(strings.HasPrefix(text, "/bid")){
			bidInput := strings.Fields(text)
			amountText := bidInput[len(bidInput)-1]
			if(amountText == "/bid"){
				//If there is no amount written after /bid
				fmt.Println("Write the amount you want to bid after /bid")
				
			} else {
				amount, err := strconv.ParseInt(amountText,10,64)
			
			
				if(err!=nil){
					log.Println(err)
				}

				bid(amount)
			}
			

		} else if (text == "/result\n"){
			result()
		} else {
			fmt.Println("Not a known command, write /bid to make a bid and /result to see the highest bid")
		}

	

		
	}

}

func bid(amount int64){
	//making the bidder wait 10 seconds on reply, otherwise there will be a deadlineexceeded error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) 
	reply, err := clients[leaderNodeId].Bid(ctx, &proto.Bid{Name: name, Amount: amount}, grpc.WaitForReady(true))


	if(err != nil){
		if status.Code(err) == codes.DeadlineExceeded {
			// server did not reply in time, notify the other servers
        	notifyBackupServers(amount, true)
			return;
    	}
		log.Println(err)
		return;
	}
	
	// server replied in time, print outcome
	fmt.Println(reply.Outcome)
	
}


func result(){
	//making the bidder wait 10 seconds on reply, otherwise there will be a deadlineexceeded error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	reply, err := clients[leaderNodeId].Result(ctx, &proto.Empty{}, grpc.WaitForReady(true))
	
	if(err != nil){
		if status.Code(err) == codes.DeadlineExceeded {
			// server did not reply in time, notify the other servers
        	notifyBackupServers(0, false)
			return;
    	}
		log.Println(err)
		return;
	}



   	// server replied in time, print result
	fmt.Println("Amount: " + strconv.Itoa(int(reply.Amount)), ", status: " + reply.Status)
}
		




func notifyBackupServers(amount int64, isBid bool){
		
		// server did not reply in time, notify the other servers
		for i := 2; i >= 0; i-- {
			if(i != leaderNodeId){
				reply, err := clients[i].ServerNotResponding(context.Background(), &proto.Empty{}, grpc.WaitForReady(true))
			
				if(err != nil){
				log.Println(err)
				}
			
					
				if(reply.GetIsLeader() == true){
					leaderNodeId = int(reply.GetId())
					if(isBid){
						bid(amount)
					} else {
						result()
					}
					return;
				}
    
}
		
	}};