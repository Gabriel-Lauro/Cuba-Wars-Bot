package bot

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"Cuba-Wars-Bot/commands"
	"github.com/bwmarrin/discordgo"
)

type CommandDef struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(*discordgo.Session, *discordgo.InteractionCreate, *sql.DB) (bool, string)
}

func GetCommands() []CommandDef {
return []CommandDef{
{
Name:        "criar_time",
Description: "Cria um novo time",
Options: []*discordgo.ApplicationCommandOption{
{
Name:        "nome_do_time",
Description: "Nome do time (2-16 caracteres, apenas letras, números e underscore)",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
{
Name:        "nick_jogo",
Description: "Seu nick no jogo (3-16 caracteres, apenas letras, números e underscore)",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
},
Handler: handleMakeTeam,
},
{
Name:        "entrar_time",
Description: "Entra em um time existente",
Options: []*discordgo.ApplicationCommandOption{
{
Name:        "nick_jogo",
Description: "Seu nick no jogo (3-16 caracteres, apenas letras, números e underscore)",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
{
Name:        "nome_do_time",
Description: "Nome do time que deseja entrar",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
},
Handler: handleJoinTeam,
},
{
Name:        "sair_time",
Description: "Sai do time atual",
Options:     []*discordgo.ApplicationCommandOption{},
Handler:     handleLeaveTeam,
},
{
Name:        "expulsar_jogador",
Description: "Remove um jogador do time (apenas líder)",
Options: []*discordgo.ApplicationCommandOption{
{
Name:        "nick_jogo",
Description: "Nick do jogador a ser removido",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
},
Handler: handleKickPlayer,
},
{
Name:        "info_time",
Description: "Mostra informações do time",
Options: []*discordgo.ApplicationCommandOption{
{
Name:        "nome_do_time",
Description: "Nome do time para obter informações",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
},
Handler: handleTeamDetails,
},
{
Name:        "deletar_time",
Description: "Deleta um time (apenas líder)",
Options: []*discordgo.ApplicationCommandOption{
{
Name:        "nome_do_time",
Description: "Nome do time a ser deletado",
Type:        discordgo.ApplicationCommandOptionString,
Required:    true,
},
},
Handler: handleDeleteTeam,
},
{
Name:        "listar_times",
Description: "Lista todos os times com a quantidade de jogadores",
Options:     []*discordgo.ApplicationCommandOption{},
Handler:     handleListTeams,
},
}
}

func RegisterCommands(session *discordgo.Session, db *sql.DB) error {
	commandDefs := GetCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandDefs))
	for i, cmd := range commandDefs {
		appCmd := &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     cmd.Options,
		}
		registeredCmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", appCmd)
		if err != nil {
			return fmt.Errorf("erro ao registrar comando %s: %v", cmd.Name, err)
		}
		registeredCommands[i] = registeredCmd
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		for _, cmd := range commandDefs {
			if cmd.Name == i.ApplicationCommandData().Name {
				success, message := cmd.Handler(s, i, db)

				emoji := "✅"
				if !success {
					emoji = "❌"
				}

				response := fmt.Sprintf("%s %s", emoji, message)

				// Responde à interação
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
					},
				})
				return
			}
		}
	})

	return nil
}

func handleMakeTeam(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	options := i.ApplicationCommandData().Options
	teamName := options[0].StringValue()
	playerNick := options[1].StringValue()

	valid, errMsg := ValidateTeamName(teamName)
	if !valid {
		return false, errMsg
	}

	valid, errMsg = ValidatePlayerNick(playerNick)
	if !valid {
		return false, errMsg
	}

	return commands.MakeTeam(s, db, teamName, playerNick, i.Member.User.ID)
}

func handleJoinTeam(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	options := i.ApplicationCommandData().Options
	playerNick := options[0].StringValue()
	teamName := options[1].StringValue()

	valid, errMsg := ValidatePlayerNick(playerNick)
	if !valid {
		return false, errMsg
	}

	maxPlayers := 5
	if val := os.Getenv("MAX_TEAM_PLAYERS"); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			maxPlayers = n
		}
	}

	return commands.JoinTeam(s, db, playerNick, teamName, i.Member.User.ID, maxPlayers)
}

func handleLeaveTeam(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	return commands.LeaveTeam(db, i.Member.User.ID)
}

func handleKickPlayer(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	options := i.ApplicationCommandData().Options
	playerNick := options[0].StringValue()

	return commands.KickPlayer(db, playerNick, i.Member.User.ID)
}

func handleTeamDetails(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	options := i.ApplicationCommandData().Options
	teamName := options[0].StringValue()

	return commands.TeamDetails(db, teamName)
}

func handleDeleteTeam(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	options := i.ApplicationCommandData().Options
	teamName := options[0].StringValue()

	return commands.DeleteTeam(db, teamName, i.Member.User.ID)
}

func handleListTeams(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) (bool, string) {
	return commands.ListTeams(db)
}
