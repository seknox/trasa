package utils

//ArrayContainsInt check if int array contains certain int
func ArrayContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
