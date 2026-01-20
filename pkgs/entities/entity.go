package entities

import "fmt"

type Entity struct {
	Name      string
	Fields    []Field
	DependsOn []string
}

type Field struct {
	Name string
	Type string
}

type EntitiesTree []Entity

func CreateEntitiesTree() EntitiesTree {
	return []Entity{}
}

func NewEntity(name string, fields []Field) Entity {
	return Entity{
		Name:      name,
		Fields:    fields,
		DependsOn: []string{},
	}
}

func (e *EntitiesTree) AppendEntity(newEntity Entity) {
	*e = append(*e, newEntity)
}

func (e *EntitiesTree) Show() {
	for _, entity := range *e {
		fmt.Printf("Entity: %s\n", entity.Name)
		for _, field := range entity.Fields {
			fmt.Printf("\tField: %s, Type: %s\n", field.Name, field.Type)
		}
	}
}
