package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"bufio"
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
var nameToIDMap map[string]string;

var currentRunNum = -1
var runs = []Run{}
var weeklyRoleID = "419110547703726090"
//TODO temp hack, make it a map

func main() {
	if (sheetID == "") {
		fmt.Println("no sheet id specified");
		return;
	}

	if (botToken == "") {
		fmt.Println("no bot token specified");
		return;
	}
	if (ranges == "") {
		fmt.Println("no ranges speciied");
		return;
	}

	readNameToIDMap();

//todo
	//params these too
	// var run1 = []Person {
	// 	Person {
	// 		Name:"wat",
	// 		Dmg:500,
	// 		Buff:SB,
	// 	},
	// 	Person {
	// 		Name:"celvie",
	// 		Dmg:50430,
	// 		Buff:BLUE,
	// 	},
	// 	Person {
	// 		Name:"wa23t",
	// 		Dmg:504320,
	// 		Buff:NONE,
	// 	},
	// }
	// runs = append(runs, Run{
	// 	Runners:run1,
	// })
	// // refreshRuns()
	// fmt.Println(generateRunString(runs[currentRunNum]))
	

	runs = GetRuns(sheetID, ranges);
	setupDiscord();

	saveNameToIDMap();
}

var nameToIDMapFilePath = "names.txt"

func readNameToIDMap() {
	file, err := os.Open(nameToIDMapFilePath)
	if (err != nil) {
		fmt.Println("no previous name to id mapping file found");
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		nameToIDMap = make(map[string]string)
		for scanner.Scan() {
			var nametoid = strings.Split(scanner.Text()," ");
			nameToIDMap[nametoid[0]] = nametoid[1];
		}
	}
}

func saveNameToIDMap() {
	f, err := os.OpenFile(nameToIDMapFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("WDF WE FAILED TO SAVE O_O");
	} else {
		defer f.Close()
		f.WriteString(generateNameToIDMapString());
		f.Sync();
	}
}

func generateNameToIDMapString() string {
	var str strings.Builder
	for name, id := range nameToIDMap { 
		str.WriteString(fmt.Sprintf("%s %s\n", name, id));
	}

	return str.String()
}

func setupDiscord() {
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("wdf couldnt make bot", err)
		return
	}

	discord.AddHandler(onMessageReceived)

	err = discord.Open()
	if err != nil {
		fmt.Println("failed to open discord", err)
		return
	}

	fmt.Println("Listening")

	//Wait on Signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func refreshRuns() {
	runs = GetRuns(sheetID, ranges)
	// fmt.Println(runs)
	// fmt.Println(botToken)
	// fmt.Println(ranges)
	// fmt.Println(sheetID)
	currentRunNum = -1

}

func generateRunString(run Run) string {
	var str strings.Builder


	for idx, person := range run.Runners {
		var mention = ""
		if id, ok := nameToIDMap[strings.ToLower(person.Name)]; ok {
		    //do something here
		    mention = fmt.Sprintf("<@%s>", id)
		}
		str.WriteString(fmt.Sprintf("Player %d: %s %s %s %s \n", idx, person.Name, person.Class, BuffToStringMapping[person.Buff], mention))
	}

	return str.String()
}



func getWeeklyRoleID(session *discordgo.Session, guildID string) string {
	//weeklyRoleID = "retardedweeklies"
	if (weeklyRoleID != "") {
		return weeklyRoleID
	}
	var roles, err = session.GuildRoles(guildID)
	if (err != nil) {
		panic(err)
	}
	for _, role := range roles {
		if (role.Name == "retardedweeklies") {
			weeklyRoleID = role.ID
			return weeklyRoleID
		}
	}
	panic("could not find weekly role id in this discord")
	// well fk
	return weeklyRoleID
}

func onMessageReceived(session *discordgo.Session, msg *discordgo.MessageCreate) {

	// This isn't required in this specific example but it's a good practice.
	if msg.Author.ID == session.State.User.ID {
		return
	}

	var sentMessage = msg.Content
	sentMessage = strings.ToLower(sentMessage);
	if strings.HasPrefix(sentMessage, "!weekly ") {
		if (strings.Contains(sentMessage, "next")) {
			currentRunNum += 1
			if currentRunNum >= len(runs) {
				session.ChannelMessageSend(msg.ChannelID, "There are no more runs... :shy:")
			} else {
				// str.WriteString()
				var response = fmt.Sprintf("Group %c of weeklies is starting. <@&%s> \n", currentRunNum+'A', getWeeklyRoleID(session, msg.GuildID)) + generateRunString(runs[currentRunNum]);

				session.ChannelMessageSend(msg.ChannelID, response)
			}
		} else if (strings.Contains(sentMessage, "reset")) {
			refreshRuns()
			var response = fmt.Sprintf("Runs reset: %d runs found.", len(runs));
			session.ChannelMessageSend(msg.ChannelID, response)
		} else if (strings.Contains(sentMessage, "status")) {
			if currentRunNum < 0 {
				session.ChannelMessageSend(msg.ChannelID, "Hasn't started yet, or didnt update. :shy:")
			} else if currentRunNum >= len(runs) {
				session.ChannelMessageSend(msg.ChannelID, "There are no more runs... :shy:")
			} else {
				var response = fmt.Sprintf("Currently on group %c of weeklies. \n", currentRunNum+'A') + generateRunString(runs[currentRunNum]);

				session.ChannelMessageSend(msg.ChannelID, response)
			}
		} else if (strings.Contains(sentMessage, "register")) {
			var tokens = strings.Split(sentMessage, " ");
		
			if (len(tokens) == 3) {
				var name = strings.ToLower(tokens[2]);
				var ownerID = msg.Author.ID;
				nameToIDMap[name] = ownerID;
				session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("registered '%s' as <@%s>, you will be mentioned. :shy:", name, ownerID));
			} else {
				session.ChannelMessageSend(msg.ChannelID, "failed to register");
			}

		}  else if (strings.Contains(sentMessage, "registrants")) {
			session.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("I know these people\n%s", generateNameToIDMapString()))
		} else {
			session.ChannelMessageSend(msg.ChannelID, "huh?? :shy:")
		}

	}
}
