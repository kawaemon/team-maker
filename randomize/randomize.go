package randomize

import (
	"github.com/kawaemon/group-maker/parser"
	"math/rand"
	"sort"
)

func Randomize(data parser.ParseResult) [][]string {
	rand.Shuffle(len(data.TeamMembers), func(i, j int) {
		data.TeamMembers[i], data.TeamMembers[j] = data.TeamMembers[j], data.TeamMembers[i]
	})

	result := make([][]string, 0, data.TeamCount)

	for i := 0; i < data.TeamCount; i++ {
		result = append(result, make([]string, 0))
	}

	currentIndex := 0
	for _, v := range data.TeamMembers {
		result[currentIndex] = append(result[currentIndex], v)

		currentIndex += 1
		if currentIndex == data.TeamCount {
			currentIndex = 0
		}
	}

	for _, v := range result {
		sort.Strings(v)
	}

	return result
}
