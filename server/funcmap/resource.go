package funcmap

import (
	"fmt"
	"path"
)

// returns a resource function that generates a path to a resource file within a domain's resources directory.
//
// Creation: createResourceMappingFunction("catalog")
//
// Usage: {{ resource "main.js" }}
//
// Output: domain/catalog/resources/catalog.js
func createResourceMappingFunction(domain string) func(resource_path string) string {
	return func(resource_path string) string {
		return path.Join(fmt.Sprintf("domain/%s/resources", domain), resource_path)
	}
}
