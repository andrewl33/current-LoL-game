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

// Fail body from riot api
type Fail struct {
	Status struct {
		Message    string `json:"message"`
		StatusCode int16  `json:"status_code"`
	}
}

// CurrentGame in progress
type CurrentGame struct {
	GameID            int64  `json:"gameId"`
	MapID             int    `json:"mapId"`
	GameMode          string `json:"gameMode"`
	GameType          string `json:"gameType"`
	GameQueueConfigID int    `json:"gameQueueConfigId"`
	Participants      []struct {
		TeamID                   int           `json:"teamId"`
		Spell1ID                 int           `json:"spell1Id"`
		Spell2ID                 int           `json:"spell2Id"`
		ChampionID               int           `json:"championId"`
		ProfileIconID            int           `json:"profileIconId"`
		SummonerName             string        `json:"summonerName"`
		Bot                      bool          `json:"bot"`
		SummonerID               string        `json:"summonerId"`
		GameCustomizationObjects []interface{} `json:"gameCustomizationObjects"`
		Perks                    struct {
			PerkIds      []int `json:"perkIds"`
			PerkStyle    int   `json:"perkStyle"`
			PerkSubStyle int   `json:"perkSubStyle"`
		} `json:"perks"`
	} `json:"participants"`
	Observers struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	PlatformID      string `json:"platformId"`
	BannedChampions []struct {
		ChampionID int `json:"championId"`
		TeamID     int `json:"teamId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bannedChampions"`
	GameStartTime int `json:"gameStartTime"`
	GameLength    int `json:"gameLength"`
}

// PlayerInfo
type PlayerInfo []struct {
	LeagueID     string `json:"leagueId"`
	LeagueName   string `json:"leagueName"`
	QueueType    string `json:"queueType"`
	Position     string `json:"position"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
}

// https://na1.api.riotgames.com
// /lol/summoner/v4/summoners/by-name/{summonerName}
// /lol/spectator/v4/active-games/by-summoner/{encryptedSummonerId}
// /lol/league/v4/positions/by-summoner/{encryptedSummonerId}
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

		// lookup user data
		encryptedSummonerID, err := getUserID(arguments[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		// lookup user game
		currentGame, err := getCurrentGame(encryptedSummonerID)
		if err != nil {
			fmt.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Not in game.")
		}

		fmt.Println(currentGame)

		// lookup teammates
		// allPlayerInfo, err := getAllPlayerInfo(currentGame)

		// send message to server
		// Team, champions playing, w/l, and rank
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
		Fail
	}{}

	json.NewDecoder(resp.Body).Decode(&result)

	switch {
	case result.Status.Message != "":
		return "", errors.New("summoner does not exist")
	case result.Summoner.AccountID != "":
		return result.Summoner.ID, nil

	}
	return "", errors.New("unreachable code in getUserId")
}

func getCurrentGame(encryptedSummonerID string) (CurrentGame, error) {
	result := struct {
		CurrentGame
		Fail
	}{}

	resp, err := http.Get("https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + encryptedSummonerID + "?api_key=" + RiotAPIKey)

	if err != nil {
		return result.CurrentGame, errors.New("could not connect to active-games API")
	}

	json.NewDecoder(resp.Body).Decode(&result)

	switch {
	case result.Status.Message != "":
		return result.CurrentGame, errors.New("summoner not in game")
	case result.CurrentGame.GameID > 0:
		return result.CurrentGame, nil

	}
	return result.CurrentGame, errors.New("unreachable code in getCurrentGame")
}

// func getAllPlayerInfo(currentGame CurrentGame) (PlayerInfo, error) {

// }
