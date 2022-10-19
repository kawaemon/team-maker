package parser

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/kawaemon/team-maker/conf"
	"github.com/kawaemon/team-maker/g"
)

type ParseResult struct {
	TeamMembers g.Slice[string]
	TeamCount   int
}

var (
	regex = regexp.MustCompile(`^(?P<type>男子?|女子?|全員で|(?P<member_count>\d+)人で)?(?P<team_count>\d+)(?:チーム|グループ)作成$`)
)

func Parse(conf conf.Configuration, text string) (result ParseResult, ok bool) {
	match := regex.FindStringSubmatch(text)

	if match == nil {
		ok = false
		return
	}

	ok = true

	type_ := match[1]
	teamCountStr := match[3]

	var err error
	result.TeamCount, err = strconv.Atoi(teamCountStr)

	// 正規表現で文字列には正しい文字しか含まれていることが保証されているはず
	// ここでエラーが起きるのは完全に想定外なので落とす
	if err != nil {
		log.Fatalf("strconv.Atoi returned error: %s, input was %s", err, teamCountStr)
	}

	result.TeamMembers = g.NewSlice[string]()

	pushTeamMember := func(i int) {
		if name, ok := conf.NameMap[i]; ok {
			result.TeamMembers.Push(name)
		} else {
			result.TeamMembers.Push(strconv.Itoa(i))
		}
	}

	switch {
	// 男子だけ
	case strings.HasPrefix(type_, "男"):
		for i := 1; i <= conf.Total; i++ {
			if g.Contains(&conf.Women, i) {
				continue
			}

			pushTeamMember(i)
		}

	// 女子だけ
	case strings.HasPrefix(type_, "女"):
		for _, v := range conf.Women.Slice() {
			pushTeamMember(v)
		}

	// 全員で
	case type_ == "" || type_ == "全員で":
		for i := 1; i <= conf.Total; i++ {
			pushTeamMember(i)
		}

	// 指定された人数で
	case strings.HasSuffix(type_, "人で"):
		memberCountStr := match[2]
		memberCount, err := strconv.Atoi(memberCountStr)

		// 正規表現で文字列には正しい文字しか含まれていることが保証されているはず
		// ここでエラーが起きるのは完全に想定外なので落とす
		if err != nil {
			log.Fatalf("strconv.Atoi returned error: %s, input was %s", err, teamCountStr)
		}

		for i := 1; i <= memberCount; i++ {
			result.TeamMembers.Push(fmt.Sprintf("%d番", i))
		}
	}

	return
}
