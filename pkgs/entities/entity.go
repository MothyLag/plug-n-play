package entities

type Entity struct {
	Name string
	Fields []Field
	DependsOn []string
}

type Field struct{
	Name string
	Type string
}


