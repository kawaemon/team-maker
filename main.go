package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kawaemon/group-maker/client"
	"github.com/kawaemon/group-maker/conf"
	"github.com/kawaemon/group-maker/g"
	"github.com/kawaemon/group-maker/parser"
	"github.com/kawaemon/group-maker/randomize"
)

func main() {
	godotenv.Load()
	config, err := conf.FromEnv()

	if err != nil {
		log.Fatalf("Failed to get config from environment varibale: %s", err)
	}

	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if channelAccessToken == "" {
		log.Fatalln("Set LINE_CHANNEL_ACCESS_TOKEN")
	}

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	if channelSecret == "" {
		log.Fatalln("Set LINE_CHANNEL_ACCESS_TOKEN")
	}

	lineClient := client.LineClient{
		ChannelAccessToken: channelAccessToken,
		ChannelSecret:      channelSecret,
		OnMessage: func(msg string) (result string) {
			parsed, ok := parser.Parse(config, msg)

			if !ok {
				return
			}

			if parsed.TeamCount <= 0 {
				log.Println("Group count was less than 1")
				return "作成するグループ数は1以上にしてください"
			}

			if parsed.TeamCount > parsed.TeamMembers.Len() {
				log.Println("Group count was bigger than team members")
				return "グループのメンバーよりチーム数の方が多いです"
			}

			log.Printf("Making group. msg: %s, count: %d\n", msg, parsed.TeamCount)

			randomized := randomize.Randomize(parsed)
			result = format(randomized)
			return
		},
	}

	err = lineClient.Start()

	if err != nil {
		log.Fatalf("Failed to start line client: %s", err)
	}
}

func format(groups g.Slice[g.Slice[string]]) (result string) {
	for index, group := range groups.Slice() {
		result += fmt.Sprintf("グループ%d\n", index+1)

		for _, member := range group.Slice() {
			result += fmt.Sprintf("%s\n", member)
		}

		result += "\n"
	}

	result = strings.TrimSpace(result)
	return
}
