package main

import (
	"flag"
	"fmt"
	"mothylag/pnp/src/input"
	"mothylag/pnp/src/output"
	"mothylag/pnp/src/parser"
	"path/filepath"
)

func main(){
	flag.Parse()
	args := flag.Args()
    if len(args) < 1 {
		fmt.Println("usage: pnp <route>")
    }
	inPath := "./"
	outPath := "./plugs"
    if len(args) >= 1 {
        inPath = args[0]
    }
    if len(args) >= 2 {
        outPath = args[1]
    }

 	absInPath,err := filepath.Abs(inPath)
	if err != nil {
		panic(err)
	}

	absOutPath,err := filepath.Abs(outPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("processing models at: ",absInPath)
	fmt.Println("output dir at: ",absOutPath)
	tree := input.CreateTree(absInPath)
	parsedTree := parser.CreateParser(tree)
	parsedTree.FilterModelFiles()
	output.ShowTree(parsedTree.Tree)
	parsedTree.GetEntities()
}
