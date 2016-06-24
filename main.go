package main

import (
	"fmt"
	"github.com/pshima/go-firewatch/firewatch"
	"log"
	"os"
)

const (
	version string = "0.0.1"
)

func main() {
	args := os.Args
	if len(args) != 3 || args[1] != "-prefix" {
		fmt.Println(help(args))
		os.Exit(1)
	}

	j := &firewatch.Job{Prefix: args[2]}
	err := j.Check()
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(255)
	}
	os.Exit(0)
}

func help(args []string) string {
	helpmsg := `%v %v Help

Required Arguments:
  -prefix <prefixname> (The prefix for alarm names in cloudwatch alarms)
`
	return fmt.Sprintf(helpmsg, args[0], version)
}
