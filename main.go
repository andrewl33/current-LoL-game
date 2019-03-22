package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Summoner info
type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PuuID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID string `json:"profileIconId"`
	RevisionDate  string `json:"revisionDate"`
	SummonerLevel string `json:"summonerLevel"`
}

// SummonerFail body from riot api
type SummonerFail struct {
	Status struct {
		Message    string `json:"message"`
		StatusCode int16  `json:"status_code"`
	}
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

		// get data
		encryptedAccountID, err := getUserID(arguments[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(encryptedAccountID)
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func getUserID(accountName string) (string, error) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + accountName + "?api_key=" + RiotAPIKey)

	if err != nil {
		return "", errors.New("could not connect to summoner api")
	}
	defer resp.Body.Close()

	// get summoner ID
	result := struct {
		Summoner
		SummonerFail
	}{}

	json.NewDecoder(resp.Body).Decode(&result)

	switch {
	case result.Status.Message != "":
		fmt.Println("summoner does not exist")
		return "", errors.New("summoner does not exist.")
	case result.Summoner.AccountID != "":
		fmt.Println(result.Summoner.AccountID)
		return result.Summoner.AccountID, nil

	}
	return "", errors.New("unreachable code in getUserId")
}
