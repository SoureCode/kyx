package tools

func GetNames(mapping Mapping) []string {
	names := make([]string, 0, len(mapping))

	for name := range mapping {
		names = append(names, name)
	}

	return names
}
