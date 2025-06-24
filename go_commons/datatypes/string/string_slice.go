package string

type StringSlice []String

func (input StringSlice) Sanitize(includeChars ...String) StringSlice {
	var sanitizedString String
	sanitizedSlice := make([]String, 0, len(input))

	for _, str := range input {
		sanitizedString = str.Sanitize(includeChars...)
		sanitizedSlice = append(sanitizedSlice, sanitizedString)
	}

	return sanitizedSlice
}

func (input StringSlice) SanitizeAndPrune(includeChars ...String) StringSlice {
	var sanitizedString String
	sanitizedSlice := make([]String, 0, len(input))

	for _, str := range input {
		sanitizedString = str.Sanitize(includeChars...)
		if len(sanitizedString) > 0 {
			sanitizedSlice = append(sanitizedSlice, sanitizedString)
		}
	}

	return sanitizedSlice
}
