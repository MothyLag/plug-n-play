package output

import (
	"fmt"
	"mothylag/pnp/pkgs/input"
)

func ShowTree(tree *input.Tree) {
	for _, f := range tree.Files {
		fmt.Println("📄 - ", f)
	}
}
