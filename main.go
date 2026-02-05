package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"	
	"database/sql"
	"Cuba-Wars-Bot/bot"
	_ "modernc.org/sqlite"
	
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func createDb(db *sql.DB) {
	_, err := db.Exec(`	
		CREATE TABLE IF NOT EXISTS time (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nome VARCHAR(32) NOT NULL UNIQUE COLLATE NOCASE,
			lider VARCHAR(32) NOT NULL COLLATE NOCASE
		);

		CREATE TABLE IF NOT EXISTS players (
			id_discord VARCHAR(32) PRIMARY KEY,
			nick_jogo VARCHAR(32) NOT NULL UNIQUE COLLATE NOCASE,
			op INTEGER DEFAULT 0,
			time_id INTEGER,
			senha VARCHAR(6),
			FOREIGN KEY (time_id) REFERENCES time(id)
		);`)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Garante que o diretório do banco exista
	if err := os.MkdirAll("db", 0755); err != nil {
		panic(err)
	}

	// Carrega o Db
	db, err := sql.Open("sqlite", "db/database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createDb(db)

	// Carrega o .env
	err = godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sessão:", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds

	err = dg.Open()
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}

	fmt.Println("🤖 Bot conectado com sucesso!")
	bot.RegisterCommands(dg, db)

	// Mantém rodando
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	dg.Close()
}