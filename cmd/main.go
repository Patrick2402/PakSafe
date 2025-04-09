package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"PakSafe/internal/dependencies"
	"PakSafe/internal/registry"
)

func main() {
	defaultPath := "../tests/package.json"

	if len(os.Args) > 1 {
		defaultPath = os.Args[1]
	}

	deps, err := dependencies.GetDependencies(defaultPath)
	if err != nil {
		log.Fatalf("error during dowloading dependencies: %v", err)
	}

	if len(deps) == 0 {
		fmt.Println("Haven't found any dependencies")
		return
	}

	fmt.Println("Found dependencies:")
	for _, dep := range deps {
		fmt.Println(dep)
	}

	registryStatuses := registry.CheckNpmRegistry(deps)

	jsonResult := dependencies.BuildJson(deps, registryStatuses)

	outputPath := filepath.Join("output", "result.json")

	err = dependencies.SaveJsonToFile(jsonResult, outputPath)
	if err != nil {
		log.Fatalf("error during saving to JSON: %v", err)
	}

	fmt.Printf("Results saved to:  %s\n", outputPath)
}