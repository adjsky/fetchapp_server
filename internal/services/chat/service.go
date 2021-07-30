package chat

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/adjsky/fetchapp_server/internal/models/user/userauth"

	"github.com/adjsky/fetchapp_server/internal/services"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

const (
	pongWait   = time.Second * 60
	pingPeriod = pongWait * 2 / 3
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type clientData struct {
	Email string
}

type chatService struct {
	clients       map[*websocket.Conn]clientData
	clientsSync   sync.RWMutex
	broadcastChan chan string
}

func NewService() services.Service {
	serv := chatService{
		clients:       make(map[*websocket.Conn]clientData),
		broadcastChan: make(chan string),
	}
	go serv.broadcast()
	return &serv
}

func (serv *chatService) Register(r *gin.RouterGroup) {
	r.GET("/ws", serv.handleWebsocket)
}

func (serv *chatService) Close() {
	for client := range serv.clients {
		client.Close()
	}
}

func (serv *chatService) broadcast() {
	message := <-serv.broadcastChan
	for client := range serv.clients {
		client.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (serv *chatService) handleWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	claims, _ := c.Get(userauth.ClaimsKey)
	userClaims, _ := claims.(*userauth.Claims)
	clientEmail := userClaims.Email
	log.Println("New client:", clientEmail)
	serv.clientsSync.Lock()
	serv.clients[conn] = clientData{
		Email: clientEmail,
	}
	serv.clientsSync.Unlock()
	go serv.writer(conn)
	serv.reader(conn)
}

func (serv *chatService) writer(conn *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	for {
		<-pingTicker.C
		err := conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			conn.SetReadDeadline(time.Now())
			break
		}
	}
}

func (serv *chatService) reader(conn *websocket.Conn) {
	defer conn.Close()
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		mType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			log.Println("Client left:", serv.clients[conn].Email)
			serv.clientsSync.Lock()
			delete(serv.clients, conn)
			serv.clientsSync.Unlock()
			break
		}
		if mType == websocket.TextMessage {
			serv.broadcastChan <- string(message)
		}
	}
}
