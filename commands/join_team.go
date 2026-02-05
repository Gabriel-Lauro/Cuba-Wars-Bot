package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func JoinTeam(s *discordgo.Session, db *sql.DB, nomePlayer string, nomeTime string, idDiscord string, maxPlayers int) (bool, string) {
	var timeID int

	err := db.QueryRow(`
		SELECT id FROM time WHERE nome = ?
	`, nomeTime).Scan(&timeID)

	if err != nil {
		return false, "Time não encontrado"
	}

	var currentTeamID sql.NullInt64
	err = db.QueryRow(`
		SELECT time_id FROM players WHERE id_discord = ?
	`, idDiscord).Scan(&currentTeamID)

	if err != nil && err != sql.ErrNoRows {
		return false, "Erro ao verificar time atual: " + err.Error()
	}

	if currentTeamID.Valid && currentTeamID.Int64 != 0 {
		return false, "Você já está em um time. Saia dele primeiro com /sair_time"
	}

	var playerCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM players WHERE time_id = ?
	`, timeID).Scan(&playerCount)

	if err != nil {
		return false, "Erro ao verificar jogadores do time: " + err.Error()
	}

	if playerCount >= maxPlayers {
		return false, "O time está cheio"
	}

	var existingNick string
	err = db.QueryRow(`
		SELECT nick_jogo FROM players WHERE nick_jogo = ? AND time_id = ?
	`, nomePlayer, timeID).Scan(&existingNick)

	if err == nil {
		return false, "Esse nick já existe neste time"
	}

	if err != sql.ErrNoRows {
		return false, "Erro ao verificar nick: " + err.Error()
	}

	// Verificar se o jogador já existe no banco
	var existingPlayer string
	err = db.QueryRow(`
		SELECT id_discord FROM players WHERE id_discord = ?
	`, idDiscord).Scan(&existingPlayer)

	if err == nil {
		// Jogador já existe, fazer UPDATE
		_, err = db.Exec(`
			UPDATE players SET nick_jogo = ?, time_id = ? WHERE id_discord = ?
		`, nomePlayer, timeID, idDiscord)

		if err != nil {
			return false, "Erro ao entrar no time: " + err.Error()
		}
	} else if err == sql.ErrNoRows {
		// Jogador não existe, fazer INSERT
		_, err = db.Exec(`
			INSERT INTO players (id_discord, nick_jogo, time_id)
			VALUES (?, ?, ?)
		`, idDiscord, nomePlayer, timeID)

		if err != nil {
			return false, "Erro ao entrar no time: " + err.Error()
		}
	} else {
		return false, "Erro ao verificar jogador: " + err.Error()
	}

	SendPlayerMessage(s, idDiscord, db, nomePlayer)
	return true, "Entrou no time com sucesso"
}
