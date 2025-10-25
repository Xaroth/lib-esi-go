package requestgen

// IgnoredParameterComponents are resolved component parameter names excluded from Input.
var IgnoredParameterComponents = map[string]bool{
	"CompatibilityDate": true,
	"AcceptLanguage":    true,
	"IfModifiedSince":   true,
	"IfNoneMatch":       true,
	"Tenant":            true,
}

func isIgnoredParameter(componentName string) bool {
	return IgnoredParameterComponents[componentName]
}
