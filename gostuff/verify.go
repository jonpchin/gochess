package gostuff

// checks if the time choices are valid
// gameType can be used to check correspondence times in minutes as well
func checkTime(choice int) bool {

	// 1440, 2880, 4320, 5760 are minutes for correspondence which are 1, 2, 3, 4 days
	var timeChoices = []int{1, 2, 3, 4, 5, 10, 15, 20, 30, 45, 1440, 2880, 4320, 5760}

	for _, v := range timeChoices {
		if choice == v {
			return true
		}
	}
	return false
}
