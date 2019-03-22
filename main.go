package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PuuID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID string `json:"profileIconId"`
	RevisionDate  string `json:"revisionDate"`
	SummonerLevel string `json:"summonerLevel"`
}

// https://na1.api.riotgames.com
// /lol/summoner/v4/summoners/by-name/{summonerName}
// /lol/spectator/v4/active-games/by-summoner/{encryptedSummonerId}
// ?api_key=<key>

func main() {
	fmt.Println("starting currentLoLBot")
	dg, err := discordgo.New("Bot " + DiscordBotKey)

	if err != nil {
		fmt.Println("error on connecting,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if strings.Contains(m.Content, "!currentLoL") {
		arguments := strings.Fields(m.Content)
		if len(arguments) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !currentLoL [summoner name]")
			return
		}

		// get summoner ID
		resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + arguments[1] + "?api_key=" + RiotAPIKey)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Problem accessing Riot Games API.")
			fmt.Println("error getting name,", err)
			return
		}
		defer resp.Body.Close()

		summoner := Summoner{}

		json.NewDecoder(resp.Body).Decode(&summoner)

		fmt.Println(summoner.ID)
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
