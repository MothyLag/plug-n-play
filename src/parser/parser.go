package parser

import (
	"fmt"
	"mothylag/pnp/src/entities"
	"mothylag/pnp/src/input"
	"os"
	"strings"
)

type Parser struct{
	Tree *input.Tree
	Entities []entities.Entity
}

func CreateParser(t *input.Tree) *Parser{
	return &Parser{Tree: t,Entities: []entities.Entity{}}
}

func (p *Parser) FilterModelFiles(){
	filteredFiles := make([]string, 0, len(p.Tree.Files))
	for _,f := range p.Tree.Files{
		if isEntityFile(f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	p.Tree.Files = filteredFiles
}

func (p *Parser) GetContent() {
    for _, file := range p.Tree.Files {
        data, err := os.ReadFile(file)
        if err != nil {
            fmt.Printf("Error reading file %s: %v\n", file, err)
            continue
        }
        fmt.Printf("File %s\nContent:\n%s\n", file, string(data))
    }
}

func parseEntity(content string){
	
}

func isEntityFile(f string) bool{
	return strings.HasSuffix(f,".entity.go")
}
