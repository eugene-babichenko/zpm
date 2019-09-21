package main

import "github.com/eugene-babichenko/zpm/commands"

// to be filled in during the build process
var version = "dev"

func main() {
	commands.Version = version
	commands.Execute()
}
