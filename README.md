# 🤖 Cuba Wars Bot

Bot Discord para gerenciamento de times.

## Requisitos

- Go 1.25.5+
- Token Discord Bot

## Setup

1. Clone o repositório
2. Renomeie `.env.example` para `.env` e configure:
   ```
   DISCORD_TOKEN=seu_token_aqui
   MAX_TEAM_PLAYERS=5
   ```
3. Execute:
   ```bash
   go run main.go
   ```

## Comandos

| Comando             | Descrição                        | Restrição         |
| ------------------- | -------------------------------- | ----------------- |
| `/criar_time`       | Cria um novo time                | -                 |
| `/entrar_time`      | Entra em um time                 | Máx. de jogadores |
| `/sair_time`        | Sai do time atual                | -                 |
| `/expulsar_jogador` | Remove um jogador do time        | Apenas líder      |
| `/info_time`        | Mostra informações do time       | -                 |
| `/listar_times`     | Lista todos os times disponíveis | -                 |
| `/deletar_time`     | Deleta o time                    | Apenas líder      |

## Banco de Dados

SQLite com tabelas:
- `time`: id, nome, lider
- `players`: id_discord, nick_jogo, op, time_id, senha

## Dependências

- discordgo
- sqlite (modernc.org)
- godotenv
