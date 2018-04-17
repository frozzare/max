package cast

// getArg returns a arg value or the default value.
func getArg(def interface{}, i int, args ...interface{}) interface{} {
	if len(args) >= i+1 && args[i] != nil {
		return args[i]
	}

	return def
}

// getArgInt returns a arg value as int.
func getArgInt(def int, i int, args ...interface{}) int {
	return getArg(def, i, args...).(int)
}
