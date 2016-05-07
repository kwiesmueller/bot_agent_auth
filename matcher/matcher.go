package matcher

func Match(requiredParts []string, parts []string) bool {
	if len(requiredParts) != len(parts) {
		return false
	}
	for i, _ := range requiredParts {

		if requiredParts[i] != parts[i] {
			if len(requiredParts[i]) <= 2 ||
				requiredParts[i][0:1] != "[" ||
				requiredParts[i][len(requiredParts[i])-1:len(requiredParts[i])] != "]" {
				return false
			}
		}
	}
	return true
}
