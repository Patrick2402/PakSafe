package dependencies

import (
	"encoding/json"
	"fmt"
	"os"
)

// Create a struct to hold dependency information
type Dependency struct {
    Name    string
    Version string
}

func GetDependencies(path string) ([]Dependency, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("error reading the file: %v", err)
    }

    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil {
        return nil, fmt.Errorf("error in parsing: %v", err)
    }

    dependencies, ok := data["dependencies"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("no dependencies found in the file")
    }

    var result []Dependency
    for name, versionInterface := range dependencies {
        var version string
        if versionStr, ok := versionInterface.(string); ok {
            version = versionStr
        } else {
            version = "unknown"
        }
        
        result = append(result, Dependency{
            Name:    name,
            Version: version,
        })
    }

    return result, nil
}

func BuildJson(dependencies []Dependency, registryStatuses map[string]bool) map[string]interface{} {
    result := make(map[string]interface{})
    packages := make([]map[string]interface{}, len(dependencies))

    for i, dep := range dependencies {
        status := "available"
        if !registryStatuses[dep.Name] {
            status = "vulnerable"
        }

        packages[i] = map[string]interface{}{
            "package": dep.Name,
            "version": dep.Version, 
            "status":  status, 
        }
    }

    result["packages"] = packages
    return result
}

func SaveJsonToFile(data map[string]interface{}, path string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error during formatting JSON: %v", err)
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("cannot save to the file: %v", err)
	}

	return nil
}