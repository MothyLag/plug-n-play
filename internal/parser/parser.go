package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"mothylag/pnp/internal/entities"
)

// FilterEntityFiles returns only files that look like entity/model files.
// Current heuristic: filename ends with ".entity.go".
func FilterEntityFiles(files []string) []string {
	out := make([]string, 0, len(files))
	for _, f := range files {
		if strings.HasSuffix(f, ".entity.go") {
			out = append(out, f)
		}
	}
	return out
}

// ParseEntities parses the provided Go source files and returns a slice of
// entities.Entity representing struct types found in those files.
// It returns an error if reading or parsing any file fails.
func ParseEntities(files []string) ([]entities.Entity, error) {
	var result []entities.Entity
	fset := token.NewFileSet()

	for _, fp := range files {
		data, err := os.ReadFile(fp)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", fp, err)
		}

		astFile, err := parser.ParseFile(fset, fp, data, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse file %s: %w", fp, err)
		}

		for _, decl := range astFile.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				ent := entities.Entity{
					Name:      typeSpec.Name.Name,
					Fields:    []entities.Field{},
					DependsOn: []string{},
				}
				depSet := map[string]struct{}{}

				if structType.Fields != nil {
					for _, f := range structType.Fields.List {
						typ := exprString(f.Type)
						// Field names may be nil for embedded fields.
						if len(f.Names) == 0 {
							// use a name derived from the type for embedded fields
							name := embeddedFieldName(f.Type)
							ent.Fields = append(ent.Fields, entities.Field{
								Name: name,
								Type: typ,
							})
							if shouldDependOn(typ) {
								depSet[baseTypeName(typ)] = struct{}{}
							}
						} else {
							for _, n := range f.Names {
								ent.Fields = append(ent.Fields, entities.Field{
									Name: n.Name,
									Type: typ,
								})
								if shouldDependOn(typ) {
									depSet[baseTypeName(typ)] = struct{}{}
								}
							}
						}
					}
				}

				// collect dependencies as a slice
				for d := range depSet {
					// avoid self-dependency
					if d != ent.Name && d != "" {
						ent.DependsOn = append(ent.DependsOn, d)
					}
				}

				result = append(result, ent)
			}
		}
	}

	return result, nil
}

// exprString converts common ast.Expr nodes into a readable type string.
// It's not a complete printer, but handles the usual cases (ident, selector,
// pointer, array, map, chan, func).
func exprString(e ast.Expr) string {
	switch tt := e.(type) {
	case *ast.Ident:
		return tt.Name
	case *ast.SelectorExpr:
		return exprString(tt.X) + "." + tt.Sel.Name
	case *ast.StarExpr:
		return "*" + exprString(tt.X)
	case *ast.ArrayType:
		if tt.Len == nil {
			return "[]" + exprString(tt.Elt)
		}
		// fixed length array
		return "[" + exprString(tt.Len) + "]" + exprString(tt.Elt)
	case *ast.MapType:
		return "map[" + exprString(tt.Key) + "]" + exprString(tt.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{...}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.ChanType:
		return "chan " + exprString(tt.Value)
	case *ast.BasicLit:
		return tt.Value
	case *ast.IndexExpr:
		return exprString(tt.X) + "[" + exprString(tt.Index) + "]"
	default:
		// fallback to a short description
		return fmt.Sprintf("%T", e)
	}
}

// embeddedFieldName returns a reasonable name for an embedded field based on the type.
func embeddedFieldName(e ast.Expr) string {
	switch tt := e.(type) {
	case *ast.Ident:
		return tt.Name
	case *ast.SelectorExpr:
		return tt.Sel.Name
	case *ast.StarExpr:
		return embeddedFieldName(tt.X)
	default:
		return exprString(e)
	}
}

// baseTypeName strips pointers/slices and package qualifiers and returns the base type name.
// Examples:
//
//	"*pkg.User" -> "User"
//	"[]*Order" -> "Order"
//	"map[string]Item" -> "Item" (returns the value type's base)
//
// This is a heuristic and not a full type resolver.
func baseTypeName(typ string) string {
	// remove pointer stars
	typ = strings.TrimLeft(typ, "*")
	// for slice notation
	if strings.HasPrefix(typ, "[]") {
		typ = typ[2:]
	}
	// for map -> get value type portion if possible
	if strings.HasPrefix(typ, "map[") {
		// try to split at first ']' and take remainder
		if idx := strings.Index(typ, "]"); idx != -1 && idx+1 < len(typ) {
			typ = typ[idx+1:]
		}
	}
	// if selector like pkg.Name, return the part after '.'
	if idx := strings.LastIndex(typ, "."); idx != -1 {
		typ = typ[idx+1:]
	}
	// trim any remaining pointer/slice markers
	typ = strings.TrimLeft(typ, "*[]")
	// remove any surrounding spaces
	typ = strings.TrimSpace(typ)
	return typ
}

// shouldDependOn determines whether the given type should be treated as a
// dependency on another entity. Basic builtin types are ignored.
func shouldDependOn(typ string) bool {
	base := baseTypeName(typ)
	if base == "" {
		return false
	}
	// common builtin types to ignore
	builtins := map[string]struct{}{
		"int": {}, "int8": {}, "int16": {}, "int32": {}, "int64": {},
		"uint": {}, "uint8": {}, "uint16": {}, "uint32": {}, "uint64": {},
		"uintptr": {}, "float32": {}, "float64": {},
		"string": {}, "bool": {}, "byte": {}, "rune": {},
		"error": {}, "complex64": {}, "complex128": {}, "interface{}": {},
	}
	_, ok := builtins[base]
	return !ok
}
