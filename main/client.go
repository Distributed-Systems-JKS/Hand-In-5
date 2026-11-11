package main;

import(
	proto "Hand-In-4/grpc"
	"fmt"
	"strconv"
	"log"
	"context"
	"math/rand/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func makeClient(port string) {

	address := fmt.Sprintf("localhost:%s", port)
	
	//conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//client := proto.NewClientServerClient(conn)

	//making the other client stubs, so we can communicate to the other servers
	var clients[2] proto.AuctionClient
	clientPosition := 0 //making sure that we only add two clients to the list in the for-loop
	for i := 0; i < 3; i++ {
		var basePort string = "localhost:505" + strconv.Itoa(i);
		if(basePort != address){
			conn, err := grpc.NewClient(basePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
			//log.Printf("Connecting to gRPC server on port 505" + strconv.Itoa(i))
			if err != nil {
				log.Fatalf("Did not connect: %v", err)
			}
			clients[clientPosition] = proto.NewAuctionClient(conn)
			clientPosition++
		}
		
	}


	bid(clients)
	result(clients)
	
}

func bid(clients[2] proto.AuctionClient){
	
	for i := range clients{
		
		reply, err := clients[i].Bid(context.Background(), &proto.Bid{Id: 2, Amount: int64(rand.IntN(300))}, grpc.WaitForReady(true))
		fmt.Println(reply.Outcome)

		if(err != nil){
			log.Println(err)
			return;
		}
	}
	
	
}

func result(clients[2] proto.AuctionClient){
	
	for i := range clients{
		reply, err := clients[i].Result(context.Background(), &proto.Empty{}, grpc.WaitForReady(true))
		fmt.Println("Amount: " + strconv.Itoa(int(reply.Amount)), ", status: " + reply.Status)

		if(err != nil){
			log.Println(err)
			return;
		}
	}
}