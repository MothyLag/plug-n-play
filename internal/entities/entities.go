package entities

// Entity represents a domain/entity model discovered in source files.
//
// Name is the type name (e.g. "User").
// Fields contains the named fields of the struct.
// DependsOn lists other entity type names this entity references (heuristic).
type Entity struct {
	// Name of the entity (Go type name).
	Name string `json:"name"`

	// Fields declared on the struct (in declaration order).
	Fields []Field `json:"fields"`

	// DependsOn contains other entity type names referenced by this entity.
	DependsOn []string `json:"depends_on"`
}

// Field represents a single struct field on an entity.
type Field struct {
	// Name of the field (as declared).
	Name string `json:"name"`

	// Type is a readable representation of the field's type (e.g. "string", "*User", "[]Order").
	Type string `json:"type"`
}
