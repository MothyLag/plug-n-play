package output

import (
	"fmt"
	"mothylag/pnp/src/input"
)

func ShowTree(tree *input.Tree){
	for _,f := range tree.Files{
		fmt.Println("📄 - ",f)	
	}
}
