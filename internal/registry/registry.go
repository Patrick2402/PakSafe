package registry

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    
    "PakSafe/internal/types"
)

// Check for dependency confusion vulnerabilities
// Logic:
// - Package in private registry, not in public -> vulnerable
// - Package in private registry AND in public registry -> suspicious (version issues)
// - Package not in private registry, exists in public -> dependency confusion free
func CheckForDependencyConfusion(dependencies []types.Dependency, privateRegistryURL string, privateScope string) map[string]types.DependencyStatus {
    results := make(map[string]types.DependencyStatus)
    publicRegistryURL := "https://registry.npmjs.org"
    
    // Default scope for testing if none provided
    if privateScope == "" {
        // privateScope = "@your-scope"
    }
    
    fmt.Printf("Using private scope: %s\n", privateScope)
    
    // If no private registry specified, just check public registry
    if privateRegistryURL == "" {
        for _, dep := range dependencies {
            inPublicRegistry, publicVersion := checkPackageInRegistry(dep.Name, publicRegistryURL)
            
            status := types.DependencyStatus{
                IsVulnerable: false,
                Status: "available",
				Reason: "Package exists only in public registry",
                PrivateVersion: dep.Version,
                PublicVersion: publicVersion,
            }
            
            if !inPublicRegistry {
                status.Status = "vulnerable"
                status.Reason = "Package not found in public registry"
            }
            
            results[dep.Name] = status
        }
        return results
    }
    
    for _, dep := range dependencies {
        status := types.DependencyStatus{
            PrivateVersion: dep.Version,
        }
        
        // Check if package exists in private registry
        inPrivateRegistry, privateVersion := checkPackageInRegistry(dep.Name, privateRegistryURL)
        
        // Check if package exists in public registry
        inPublicRegistry, publicVersion := checkPackageInRegistry(dep.Name, publicRegistryURL)

        status.PrivateVersion = privateVersion
        status.PublicVersion = publicVersion
        
        // Apply the core logic based on the comments
        if inPrivateRegistry && !inPublicRegistry {
            // Case 1: In private registry but not in public -> vulnerable
            status.IsVulnerable = true
            status.Status = "vulnerable"
            status.Reason = "Package exists only in private registry (possible name squatting target)"
        } else if inPrivateRegistry && inPublicRegistry {
            // Case 2: In both registries -> suspicious
            status.Status = "suspicious"
            status.Reason = "Package exists in both registries"
            
            if compareVersions(publicVersion, privateVersion) > 0 {
                status.Reason += " (public version is higher - possible version bombing)"
            } else if compareVersions(publicVersion, privateVersion) < 0 {
                status.Reason += " (private version is higher)"
            } else {
                status.Reason += " (same versions - possible proxy caching)"
            }
        } else if !inPrivateRegistry && inPublicRegistry {
            // Case 3: Not in private registry but in public -> dependency confusion free
            status.Status = "available"
            status.Reason = "Package exists only in public registry (dependency confusion free)"
        } else {
            // Not in either registry
            status.Status = "not found"
            status.Reason = "Package not found in either registry"
        }
        
        // Special handling for packages with your organization's scope
        if belongsToPrivateScope(dep.Name, privateScope) {
            if !inPrivateRegistry && inPublicRegistry {
                // Your private-scoped package only exists in public registry
                status.IsVulnerable = true
                status.Status = "vulnerable"
                status.Reason = "Private-scoped package exists in public registry but not in private registry (high risk)"
            }
        }
        
        results[dep.Name] = status
    }
    
    return results
}

// Check if a package belongs to the private scope
func belongsToPrivateScope(packageName, privateScope string) bool {
    return strings.HasPrefix(packageName, privateScope)
}

// Check if a package name is organization scoped (starts with @)
func isOrganizationScoped(packageName string) bool {
    return strings.HasPrefix(packageName, "@")
}

// Check if a package exists in a registry and get its latest version
func checkPackageInRegistry(packageName, registryURL string) (bool, string) {
    url := fmt.Sprintf("%s/%s", strings.TrimRight(registryURL, "/"), packageName)
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error checking %s in registry %s: %v\n", packageName, registryURL, err)
        return false, ""
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 {
        version := extractLatestVersion(resp)
        return true, version
    }
    
    return false, ""
}

// Extract latest version from npm registry response
func extractLatestVersion(resp *http.Response) string {
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body: %v\n", err)
        return ""
    }
    
    var npmResp map[string]interface{}
    if err := json.Unmarshal(body, &npmResp); err != nil {
        fmt.Printf("Error parsing JSON response: %v\n", err)
        return ""
    }
    
    // Get the latest version from the "dist-tags" field
    if distTags, ok := npmResp["dist-tags"].(map[string]interface{}); ok {
        if latest, ok := distTags["latest"].(string); ok {
            return latest
        }
    }
    
    // If we can't find the latest version, return an empty string
    return ""
}

// Compare two semantic versions (simplified)
func compareVersions(version1, version2 string) int {
    // This is a simplified version comparison
    // For a complete solution, use a semantic versioning library
    
    // If either version is empty, consider it "lower"
    if version1 == "" {
        return -1
    }
    if version2 == "" {
        return 1
    }
    
    // For now, just compare strings (inadequate for real semver)
    if version1 > version2 {
        return 1
    } else if version1 < version2 {
        return -1
    }
    return 0
}