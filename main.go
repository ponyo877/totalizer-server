package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ponyo877/totalizer-server/repository"
	"github.com/ponyo877/totalizer-server/usecase/session"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type ReqMsg struct {
	Type string `json:"type"`
}

type ResMsg struct {
	Value int `json:"value"`
}

func NewResMsg(value int) *ResMsg {
	return &ResMsg{value}
}

func wsConnection(service *session.Service) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			var req ReqMsg
			if err := websocket.JSON.Receive(ws, &req); err != nil {
				log.Printf("Receive failed: %s; closing connection...", err.Error())
				if err = ws.Close(); err != nil {
					log.Println("Error closing connection:", err.Error())
				}
				break
			}
			questions, _ := service.ListQuestion()
			log.Printf("Questions: %v\n", questions)
			value, err := service.Incriment("counter")
			if err != nil {
				log.Printf("Error incrementing value: %s\n", err.Error())
			}
			if err = websocket.JSON.Send(ws, NewResMsg(value)); err != nil {
				log.Printf("Send failed: %s; closing connection...", err.Error())
				if err = ws.Close(); err != nil {
					log.Println("Error closing connection:", err.Error())
				}
				break
			}
		}
	}
}

func main() {
	flag.Parse()
	redisURL := os.Getenv("REDIS_URL")
	redisToken := os.Getenv("REDIS_TOKEN")

	opt, _ := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:6379", redisToken, redisURL))
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pgbouncer=true sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Printf("dsn: %s\n", dsn)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	orgdb, _ := db.DB()
	defer orgdb.Close()
	repository := repository.NewSessionRepository(db, redis.NewClient(opt))
	service := session.NewService(repository)

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		websocket.Server{Handler: websocket.Handler(wsConnection(service))}.ServeHTTP(w, req)
	})

	log.Printf("Server listening on port %d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
