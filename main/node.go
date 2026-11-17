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



func main(){

	//makeCSVfile()
	
	portFlag := flag.String("port", "5050", "The port for this server to listen on")
	flag.Parse()
	port := *portFlag
	


	
	idInt, err := strconv.ParseInt(port[len(port)-1:], 10, 64) //get the last digit of the port and set as id
	if(err != nil){
		log.Printf("Could not convert id %v", err)
	}

	log.Println((idInt))

	
	database := makeDatabase(int64(idInt))
	go makeServer(*portFlag, idInt, database)

	go makeClient(*portFlag, idInt, database)

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
        {"Name", "amount"},
    }

	//write the header in csv
	w.Write(header[0])

	//Ensures that the written data is in the csv file
	w.Flush()

	return f;
}

func makeCSVfile(){ //Er her bare for at se hvordan man læser i en csv fil
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

	log.Println(err)

	fmt.Println(len(records))

	lastRecord := records[len(records)-1]

	fmt.Println(lastRecord)

	

}