package main

import (
	"fmt"
	"javilobo8-bot/config"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

var (
	client  *twitch.Client
	botUser = os.Getenv("TWITCH_CHAT_USER")
	botPass = os.Getenv("TWITCH_CHAT_PASS")
)

const (
	CmdChar    = '!'
	CmdConfig  = "c"
	CmdBanCopy = "bc"
	CmdOn      = "on"
	CmdOff     = "off"
)

type LastMessage struct {
	Message   string
	Timestamp int64
}

var lastMessages = make(map[string]LastMessage)

func saveLastMessage(channel, message string, now int64) {
	lastMessages[channel] = LastMessage{Message: message, Timestamp: now}
}

func getLastMessage(channel string) LastMessage {
	return lastMessages[channel]
}

func getChannelsToConnect() []string {
	channels := []string{config.Config.MainChannel}
	for _, item := range config.Config.Channels {
		channels = append(channels, item.Channel)
	}
	return channels
}

func main() {
	config.ReadConfig()
	client = twitch.NewClient(botUser, botPass)
	client.OnConnect(func() {
		log.Println("Twitch Client connected")
	})
	client.OnPrivateMessage(onPrivateMessage)
	client.Join(getChannelsToConnect()...)
	client.Connect()
}

func isSameMessage(a, b string) bool {
	return strings.Contains(a, b) || a == ("Chatting "+b)
}

func findChannelConf(channel string) *config.ChannelConfig {
	for _, item := range config.Config.Channels {
		if item.Channel == channel && item.Enabled {
			return &item
		}
	}
	return new(config.ChannelConfig)
}

func handleMainCommand(message twitch.PrivateMessage) {
	messages := strings.Split(message.Message, " ")
	if messages[0][0] != CmdChar {
		return
	}

	switch messages[0] {
	// !c bc channel on
	// !c bc channel off
	case (string(CmdChar) + CmdConfig):
		{
			if messages[1] == CmdBanCopy && len(messages) == 4 {
				// Set config
				// Save config
			}
		}
	}
}

func onPrivateMessage(message twitch.PrivateMessage) {
	now := time.Now().UnixMilli()

	if message.User.Name == config.Config.MainChannel {
		saveLastMessage(message.Channel, message.Message, now)
		handleMainCommand(message)
		return
	}

	channelConfig := findChannelConf(message.Channel)
	lastMessage := getLastMessage(message.Channel)

	// Has channelConfig and lastMessage
	if channelConfig != new(config.ChannelConfig) && &lastMessage != new(LastMessage) {
		// BanCopy action
		if channelConfig.BanCopyEnabled &&
			now < (channelConfig.BanCopyWindow*100)+lastMessage.Timestamp && // Check
			isSameMessage(message.Message, lastMessage.Message) {
			command := fmt.Sprintf("/timeout %s %d", message.User.Name, channelConfig.BanCopyTime)
			log.Println(message.Channel, command)
			client.Say(message.Channel, command)
		}
	}
}
