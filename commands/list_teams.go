package commands

import (
	"database/sql"
	"fmt"
	"strings"
)

func ListTeams(db *sql.DB) (bool, string) {
	rows, err := db.Query(`
		SELECT t.nome, COUNT(p.id_discord) as players
		FROM time t
		LEFT JOIN players p ON p.time_id = t.id
		GROUP BY t.id, t.nome
		ORDER BY t.nome ASC
	`)
	if err != nil {
		return false, "Erro ao listar times"
	}
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var nome string
		var players int
		if err := rows.Scan(&nome, &players); err != nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("**%s** [%d]", nome, players))
	}

	if len(lines) == 0 {
		return true, "Nenhum time encontrado"
	}

	sep := "\n------\n"
	return true, strings.Join(lines, sep)
}
