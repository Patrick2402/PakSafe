package types

type Dependency struct {
    Name    string
    Version string
}

type DependencyStatus struct {
    IsVulnerable   bool
    Status         string
    Reason         string
    PrivateVersion string
    PublicVersion  string
}