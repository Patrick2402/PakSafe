package utils

import (
	"PakSafe/internal/dependencies"
	"fmt"

	"github.com/fatih/color"
)


func ListDependencies(dependencies []dependencies.Dependency) {
	color.Cyan("Found dependencies")
	for _, dependencies := range dependencies {
		fmt.Println("- ", dependencies)
	}
}