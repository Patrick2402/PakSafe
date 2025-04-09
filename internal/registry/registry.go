package registry

import (
	"fmt"
	"net/http"
)

func CheckNpmRegistry(dependencies []string) map[string]bool {
	registryStatuses := make(map[string]bool)

	for _, dep := range dependencies {
		url := fmt.Sprintf("https://registry.npmjs.org/%s", dep)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error in request %s: %v\n", dep, err)
			registryStatuses[dep] = false
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			fmt.Printf("Dependency %s exists in public npm registry\n", dep)
			registryStatuses[dep] = true
		} else {
			fmt.Printf("Dependency %s does not exist in public npm registry\n", dep)
			registryStatuses[dep] = false
		}
	}

	return registryStatuses
}