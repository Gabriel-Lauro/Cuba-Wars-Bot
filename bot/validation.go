package bot

import (
	"regexp"
	"strings"
)

// ValidatePlayerNick valida o nick do jogador de acordo com as regras do Minecraft
func ValidatePlayerNick(nick string) (bool, string) {
	// Verificar tamanho
	if len(nick) < 3 || len(nick) > 16 {
		return false, "O nick do player deve ter entre 3 e 16 caracteres."
	}

	// Verificar caracteres permitidos (A-Z, a-z, 0-9, _)
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validChars.MatchString(nick) {
		return false, "O nick só pode conter letras (A-Z), números (0-9) e underscore (_). Sem espaços, acentos ou símbolos especiais."
	}

	// Verificar se começa ou termina com underscore
	if strings.HasPrefix(nick, "_") || strings.HasSuffix(nick, "_") {
		return false, "O nick não pode começar ou terminar com underscore (_)."
	}

	// Verificar se tem dois underscores seguidos
	if strings.Contains(nick, "__") {
		return false, "O nick não pode ter dois underscores seguidos (__)."
	}

	return true, ""
}

// ValidateTeamName valida o nome do time de acordo com as regras
func ValidateTeamName(name string) (bool, string) {
	// Verificar tamanho
	if len(name) < 2 || len(name) > 16 {
		return false, "O nome do time deve ter entre 3 e 16 caracteres."
	}

	// Verificar caracteres permitidos (A-Z, a-z, 0-9, _)
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validChars.MatchString(name) {
		return false, "O nome do time só pode conter letras (A-Z), números (0-9) e underscore (_). Sem espaços, acentos ou símbolos especiais."
	}

	// Verificar se começa ou termina com underscore
	if strings.HasPrefix(name, "_") || strings.HasSuffix(name, "_") {
		return false, "O nome do time não pode começar ou terminar com underscore (_)."
	}

	// Verificar se tem dois underscores seguidos
	if strings.Contains(name, "__") {
		return false, "O nome do time não pode ter dois underscores seguidos (__)."
	}

	return true, ""
}

