/*
Copyright © 2024 Miguel Chateloin
*/
package main

import "github.com/Cenergistic/cmdeagle/cmd"

// version is set during build by goreleaser
var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
