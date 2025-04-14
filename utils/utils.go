package utils

import (
	"PakSafe/internal/types"
	"fmt"

	"github.com/fatih/color"
)

type Dependency struct {
	// Define the fields of the Dependency struct
}

func ListDependencies(dependencies []types.Dependency) {
	color.Cyan("Found dependencies")
	for _, dependency := range dependencies {
		fmt.Println("- ", dependency)
	}
}