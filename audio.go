package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var buffer = make([][]byte, 0)

func PlayLocalAudio(s *discordgo.Session, m *discordgo.MessageCreate, filePath string, stopChan chan bool) {
	vs, err := findUserVoiceState(s, m.Author.ID)
	if err != nil {
		log.Fatalln("Could not find user channel, ", err)
	}

	// Connect to voice channel.
	// NOTE: Setting mute to false, deaf to true.
	err = playSound(filePath, s, m.GuildID, vs.ChannelID, stopChan)
	if err != nil {
		log.Fatal("Error playing sound, ", err)
		return
	}

	return
}

func playSound(fp string, s *discordgo.Session, guildID, channelID string, stop <-chan bool) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	file, err := os.Open(fp)
	if err != nil {
		log.Fatalln("Error opening dca file :", err)
		return err
	}
	defer file.Close()

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	//when stop is sent, kill ffmpeg
	go func() {
		<-stop
		vc.Disconnect()
	}()

	// Start speaking.
	vc.Speaking(true)

	//buff := make([]byte, 100)
	var opuslen int16

	// Send the buffer data.
	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		if opuslen < 0 {
			opuslen = 0
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			log.Fatal("Error reading streaming buffer, ", err)
		}

		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		vc.OpusSend <- InBuf
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}
