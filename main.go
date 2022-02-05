package main

import(
	"os"
	"fmt"
)

func main(){
	if len(os.Args) != 3 {
		fmt.Println("Error: Usage: ./ctf filename")
		os.Exit(0)
	}
	parse_file(os.Args[1])
}
