package env

func mergeEnvMaps(baseMap, newMap map[string]string) {
	for key, value := range newMap {
		baseMap[key] = value
	}
}
