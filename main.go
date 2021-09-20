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
	htgotts "github.com/hegedustibor/htgo-tts"
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
	//dg.LogLevel = discordgo.LogDebug

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	dgv, err := dg.ChannelVoiceJoin(cfg.guildID, cfg.voiceChannelID, false, false)
	if err != nil {
		log.Println("error connecting:", err)
		return
	}
	defer dgv.Close()

	b := bot{
		session: dg,
		voice:   dgv,
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(b.messageCreate)
	dgv.AddHandler(b.voiceHandler)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (b *bot) voiceHandler(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	log.Println("new voice event:", vs.UserID, vs)

	switch vs.UserID {
	case "yyyyyy":
		// FIXME can't send TTS message in a voice channel
		// _, err := b.session.ChannelMessageSendTTS(vc.ChannelID, "uuuuuuwwwww")
		// if err != nil {
		// 	log.Println("err:", err)
		// }
	case "xxxxxx":
		//_, err := b.session.ChannelMessageSendTTS(vc.ChannelID, "la funk")
		//if err != nil {
		//	log.Println("err:", err)
		//}
	case "zzzzzz":
		stop := make(chan bool)
		dgvoice.PlayAudioFile(b.voice, "./duc.mp3", stop)
		<-stop
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//return func(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Author.Username {
	case "dolanor":
		log.Println("dolanor")
		_, err := s.ChannelMessageSend(m.ChannelID, "C'est pas faux.")
		if err != nil {
			log.Println("err:", err)
		}
		content := strings.ReplaceAll(m.Content, "je", "il")
		content = strings.ReplaceAll(content, "Je", "Il")
		content = strings.ReplaceAll(content, "JE", "IL")
		_, err = s.ChannelMessageSendTTS(m.ChannelID, content)
		if err != nil {
			log.Println("err:", err)
		}

		// play audio file
		stop := make(chan bool)
		dgvoice.PlayAudioFile(b.voice, "./duc.mp3", stop)
		<-stop

	default:
		log.Println("les autres")
		_, err := s.ChannelMessageSend(m.ChannelID, "EN FAIT "+strings.ToUpper(m.Content))
		if err != nil {
			log.Println("err:", err)
		}

	}
}

func findVoiceChannelID(guild *discordgo.Guild, m *discordgo.MessageCreate) string {
	log.Println("voice states start")
	for _, voiceState := range guild.VoiceStates {
		log.Println("voice states:", voiceState)
		if voiceState.UserID == m.Author.ID {
			return voiceState.ChannelID
		}
	}
	return ""
}

func Speak() {
	speech := htgotts.Speech{Folder: "audio", Language: "fr"}
	// TODO the: need to get the audio bytes
	speech.Speak("oulala")

}
