package dependencies

import (
	"encoding/json"
	"fmt"
	"os"
)


func GetDependencies(path string) ([]string, error) {

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
		return nil, fmt.Errorf("no depencencies found in the file")
	}

	var result []string
	for dependency := range dependencies {
		result = append(result, dependency)
	}

	return result, nil
}

func BuildJson(dependencies []string, registryStatuses map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	packages := make([]map[string]interface{}, len(dependencies))

	for i, dep := range dependencies {
		status := "available"
		if !registryStatuses[dep] {
			status = "vulnerable"
		}

		packages[i] = map[string]interface{}{
			"package": dep,
			"version": "1.0.0", 
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