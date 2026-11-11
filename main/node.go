package main;

import(
	"flag"
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

	makeCSVfile()
	
	portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()
	port := *portFlag
	


	
	idInt, err := strconv.ParseInt(port[len(port)-1:], 10, 64) //get the last digit of the port and set as id
	if(err != nil){
		log.Printf("Could not convert id %v", err)
	}

	log.Println((idInt))

	go makeServer(*portFlag, int64(idInt))

	go makeClient(int64(idInt))

	select{} //so our main continues to run even if the goroutine finishes temporarily

}