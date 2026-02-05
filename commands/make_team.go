package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func MakeTeam(s *discordgo.Session, db *sql.DB, nome string, nomePlayer string, liderID string) (bool, string) {
	var existingTeam string
	err := db.QueryRow(`
		SELECT nome FROM time WHERE lider = ?
	`, liderID).Scan(&existingTeam)

	if err == nil {
		return false, "Você já lidera o time: " + existingTeam
	}

	if err != sql.ErrNoRows {
		return false, "Erro ao verificar times: " + err.Error()
	}

	var currentTeamID sql.NullInt64
	err = db.QueryRow(`
		SELECT time_id FROM players WHERE id_discord = ?
	`, liderID).Scan(&currentTeamID)

	if err != nil && err != sql.ErrNoRows {
		return false, "Erro ao verificar time atual: " + err.Error()
	}

	if currentTeamID.Valid && currentTeamID.Int64 != 0 {
		return false, "Você já está em um time. Saia dele primeiro"
	}

	result, err := db.Exec(`
		INSERT INTO time (nome, lider)
		VALUES (?, ?)
	`, nome, liderID)

	if err != nil {
		return false, "Erro ao criar time (nome pode estar em uso)"
	}

	timeID, err := result.LastInsertId()
	if err != nil {
		return false, "Erro ao obter ID do time: " + err.Error()
	}

	_, err = db.Exec(`
		INSERT INTO players (id_discord, nick_jogo, time_id)
		VALUES (?, ?, ?)
	`, liderID, nomePlayer, timeID)

	if err != nil {
		db.Exec(`DELETE FROM time WHERE id = ?`, timeID)
		return false, "Erro ao adicionar você ao time: " + err.Error()
	}

	SendPlayerMessage(s, liderID, db, nomePlayer)
	return true, "Time criado com sucesso e você entrou como jogador"
}
