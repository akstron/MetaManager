package main

import (
	"log"

	"github.com/heroku/self/MetaManager/cmd" // import your root command

	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, ".")
	if err != nil {
		log.Fatal(err)
	}
}
