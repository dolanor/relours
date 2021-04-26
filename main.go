package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type config struct {
	token          string
	guildID        string
	voiceChannelID string
}

type bot struct {
	session *discordgo.Session
	voice   *discordgo.VoiceConnection
}

func loadConfig() config {
	var cfg config
	flag.StringVar(&cfg.token, "t", "", "Bot Token")
	flag.StringVar(&cfg.guildID, "g", "", "guild ID")
	flag.StringVar(&cfg.voiceChannelID, "c", "", "voice channel ID")
	flag.Parse()
	return cfg
}

func main() {

	cfg := loadConfig()
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}
	defer dg.Close()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	dgv, err := dg.ChannelVoiceJoin(cfg.guildID, cfg.voiceChannelID, false, true)
	if err != nil {
		if err, ok := dg.VoiceConnections[cfg.guildID]; ok {
			dgv = dg.VoiceConnections[cfg.guildID]
			if err != nil {
				log.Println("error connecting:", err)
				return
			}
		} else {
			log.Println("error connecting:", err)
			return
		}
	}

	defer dgv.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(b *bot) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		switch m.Author.Username {
		case "dolanor":
			_, err := s.ChannelMessageSend(m.ChannelID, "C'est pas faux.")
			if err != nil {
				log.Println("err:", err)
			}
			fallthrough
		default:
			_, err := s.ChannelMessageSend(m.ChannelID, "EN FAIT "+strings.ToUpper(m.Content))
			if err != nil {
				log.Println("err:", err)
			}

			dgvoice.PlayAudioFile(b.voice, "file.ogg", make(chan bool))

		}
	}
}
