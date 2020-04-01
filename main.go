package main

import "fmt"

func main() {
	// server := &DataAnalysisServer{}
	// if err := http.ListenAndServe(":5000", server); err != nil {
	// 	log.Fatalf("could not listen on port 5000 %v", err)
	// }
	err := CheckCSV("clientData.csv", "ourData.csv")
	if err != nil {
		fmt.Println(err)
	}
}
