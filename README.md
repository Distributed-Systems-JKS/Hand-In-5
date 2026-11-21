# Hand-In-5

## Run the program

Our program works for three nodes, no more or less. 

The port should be specified in the command line when starting up a Node.

The name of the bidder should be specified in the command line when starting up a client. 
The different names must be unique.

When starting the program open up 5 terminals.

Write in the first terminal:
```bash
go run ./main  -port=5050
```

Write in the second terminal:
```bash
go run ./main  -port=5051
```

Write in the third terminal:
```bash
go run ./main  -port=5052
```

Write in the fourth terminal

```bash
go run bidder/bidder.go  -name=Bidder1
```
Write in the fifth terminal:

```bash
go run bidder/bidder.go  -name=Bidder2
```

