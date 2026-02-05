package commands

import (
	"database/sql"
	"fmt"
)

func TeamDetails(db *sql.DB, nomeTime string) (bool, string) {
	var timeID int
	var lider string

	err := db.QueryRow(`
		SELECT id, lider FROM time WHERE nome = ?
	`, nomeTime).Scan(&timeID, &lider)

	if err != nil {
		return false, "Time não encontrado"
	}

	rows, err := db.Query(`
		SELECT nick_jogo FROM players WHERE time_id = ?
	`, timeID)
	if err != nil {
		return false, "Erro ao buscar jogadores"
	}
	defer rows.Close()

	details := fmt.Sprintf("Time: %s\nJogadores:", nomeTime)
	count := 0
	for rows.Next() {
		var nick string
		if err := rows.Scan(&nick); err == nil {
			details += "\n- " + nick
			count++
		}
	}
	if count == 0 {
		details += "\n(nenhum)"
	}

	return true, details
}
