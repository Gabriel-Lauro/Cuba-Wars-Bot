package commands

import (
	"database/sql"
)

func KickPlayer(db *sql.DB, nomePlayer string, idDiscord string) (bool, string) {
	var timeID int
	var lider string
	err := db.QueryRow(`
		SELECT id, lider FROM time WHERE lider = ?
	`, idDiscord).Scan(&timeID, &lider)

	if err != nil {
		return false, "Você não lidera nenhum time"
	}

	var playerID string
	var playerLiderID string
	err = db.QueryRow(`
		SELECT id_discord, COALESCE((SELECT lider FROM time WHERE id = ?), '') as lider_id
		FROM players WHERE nick_jogo = ? AND time_id = ?
	`, timeID, nomePlayer, timeID).Scan(&playerID, &playerLiderID)

	if err != nil {
		return false, "Jogador não encontrado no seu time"
	}

	// Impedir que o líder seja expulso
	if playerID == lider {
		return false, "Você não pode expulsar o líder do time"
	}

	_, err = db.Exec(`
		DELETE FROM players 
		WHERE nick_jogo = ? AND time_id = ?
	`, nomePlayer, timeID)

	if err != nil {
		return false, "Erro ao expulsar jogador: " + err.Error()
	}

	return true, "Jogador expulso com sucesso"
}
