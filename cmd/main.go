package main

import (
	// "github.com/cbednarski/hostess"
	"os"
)

func getCommand() string {
	return os.Args[1]
}

func getArgs() []string {
	return os.Args[2:]
}

func main() {

}
