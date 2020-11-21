package parser

import (
	"fmt"
	"github.com/kawaemon/group-maker/conf"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ParseResult struct {
	TeamMembers []string
	TeamCount   int
}

var (
	regex = regexp.MustCompile(`^(?P<type>男子?|女子?|全員で|(?P<member_count>\d+)人で)?(?P<team_count>\d+)チーム作成$`)
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

	result.TeamMembers = []string{}
	switch {
	//男子だけ
	case strings.HasPrefix(type_, "男"):
		for i := 1; i <= conf.Total; i++ {
			// if confManNumbers.contains(i) { continue; }
			flag := false
			for _, man := range conf.WomanNumbers {
				if i == man {
					flag = true
				}
			}

			if flag {
				continue
			}

			name, ok := conf.NameMap[i]

			if ok {
				result.TeamMembers = append(result.TeamMembers, name)
			} else {
				result.TeamMembers = append(result.TeamMembers, strconv.Itoa(i))
			}
		}

	// 女子だけ
	case strings.HasPrefix(type_, "女"):
		for _, v := range conf.WomanNumbers {
			name, ok := conf.NameMap[v]

			if !ok {
				result.TeamMembers = append(result.TeamMembers, strconv.Itoa(v))
				continue
			}

			result.TeamMembers = append(result.TeamMembers, name)
		}

	// 全員で
	case type_ == "" || type_ == "全員で":
		for i := 1; i <= conf.Total; i++ {
			name, ok := conf.NameMap[i]

			if !ok {
				result.TeamMembers = append(result.TeamMembers, strconv.Itoa(i))
				continue
			}

			result.TeamMembers = append(result.TeamMembers, name)
		}

	// 指定された人数で
	case strings.HasSuffix(type_, "人で"):
		memberCountStr := match[2]
		memberCount, perr := strconv.Atoi(memberCountStr)

		if perr != nil {
			err = perr
			return
		}

		for i := 1; i <= memberCount; i++ {
			result.TeamMembers = append(result.TeamMembers, fmt.Sprintf("%d番", i))
		}
	}

	return
}
