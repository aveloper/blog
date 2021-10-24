package main

import "github.com/aveloper/blog/internal/cmd"

//Version denotes the application version. It is populated on build using ld flags
var Version = "v0.0.0"

func main() {
	cmd.Execute(Version)
}
