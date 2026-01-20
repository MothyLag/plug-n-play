package parser

import (
	"fmt"
	"mothylag/pnp/pkgs/entities"
	"mothylag/pnp/pkgs/input"
	"os"
	"regexp"
	"strings"
)

type GoParser struct {
	Tree *input.Tree
}

func CreateParser(t *input.Tree) *GoParser {
	return &GoParser{Tree: t}
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

func parse2Entity(content string) entities.EntitiesTree {
	tree := entities.CreateEntitiesTree()
	//pattern to get entity name and content into brackets block
	pattern := `(?:type)\s+(?P<Entity>\w+)\s+(?:struct)(?:\s*)(?:\n*)(?:{)(?:\n*)(?P<Content>[^}]*)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		entityName := match[1]
		fields := parseFields(match[2])
		entity := entities.NewEntity(entityName, fields)
		tree.AppendEntity(entity)
	}
	return tree
}

func parseFields(content string) []entities.Field {
	fields := []entities.Field{}
	//pattern to get field name and type from content
	pattern := `(?P<Key>\w+)\s+(?P<Value>\w+);`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		fields = append(fields, entities.Field{Name: key, Type: value})
	}
	return fields
}

func isEntityFile(f string) bool {
	return strings.HasSuffix(f, ".entity.go")
}
