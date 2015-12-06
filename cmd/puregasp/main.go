package main

import (
	"fmt"
	"os"

	"github.com/anacrolix/gasp"
	"github.com/anacrolix/tagflag"
)

var args struct {
	ProjectDir string `type:"pos"`
}

func main() {
	tagflag.Parse(&args)
	env := gasp.NewStandardEnv()
	err := env.RunProject(args.ProjectDir)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
