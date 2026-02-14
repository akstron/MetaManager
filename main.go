/*
Copyright © 2023 NAME HERE <dev.alok.singh123@gmail.com>
*/
package main

import (
	_ "embed"

	"github.com/heroku/self/MetaManager/cmd"
	"github.com/heroku/self/MetaManager/internal/services"
)

//go:embed credentials.json
var embeddedCredentials []byte

func main() {
	services.SetEmbeddedCredentials(embeddedCredentials)
	cmd.Execute()
}
