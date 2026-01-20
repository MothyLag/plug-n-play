package parser

import (
	"fmt"
	"mothylag/pnp/pkgs/entities"
	"mothylag/pnp/pkgs/input"
	"os"
	"strings"
)

type GoParser struct {
	Tree     *input.Tree
	Entities []entities.Entity
}

func CreateParser(t *input.Tree) *GoParser {
	return &GoParser{Tree: t, Entities: []entities.Entity{}}
}

func (p *GoParser) FilterModelFiles() {
	filteredFiles := make([]string, 0, len(p.Tree.Files))
	for _, f := range p.Tree.Files {
		if isEntityFile(f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	p.Tree.Files = filteredFiles
}

func (p *GoParser) GetContent() {
	for _, file := range p.Tree.Files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}
		fmt.Printf("File %s\nContent:\n%s\n", file, string(data))
	}
}

func parseEntity(content string) {

}

func isEntityFile(f string) bool {
	return strings.HasSuffix(f, ".entity.go")
}
