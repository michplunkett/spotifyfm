package utility


type HelperFunctions interface {
	ArrayHasNoEmptyStrings(envVars []string) bool
}

func ArrayHasNoEmptyStrings(envVars []string) bool {
	for _, value := range envVars {
		if value == EmptyString {
			return false
		}
	}

	return true
}
