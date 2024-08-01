package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ponyo877/totalizer-server/repository"
	"github.com/ponyo877/totalizer-server/usecase/session"
	socket "github.com/ponyo877/totalizer-server/websocket"
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

type EnterMsg struct {
	RoomNumber string `json:"room_number"`
}

type AskMsg struct {
	RoomID   string `json:"room_id"`
	Question string `json:"question"`
}

type VoteMsg struct {
	RoomID     string `json:"room_id"`
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type ReleaseMsg struct {
	RoomID     string `json:"room_id"`
	QuestionID string `json:"question_id"`
}

func NewResMsg(value int) *ResMsg {
	return &ResMsg{value}
}

func parseMsg[T any](msg []byte) (T, error) {
	var anyMsg T
	if err := json.Unmarshal(msg, &anyMsg); err != nil {
		log.Printf("Unmarshal failed: %s\n", err.Error())
		return anyMsg, err
	}
	return anyMsg, nil
}

// A: Open ->       -> Ask ->      -> Vote ->
// B:      -> Enter ->     -> Vote ->      -> Release
func wsConnection(service session.UseCase) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		s := socket.NewSocket(ws, service)
		for {
			var msg []byte
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				log.Printf("Receive failed: %s; closing connection...", err.Error())
				if err = ws.Close(); err != nil {
					log.Println("Error closing connection:", err.Error())
				}
				break
			}
			var req ReqMsg
			if err := json.Unmarshal(msg, &req); err != nil {
				log.Printf("Unmarshal failed: %s\n", err.Error())
			}
			switch req.Type {
			// 開室
			case "open":
				s.Open()
			// 入室
			case "enter":
				enterMsg, err := parseMsg[EnterMsg](msg)
				if err != nil {
					log.Printf("Message failed: %s\n", err)
					break
				}
				s.Enter(enterMsg.RoomNumber)
			// 出題
			case "ask":
				askMsg, err := parseMsg[AskMsg](msg)
				if err != nil {
					log.Printf("Message failed: %s\n", err)
					break
				}
				s.Ask(askMsg.RoomID, askMsg.Question)
			// 投票
			case "vote":
				voteMsg, err := parseMsg[VoteMsg](msg)
				if err != nil {
					log.Printf("Message failed: %s\n", err)
					break
				}
				s.Vote(voteMsg.RoomID, voteMsg.QuestionID, voteMsg.Answer)
			// 公表
			case "release":
				releaseMsg, err := parseMsg[ReleaseMsg](msg)
				if err != nil {
					log.Printf("Message failed: %s\n", err)
					break
				}
				s.Release(releaseMsg.RoomID, releaseMsg.QuestionID)
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
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Printf("dsn: %s\n", dsn)
	db, _ := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	orgdb, _ := db.DB()
	defer orgdb.Close()
	repository := repository.NewSessionRepository(db, redis.NewClient(opt))
	service := session.NewService(repository)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		websocket.Server{Handler: websocket.Handler(wsConnection(service))}.ServeHTTP(w, req)
	})

	log.Printf("Server listening on port %d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
