package dependencies

import (
    "encoding/json"
    "fmt"
    "os"
    
    "PakSafe/internal/types"
)

func GetDependencies(path string) ([]types.Dependency, error) {
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

    var result []types.Dependency
    for name, versionInterface := range dependencies {
        var version string
        if versionStr, ok := versionInterface.(string); ok {
            version = versionStr
        } else {
            version = "unknown"
        }
        
        result = append(result, types.Dependency{
            Name:    name,
            Version: version,
        })
    }

    return result, nil
}

func BuildJson(dependencies []types.Dependency, registryStatuses map[string]types.DependencyStatus) map[string]interface{} {
    result := make(map[string]interface{})
    packages := make([]map[string]interface{}, len(dependencies))

    for i, dep := range dependencies {
        status, exists := registryStatuses[dep.Name]
        statusText := "unknown"
        reason := ""
        
        if exists {
            statusText = status.Status
            reason = status.Reason
        }

        packages[i] = map[string]interface{}{
            "package": dep.Name,
            "version": dep.Version,
            "status": statusText,
            "reason": reason,
            "privateVersion": status.PrivateVersion,
            "publicVersion": status.PublicVersion,
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