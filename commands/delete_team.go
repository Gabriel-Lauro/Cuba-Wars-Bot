package commands

import (
	"database/sql"
)

func DeleteTeam(db *sql.DB, nomeTime string, idDiscord string) (bool, string) {
	var timeID int
	var lider string

	err := db.QueryRow(`
		SELECT id, lider FROM time WHERE nome = ?
	`, nomeTime).Scan(&timeID, &lider)

	if err != nil {
		return false, "Time não encontrado"
	}

	if lider != idDiscord {
		return false, "Apenas o líder pode excluir o time"
	}

	_, _ = db.Exec(`
		DELETE FROM players WHERE time_id = ?
	`, timeID)

	_, err = db.Exec(`
		DELETE FROM time WHERE id = ?
	`, timeID)

	if err != nil {
		return false, "Erro ao excluir time: " + err.Error()
	}

	return true, "Time excluído com sucesso"
}
