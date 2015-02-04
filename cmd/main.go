package main

import (
	"fmt"
	"github.com/cbednarski/hostess"
	"os"
)

func getCommand() string {
	return os.Args[1]
}

func getArgs() []string {
	return os.Args[2:]
}

func main() {
	hostfile := hostess.NewHostfile(hostess.GetHostsPath())
	hostfile.Load()
	hostfile.Parse()
	fmt.Println(hostfile.Format())
}
