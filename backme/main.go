package main

import "github.com/dimus/backme/backme/cmd"

var buildDate, buildVersion string

func main() {
	cmd.Execute(buildVersion, buildDate)
}
