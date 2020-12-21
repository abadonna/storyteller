package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"regexp"
	"storyteller/engine"
	"storyteller/game"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type gameData struct {
	game      *game.SampleGame
	timestamp time.Time
}

var (
	token string
	games map[string]*gameData
)

const imageFolder = "game/"

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
	games = make(map[string]*gameData)
	time.AfterFunc(time.Minute*15, checkExpired)
}

func mainDiscord() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
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

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if message.Author.ID == session.State.User.ID {
		return
	}

	//Ingore all messages in public channels, DM only
	ch, err := session.Channel(message.ChannelID)
	if err != nil || ch.Type != 1 {
		return
	}

	input := strings.ToLower(message.Content)

	if input == "start" || input == "restart" {
		context := &gameData{
			game:      game.Sample(),
			timestamp: time.Now()}

		games[message.Author.ID] = context

		messages := parse(context.game.Intro())
		/*

			ms := &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Image: &discordgo.MessageEmbedImage{
						URL: "https://i.imgur.com/ar18tWi.jpg"}},
				Content: msg}

			//URL: "attachment://" + fileName,

			/*,
				Files: []*discordgo.File{
					&discordgo.File{
						Name:   fileName,
						Reader: f,
					},
				},
			}*/

		for _, msg := range messages {
			session.ChannelMessageSendComplex(message.ChannelID, msg)
		}
		return
	}

	context, ok := games[message.Author.ID]
	if !ok {
		session.ChannelMessageSend(message.ChannelID, "No active game found.\nType \"start\" to play.")
		return
	}

	context.timestamp = time.Now()

	messages := parse(engine.Process(context.game, message.Content))

	for _, msg := range messages {
		session.ChannelMessageSendComplex(message.ChannelID, msg)
	}
}

func checkExpired() {
	t := time.Now()
	for key, context := range games {
		diff := t.Sub(context.timestamp).Hours()
		if diff > 1 {
			delete(games, key)
		}
	}
	time.AfterFunc(time.Minute*15, checkExpired)
}

func parse(src string) []*discordgo.MessageSend {
	result := []*discordgo.MessageSend{}
	re := regexp.MustCompile(`\[\[([^\[\]]*)\]\]`)
	if !re.MatchString(src) {
		result = append(result, &discordgo.MessageSend{Content: src})
		return result
	}

	submatchall := re.FindAllString(src, -1)
	ms := &discordgo.MessageSend{}
	hasImage := false
	text := src

	for _, element := range submatchall {
		command := strings.Trim(element, "[")
		command = strings.Trim(command, "]")
		command = strings.Replace(command, " ", "", -1)
		if strings.HasPrefix(command, "img=") {
			img := strings.Replace(command, "img=", "", 1)
			idx := strings.Index(text, element)
			var (
				f        *os.File
				fileName string
			)

			if !strings.HasPrefix(strings.ToLower(img), "http") {
				current, _ := os.Executable()
				f, _ = os.Open(path.Join(path.Dir(current), imageFolder, img))

				fileName = img
				img = "attachment://" + img
			}

			if hasImage {
				hasImage = false
				ms.Content = text[:idx]
				result = append(result, ms)
				ms = &discordgo.MessageSend{
					Embed: &discordgo.MessageEmbed{
						Image: &discordgo.MessageEmbedImage{
							URL: img}}}
				text = strings.Replace(text[idx:], element, "", 1)
			} else {
				hasImage = true
				ms.Embed = &discordgo.MessageEmbed{
					Image: &discordgo.MessageEmbedImage{
						URL: img}}
				text = strings.Replace(text, element, "", 1)
			}

			if f != nil {
				ms.Files = []*discordgo.File{
					&discordgo.File{
						Name:   fileName,
						Reader: f,
					}}
				defer f.Close()
			}

		} else if command == "break" {
			idx := strings.Index(text, element)
			hasImage = false
			ms.Content = text[:idx]
			result = append(result, ms)
			ms = &discordgo.MessageSend{}
			text = strings.Replace(text[idx:], element, "", 1)
		}
	}

	ms.Content = text
	result = append(result, ms)

	return result
}
