package main;

import(
	"strconv"
	"log"
	"encoding/csv"
	"os"
	"fmt"
	"io"
)

func makeCSVfile(serverID int64){
	//Create a csv file

	filename := "Data" + strconv.Itoa(int(serverID)) + ".csv"
	f, err := os.Create(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

	//making a csv writer
   	w := csv.NewWriter(f)

	//Make a header in csv
    header := [][]string{
        {"Id", "amount", "version"},
    }

	//write the header in csv
	w.Write(header[0])

	//write other info in csv
	nodes := [][]string{
        {"1", "10", "1"},
		{"2", "20", "2"},
		{"3", "30", "3"},
		{"2", "40", "4"},
		{"1", "90", "5"},
    }

	for _, node := range nodes {
		w.Write(node)
	}

	//Ensures that the written data is in the csv file
	w.Flush()
	
	//Read from a csv file
	f.Seek(0, io.SeekStart) //make sure we are at the top of the file
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if (err != nil) {
		log.Printf("Error when reading csv file %v", err)
	}

	fmt.Println(len(records))

	lastRecord := records[len(records)-1]

	fmt.Println(lastRecord)

}

func main(){

	// Hardcoded to make 3 servers at the moment
	basePort := "505"
	amount := 3
	for i := 0; i < amount; i++ {
		
		/*server := &AuctionServer{
		id: 		int64(i),
		isLeader: 	i == amount,	// Midlertidigt hardcoded til at være højeste Id ud af portene
		leader:		basePort + strconv.Itoa(amount),	// Midlertidigt hardcoded
		peers: 		[]string{"localhost:5050", "localhost:5051", "localhost:5052"}, // Midlertidigt hardcoded til at være på de 3 porte 
		highestBid: 0,
		allBidders: make(map[string]bool),	// Skal initialiseres for at kunne bruges senere
		endTime: 	time.Now().Add(100 * time.Second), // Skal dette være lamport clock i stedet? Ved ikke med synkronisering osv.
	}*/

		go makeServer(basePort + strconv.Itoa(i), int64(i))
		go makeClient(int64(i))
	}

	select{} //so our main continues to run even if the goroutine finishes temporarily

}