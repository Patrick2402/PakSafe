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
// Logic based on comprehensive table of possible scenarios
func CheckForDependencyConfusion(dependencies []types.Dependency, privateRegistryURL string, privateScope string) map[string]types.DependencyStatus {
    results := make(map[string]types.DependencyStatus)
    publicRegistryURL := "https://registry.npmjs.org"
    
    if privateScope != "" {
        fmt.Printf("Using private scope: %s\n", privateScope)
    }
    
    for _, dep := range dependencies {
        status := types.DependencyStatus{
            PrivateVersion: dep.Version,
        }
        
        // Check if package exists in registries
        inPrivateRegistry := false
        inPublicRegistry := false
        var privateVersion, publicVersion string
        
        if privateRegistryURL != "" {
            inPrivateRegistry, privateVersion = checkPackageInRegistry(dep.Name, privateRegistryURL)
            status.PrivateVersion = privateVersion
        }
        
        inPublicRegistry, publicVersion = checkPackageInRegistry(dep.Name, publicRegistryURL)
        status.PublicVersion = publicVersion
        
        // Check if the package is scoped and if it belongs to the private scope
        isScoped := isOrganizationScoped(dep.Name)
        isOwnScope := belongsToPrivateScope(dep.Name, privateScope)
        
        if !isScoped && !isOwnScope {
            // Rows 11-16: Not scoped and not own scope
            if inPrivateRegistry && inPublicRegistry {
                status.Status = "suspicious"
                status.Reason = "Unscoped package exists in both registries"
                if compareVersions(publicVersion, privateVersion) > 0 {
                    status.Reason += " (public version is higher - possible version bombing)"
                }
            } else if inPrivateRegistry && !inPublicRegistry {
                status.Status = "vulnerable"
                status.Reason = "Unscoped package exists only in private registry (high risk for name squatting)"
                status.IsVulnerable = true
            } else if !inPrivateRegistry && inPublicRegistry {
                status.Status = "secure"
                status.Reason = "Public package exists only in public registry (normal scenario)"
            } else if !inPrivateRegistry && !inPublicRegistry {
                status.Status = "no package"
                status.Reason = "Package not found in any registry"
            }
        } else if isScoped && isOwnScope {
            // Rows 1-4: Scoped and own scope
            if inPrivateRegistry && inPublicRegistry {
                status.Status = "secure"
                status.Reason = "Own-scoped package exists in both registries (likely intentional)"
            } else if inPrivateRegistry && !inPublicRegistry {
                status.Status = "secure"
                status.Reason = "Own-scoped package exists only in private registry (secure scenario)"
            } else if !inPrivateRegistry && inPublicRegistry {
                status.Status = "vulnerable"
                status.Reason = "Own-scoped package exists only in public registry (high risk dependency confusion)"
                status.IsVulnerable = true
            } else if !inPrivateRegistry && !inPublicRegistry {
                status.Status = "not found"
                status.Reason = "Own-scoped package not found in any registry"
            }
        } else if isScoped && !isOwnScope {
            // Rows 5, 7, 9: Scoped but not own scope
            if inPrivateRegistry && inPublicRegistry {
                status.Status = "suspicious"
                status.Reason = "External-scoped package exists in both registries"
                if compareVersions(publicVersion, privateVersion) > 0 {
                    status.Reason += " (public version is higher - check if legitimate)"
                }
            } else if inPrivateRegistry && !inPublicRegistry {
                status.Status = "suspicious"
                status.Reason = "External-scoped package exists only in private registry (unusual scenario)"
            } else if !inPrivateRegistry && inPublicRegistry {
                status.Status = "secure"
                status.Reason = "External-scoped public package (normal scenario)"
            } else if !inPrivateRegistry && !inPublicRegistry {
                status.Status = "not possible"
                status.Reason = "External-scoped package not found anywhere"
            }
        } else if !isScoped && isOwnScope {
            // Rows 6, 8, 10: Not scoped but marked as own scope (should be impossible)
            status.Status = "not possible"
            status.Reason = "Package cannot be unscoped and have an own scope simultaneously"
        }
        
        results[dep.Name] = status
    }
    
    return results
}

// Check if a package belongs to the private scope
func belongsToPrivateScope(packageName, privateScope string) bool {
    if privateScope == "" { 
        return false
    }
    return strings.HasPrefix(packageName, privateScope)
}

// Check if a package name is organization scoped (starts with @)
func isOrganizationScoped(packageName string) bool {
    return strings.HasPrefix(packageName, "@")
}

// Check if a package exists in a registry and get its latest version
func checkPackageInRegistry(packageName, registryURL string) (bool, string) {
    if registryURL == "" {
        return false, ""
    }
    
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
    
    return ""
}

// Compare two semantic versions (simplified)
func compareVersions(version1, version2 string) int {
    if version1 == "" {
        return -1
    }
    if version2 == "" {
        return 1
    }
    
    if version1 > version2 {
        return 1
    } else if version1 < version2 {
        return -1
    }
    return 0
}