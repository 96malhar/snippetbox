package main

type fieldEntry struct {
	valid     bool
	key       string
	errString string
}

func checkFields(entries ...fieldEntry) map[string]string {
	fieldErrors := map[string]string{}
	for _, e := range entries {
		if !e.valid {
			fieldErrors[e.key] = e.errString
		}
	}
	return fieldErrors
}
