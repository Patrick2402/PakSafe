package registry

import (
	"fmt"
	"net/http"
	"PakSafe/internal/dependencies"
)

func CheckNpmRegistry(dependencies []dependencies.Dependency) map[string]bool {
    registryStatuses := make(map[string]bool)

    for _, dep := range dependencies {
        url := fmt.Sprintf("https://registry.npmjs.org/%s", dep.Name)
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf("Error in request %s: %v\n", dep.Name, err)
            registryStatuses[dep.Name] = false
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode == 200 {
            registryStatuses[dep.Name] = true
        } else {
            registryStatuses[dep.Name] = false
        }
    }

    return registryStatuses
}