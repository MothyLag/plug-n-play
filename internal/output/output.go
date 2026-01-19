package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"mothylag/pnp/internal/entities"
)

// ShowSummary writes a human-readable summary of the parsed entities to w.
func ShowSummary(w io.Writer, ents []entities.Entity) {
	if len(ents) == 0 {
		fmt.Fprintln(w, "no entities found")
		return
	}
	for _, e := range ents {
		fmt.Fprintf(w, "Entity: %s (fields: %d)\n", e.Name, len(e.Fields))
		for _, f := range e.Fields {
			fmt.Fprintf(w, "  - %s %s\n", f.Name, f.Type)
		}
		if len(e.DependsOn) > 0 {
			fmt.Fprintf(w, "  depends on: %s\n", strings.Join(e.DependsOn, ", "))
		}
		fmt.Fprintln(w)
	}
}

// WriteEntities writes the parsed entities to outDir. It creates the directory if
// necessary, writes an `entities.json` containing all entities and one markdown
// file per entity named `<entity>.md` with a simple representation.
//
// Returns an error if any filesystem operation fails.
func WriteEntities(outDir string, ents []entities.Entity) error {
	if len(ents) == 0 {
		// Nothing to write; ensure directory exists and return.
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("ensure outdir: %w", err)
		}
		return nil
	}

	// create output directory
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create outdir: %w", err)
	}

	// write consolidated JSON
	jsonPath := filepath.Join(outDir, "entities.json")
	f, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", jsonPath, err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ents); err != nil {
		f.Close()
		return fmt.Errorf("encode json: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close %s: %w", jsonPath, err)
	}

	// write one markdown file per entity
	for _, e := range ents {
		name := sanitizeFileName(e.Name)
		if name == "" {
			continue
		}
		mdPath := filepath.Join(outDir, name+".md")
		if err := writeEntityMarkdown(mdPath, &e); err != nil {
			return fmt.Errorf("write entity %s: %w", e.Name, err)
		}
	}

	return nil
}

func writeEntityMarkdown(path string, e *entities.Entity) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Simple markdown representation
	fmt.Fprintf(f, "# %s\n\n", e.Name)
	if len(e.Fields) == 0 {
		fmt.Fprintln(f, "_no fields_")
	} else {
		fmt.Fprintln(f, "## Fields")
		for _, fld := range e.Fields {
			fmt.Fprintf(f, "- `%s` — %s\n", fld.Name, fld.Type)
		}
		fmt.Fprintln(f)
	}

	if len(e.DependsOn) > 0 {
		fmt.Fprintln(f, "## Depends on")
		fmt.Fprintf(f, "%s\n", strings.Join(e.DependsOn, ", "))
	}

	return nil
}

// sanitizeFileName produces a filesystem-friendly name for an entity file.
// It lowercases the name and replaces spaces/backslashes with underscores.
// This is intentionally simple; callers should ensure names do not collide.
func sanitizeFileName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// replace path separators and spaces with underscore
	repl := strings.NewReplacer("/", "_", "\\", "_", " ", "_")
	s = repl.Replace(s)
	// remove any characters that are problematic for filenames by keeping a safe subset
	// a quick heuristic: keep letters, digits, dash and underscore and dot
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
		} else {
			// substitute anything else with underscore
			b.WriteByte('_')
		}
	}
	out := strings.ToLower(b.String())
	return out
}
