package randomize

import (
	"sort"

	"github.com/kawaemon/team-maker/g"
	"github.com/kawaemon/team-maker/parser"
)

func Randomize(data parser.ParseResult) g.Slice[g.Slice[string]] {
	data.TeamMembers.Shuffle()

	result := g.NewSliceWithCapacity[g.Slice[string]](data.TeamCount)

	for i := 0; i < data.TeamCount; i++ {
		result.Push(g.NewSlice[string]())
	}

	currentIndex := 0
	for _, v := range data.TeamMembers.Slice() {
		result.GetRef(currentIndex).Push(v)

		currentIndex += 1
		if currentIndex == data.TeamCount {
			currentIndex = 0
		}
	}

	for _, v := range result.Slice() {
		sort.Strings(v.Slice())
	}

	return result
}
