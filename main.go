package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nishi-yuki/pongdbot"
	"github.com/nishi-yuki/pongdbot/taskrunner"
)

var targetMacAddr string

func main() {
	loadEnv()
	targetMacAddr = os.Getenv("TARGET_MAC_ADDR")
	tr := makeTaskRunner()
	bot := pongdbot.GatherEnvVar(tr)
	bot.Start()
}

func loadEnv() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	err = load(
		".env",
		path.Join(home, ".raspizdbot.conf"),
		path.Join(home, "raspizdbot.conf"),
		path.Join(home, ".config", "raspizdbot.conf"),
	)

	if err != nil {
		fmt.Println(err)
	}
}

func load(filenames ...string) (err error) {
	for _, filename := range filenames {
		err = godotenv.Load(filename)
		if err == nil {
			return // return early on a spazout
		}
	}
	return
}

func makeTaskRunner() *taskrunner.TaskRunner {
	tr := taskrunner.New()
	tr.Add("ping", ping)
	tr.Add("wake", wake)
	return tr
}

func ping(m *discordgo.MessageCreate, _ []string) string {
	return "Pong! " + m.Author.String()
}

func wake(m *discordgo.MessageCreate, _ []string) string {
	err := exec.Command("wakeonlan", targetMacAddr).Run()
	if err != nil {
		return "Error: Failed to execute the wakeonlan command.\n" + err.Error()
	}
	return "wakeonlan command maybe ok."
}
