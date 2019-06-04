package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	botToken string
	sheetID  string
	ranges   string
)

func init() {
	flag.StringVar(&botToken, "t", "", "Bot Token")
	flag.StringVar(&sheetID, "s", "", "sheet id")
	flag.StringVar(&ranges, "r", "", "ranges")
	flag.Parse()
}

//https://github.com/Chikachi/DiscordIntegration/wiki/How-to-get-a-token-and-channel-ID-for-Discord

var currentRunNum = 0
var runs = []Run{}

func main() {
	// spreadsheetId := "1Ej6FpGlhQs8wK6It60675tkwpxnpXOZGF33WmlMtJDE"
	// ranges := []string{"A2:C6", "E2:G6"}

	//"1Ej6FpGlhQs8wK6It60675tkwpxnpXOZGF33WmlMtJDE"
	//"A2:C6,E2:G6,I2:K6"

	//params these too

	refreshRuns()
	// fmt.Println("wat")

	// discord, err := discordgo.New("Bot " + botToken)
	// if err != nil {
	// 	fmt.Println("wdf couldnt make bot", err)
	// 	return
	// }

	// discord.AddHandler(onMessageReceived)

	// err = discord.Open()
	// if err != nil {
	// 	fmt.Println("failed to open discord", err)
	// 	return
	// }

	// fmt.Println("Listening")

	// //Wait on Signal
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc

	// discord.Close()
}

func refreshRuns() {
	runs = GetRuns(sheetID, ranges)
	// fmt.Println(runs)
	// fmt.Println(botToken)
	// fmt.Println(ranges)
	// fmt.Println(sheetID)
	currentRunNum = 0

}

func generateRunString(run Run) string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Group %c of weeklies is starting. @TODO", currentRunNum+'A'))

	for idx, person := range run.Runners {
		str.WriteString(fmt.Sprintf("%d: %s %s %s \n", idx, person.Name, BuffToStringMapping[person.Buff], "at meTODO"))
	}

	return str.String()
}

func onMessageReceived(session *discordgo.Session, msg *discordgo.MessageCreate) {

	// This isn't required in this specific example but it's a good practice.
	if msg.Author.ID == session.State.User.ID {
		return
	}

	var sentMessage = msg.Content
	if strings.HasPrefix(sentMessage, "!weekly") {
		if currentRunNum > len(runs) {
			session.ChannelMessageSend(msg.ChannelID, "There are no more runs... :shy:")
		} else {
			session.ChannelMessageSend(msg.ChannelID, generateRunString(runs[currentRunNum]))
		}

		currentRunNum += 1
	}
}
