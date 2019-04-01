package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mrletourneau/go-minebot-go/minebot"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

/*

This is a little bot for discord that can be used to spin up/down an EC2 instance that is presumeably running a
minecraft server (hence the name, minebot.)

It also can be used to keep a watch of player activity and spin down the server after a certain amount of time has
passed since the last player left. This is intended to save money.

Uses keapler's Batchcraft, bwmarrins's discordgo & aws's aws-sdk-go libraries (among others)

 */
const botPrefix string = "!minebotjr"
var commandsMap = make(map[string]func(str []string, s *discordgo.Session, m *discordgo.MessageCreate))
var Config minebot.Config

func main() {
	authToken := os.Getenv("DISCORD_AUTH_TOKEN")

	minebot.LoadConfig(&Config)

	if authToken == "" {
		fmt.Println( "No discord authentication token found. Exiting...")
		os.Exit(0)
	}

	fmt.Println("hello")
	discord, err := discordgo.New( "Bot " + authToken )

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	discord.AddHandler(messageCreate)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, _ := s.Channel(m.ChannelID)
	channelType := channel.Type

	dispatchCommand(m.Content, channelType, s, m)
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func dispatchCommand(message string, channelType discordgo.ChannelType, s *discordgo.Session, m *discordgo.MessageCreate) {
	arguments := strings.Split(message, " ")

	// If the command is in a public channel, make sure it is addressing minebot
	if channelType == discordgo.ChannelTypeGuildText && arguments[0] != botPrefix {
		return
	} else if channelType == discordgo.ChannelTypeGuildText && arguments[0] == botPrefix {
		arguments = arguments[1:]
	}

	// If the command is empty, ignore
	if len(arguments) < 1 {
		return
	}

	command, ok := commandsMap[arguments[0]]

	if !ok {
		help(arguments, s, m)
	} else {
		command(arguments, s, m)
	}

	// Determine if user is in private channel or not (in public channel, !minebot prefix needed)
	// Verify user has proper access to bring server up or down
	// Check status of server to see if up or down
	// Bring server up or down
	// Notify channel
}

func help(arguments []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	helpText := "Error, invalid command [%s]"
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(helpText, strings.Join(arguments, " ")))
}