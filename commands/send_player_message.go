package commands

import (
	"database/sql"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func SendPlayerMessage(s *discordgo.Session, userID string, db *sql.DB, playerNick string) {
	password := generatePassword(db)
	message := "**Nick:** " + playerNick + "\n**Senha:** " + password

	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return
	}

	s.ChannelMessageSend(channel.ID, message)
}

func generatePassword(db *sql.DB) string {
	chars := "abcdefghjkmnpqrstuvwxyz23456789"

	for {
		password := ""
		for i := 0; i < 6; i++ {
			password += string(chars[rand.Intn(len(chars))])
		}

		var exists string
		err := db.QueryRow(`SELECT senha FROM players WHERE senha = ?`, password).Scan(&exists)

		if err == sql.ErrNoRows {
			return password
		}
	}
}
