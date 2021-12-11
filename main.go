package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/olahol/melody.v1"
)

// GopherInfo contains information about the gopher on screen
type UserSession struct {
	ID string
}

func main() {
	port := os.Getenv("PORT")
	router := gin.Default()
	mrouter := melody.New()
	mrouter.Config = &melody.Config{
		WriteWait:         mrouter.Config.WriteWait,
		PongWait:          mrouter.Config.PongWait,
		PingPeriod:        (time.Duration(5) * time.Second),
		MaxMessageSize:    mrouter.Config.MaxMessageSize,
		MessageBufferSize: mrouter.Upgrader.ReadBufferSize,
	}
	mrouter.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	userSession := make(map[*melody.Session]*UserSession)
	lock := new(sync.Mutex)
	counter := 0

	router.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	router.GET("/ws-pingpong-latency", func(c *gin.Context) {
		mrouter.HandleRequest(c.Writer, c.Request)
	})

	mrouter.HandlePong(func(s *melody.Session) {
		s.Write([]byte(fmt.Sprint(websocket.PingMessage)))
	})

	mrouter.HandleConnect(func(s *melody.Session) {
		lock.Lock()
		fmt.Println("client connected", s.Request.UserAgent())
		userSession[s] = &UserSession{strconv.Itoa(counter)}
		counter++
		lock.Unlock()
		s.Write([]byte("sessionID: " + userSession[s].ID))
	})

	mrouter.HandleDisconnect(func(s *melody.Session) {
		lock.Lock()
		delete(userSession, s)
		lock.Unlock()
	})

	mrouter.HandleMessage(func(s *melody.Session, msg []byte) {
		userTime, err := strconv.ParseInt(string(msg), 10, 64)
		if err != nil {
			panic(err)
		}
		serverTime := time.Now().UnixNano() / int64(time.Millisecond)
		s.Write([]byte("latency: " + fmt.Sprint(serverTime-userTime)))
	})

	router.Run(":" + port)
}
