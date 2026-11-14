package main;

import(
	"strconv"
	"log"
	"encoding/csv"
	"os"
	"fmt"
	"io"

)

func makeCSVfile(){
	//Create a csv file
	f, err := os.Create("Data.csv")
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

	//makeCSVfile()
	
	/* portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()
	port := *portFlag
	
	idInt, err := strconv.ParseInt(port[len(port)-1:], 10, 64) //get the last digit of the port and set as id
	if(err != nil){
		log.Printf("Could not convert id %v", err)
	}*/


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
		database := makeDatabase(int64(i))
		go makeServer(basePort + strconv.Itoa(i), int64(i), database)
		go makeClient(int64(i), database)
	}

	select{} //so our main continues to run even if the goroutine finishes temporarily

}

func makeDatabase(id int64) *os.File{
	//Create a csv file
	f, err := os.Create("Data" + strconv.Itoa(int(id)) + ".csv");
    if err != nil {
        log.Fatal(err)
    }
    //defer f.Close()

	//making a csv writer
   	w := csv.NewWriter(f)

	//Make a header in csv
    header := [][]string{
        {"Id", "amount"},
    }

	//write the header in csv
	w.Write(header[0])

	//Ensures that the written data is in the csv file
	w.Flush()

	return f;
}

