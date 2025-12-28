package utils

func CalculateAverage(ratings []int) int {
	if len(ratings) == 0 {
		return 0
	}

	sum := 0

	for _, rating := range ratings {
		sum += rating
	}

	return sum / len(ratings)
}
