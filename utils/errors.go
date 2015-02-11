package utils

/**
 * Checks if passed error object is nil, panics otherwise
 */
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
