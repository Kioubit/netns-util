package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		printUsage(args[0])
	}

	switch os.Args[1] {
	case "enter":
		err := RunInNetns(args[2], args[3:], false)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "enter-netadmin":
		err := RunInNetns(args[2], args[3:], true)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "moveto":
		err := setLinkToNetns(os.Args[3], os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		printUsage(args[0])
	}

}

func printUsage(name string) {
	fmt.Println("Usage:")
	fmt.Println(name, "enter [netns name] [command]")
	fmt.Println(name, "enter-netadmin [netns name] [command]")
	fmt.Println(name, "moveto [netns name] [link]")
	os.Exit(2)
}
