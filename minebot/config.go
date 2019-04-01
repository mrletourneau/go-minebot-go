package minebot

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

var minebotConfigFileDir string

type ServerConfig struct{
	Name string
	IdleShutdownTime int `toml:"idle_shutdown_time"`
	ServerId string `toml:"server_id"`
}

type Config struct{
	CommandPrefix string `toml:"command_prefix"`
	DiscordAuthKey string `toml:"discord_auth_key"`
	Admins []string
	ChannelAnnounce string `toml:"channel_announce"`
	Server []ServerConfig
}

func init() {
	minebotConfigFileDir = os.Getenv("MINEBOT_CONFIG_FILE_DIR")

	if minebotConfigFileDir == "" {
		minebotConfigFileDir = os.Getenv("HOME")
	}

	if minebotConfigFileDir != "" {
		minebotConfigFileDir += "/.minebot"
	}

	var fileMode os.FileMode
	fileMode = 0755
	os.Mkdir(minebotConfigFileDir, fileMode)
	fmt.Println(minebotConfigFileDir)
}

func createConfig() {
	serverConfig := ServerConfig {
		"Example Server #1",
		30,
		"XXXXXXXXX",
	}
	config := Config {
		"!minebot",
		"XXXXXXXXX",
		[]string {"ExampleUser1"},
		"General",
		[]ServerConfig {serverConfig},
	}
	buf := bytes.NewBufferString("")
	encoder := toml.NewEncoder(buf)
	encoder.Encode(config)
	err := ioutil.WriteFile(minebotConfigFileDir + "/config.toml", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error creating config file", err)
		os.Exit(0)
	}
}

func LoadConfig(c *Config) {
	config, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	if _, err := toml.Decode(string(config), c); err != nil {
		log.Fatal(err)
	}
}
