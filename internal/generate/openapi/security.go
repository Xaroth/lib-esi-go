package openapi

import "sort"

// RequiredOAuth2Scopes returns the OAuth scopes declared on an operation.
// Scopes from every security requirement are combined; any listed scope qualifies.
func RequiredOAuth2Scopes(security []SecurityRequirement) []string {
	if len(security) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var scopes []string
	for _, req := range security {
		for _, schemeScopes := range req {
			for _, scope := range schemeScopes {
				if seen[scope] {
					continue
				}
				seen[scope] = true
				scopes = append(scopes, scope)
			}
		}
	}
	sort.Strings(scopes)
	return scopes
}
