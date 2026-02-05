package commands

import (
	"database/sql"
)

func LeaveTeam(db *sql.DB, idDiscord string) (bool, string) {
	var timeID sql.NullInt64
	var isLider bool
	err := db.QueryRow(`
		SELECT p.time_id, CASE WHEN t.lider = ? THEN 1 ELSE 0 END as is_lider
		FROM players p
		LEFT JOIN time t ON p.time_id = t.id
		WHERE p.id_discord = ?
	`, idDiscord, idDiscord).Scan(&timeID, &isLider)

	if err != nil {
		return false, "Erro ao verificar seu time: " + err.Error()
	}

	if !timeID.Valid || timeID.Int64 == 0 {
		return false, "Você não está em nenhum time"
	}

	if isLider {
		return false, "O líder não pode sair do time."
	}

	res, err := db.Exec(`
		DELETE FROM players WHERE id_discord = ?
	`, idDiscord)

	if err != nil {
		return false, "Erro ao sair do time: " + err.Error()
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false, "Você não está em nenhum time"
	}

	return true, "Você saiu do time com sucesso"
}
