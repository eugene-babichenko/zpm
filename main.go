package main

import "github.com/eugene-babichenko/zpm/cmd"

// to be filled in during the build process
var version = "dev"

func main() {
	cmd.Version = version
	cmd.Execute()
}
