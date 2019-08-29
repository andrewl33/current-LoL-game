# current LoL game

Discord bot that prints out a summoner's current game.

Input: `!currentLoL [username]`

Output:

```
Blue team
Summoner                 Champion        Rank      W/L    %
+--------------------+-----------+-----------+--------+----+
|Roy Pherae          |      Diana|    GOLD II|    8/11| 42%|
|CeciIia             |    Morgana|     GOLD I|   55/56| 49%|
|JigglyMuscles       |       Jhin|    GOLD IV|   18/19| 48%|
|Velvet Hoop         |      Garen|    GOLD IV|   26/21| 55%|
|Blargyn             |    Kindred|   GOLD III|   26/23| 53%|
+--------------------+-----------+-----------+--------+----+
Purple Team
Summoner                 Champion        Rank      W/L    %
+--------------------+-----------+-----------+--------+----+
|MidEvilKnight3k     |      Kayle|   GOLD III|   40/31| 56%|
|Twigman20           |         Vi|  SILVER II|   30/38| 44%|
|skittles1129        |     Ezreal| DIAMOND IV|   70/69| 50%|
|Riku Senpai         |      Janna|   BRONZE I|     4/7| 36%|
|SlipOnYogurt        |     Veigar|    GOLD II|   33/19| 63%|
+--------------------+-----------+-----------+--------+----+
```

## Requirements:

Needs global variables `RiotAPIKey` from [Riot Developer portal](https://developer.riotgames.com/) and `DiscordBotKey` from [Discord Developers](https://discordapp.com/developer).

## Flags:

`--updatechampions`: updates champion list, used when a new champion is released
