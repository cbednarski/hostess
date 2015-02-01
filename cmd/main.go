package main

import (
	"github.com/cbednarski/hostess"
)

func getCommand() string {
	return os.Args[1]
}

func getArgs() []string {
	return os.Args[2:]
}

func main() {
	hostess.Hostess()
}
