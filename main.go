/*
Copyright © 2023 NAME HERE <dev.alok.singh123@gmail.com>
*/
package main

import (
	_ "embed"

	"github.com/heroku/self/MetaManager/cmd"
)

//go:embed credentials.json
var embeddedCredentials []byte

func main() {
	cmd.SetEmbeddedCredentials(embeddedCredentials)
	cmd.Execute()
}
