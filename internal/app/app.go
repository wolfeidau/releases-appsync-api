package app

// These are set at build time - see makefile
var (
	Name      = "unknown"
	BuildDate = "unknown"
	Commit    = "unknown"
)

// Fields generate a fields map
func Fields() map[string]interface{} {
	return map[string]interface{}{
		"app_name":   Name,
		"build_date": BuildDate,
		"commit":     Commit,
	}
}

// MergeFields merge fields into one map, this happens in order they are provided
func MergeFields(maps ...map[string]interface{}) map[string]interface{} {

	fields := map[string]interface{}{}

	for _, m := range maps {
		for k, v := range m {
			fields[k] = v
		}
	}

	return fields
}