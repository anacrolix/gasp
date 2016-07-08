package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anacrolix/tagflag"

	"github.com/anacrolix/gasp"
)

var args struct {
	tagflag.StartPos
	ProjectDir string `type:"pos"`
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	tagflag.Parse(&args)
	env := gasp.NewStandardEnv()
	err := env.RunProject(args.ProjectDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
