package entities

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
