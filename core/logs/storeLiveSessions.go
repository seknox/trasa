package logs

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/seknox/trasa/utils"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

//AddNewActiveSession adds new sesssion to active session list in redis
func (s LogStore) AddNewActiveSession(session *AuthLog, connID, appType string) error {
	client := s.RedisClient
	ctx := context.Background()

	newSession, err := newActiveSession(session, connID, appType)

	if err != nil {
		return err
	}

	err = client.LPush(ctx, "activeSessionKeys", connID).Err()
	if err != nil {
		return errors.Errorf("push error: %v", err)
	}

	err = client.HSet(ctx, "activeSessions", connID, newSession).Err()
	if err != nil {
		return errors.Errorf("HSET error: %v", err)
	}
	err = client.Expire(ctx, "activeSessions", time.Second*86400).Err()
	if err != nil {
		return errors.Errorf("expire error: %v", err)
	}
	err = client.Expire(ctx, "activeSessionKeys", time.Second*86400).Err()

	if err != nil {
		return err
	}
	err = client.Publish(ctx, "updateActiveSession", "").Err()

	return err

}

func (s LogStore) RemoveActiveSession(connID string) error {
	client := s.RedisClient
	ctx := context.Background()

	//remove connID from list
	err := client.LRem(ctx, "activeSessionKeys", 0, connID).Err()
	if err != nil {
		return errors.Errorf("lrem error: %v", err)
	}

	err = client.HDel(ctx, "activeSessions", "0", connID).Err()
	if err != nil {
		return errors.Errorf("hrem error: %v", err)

	}

	err = client.Publish(ctx, "updateActiveSession", "").Err()

	return err
}

func (s LogStore) RemoveAllActiveSessions() {
	ctx := context.Background()
	sessions, err := s.getAllActiveSessions(ctx)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, session := range sessions {
		err = s.RemoveActiveSession(session)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (s LogStore) getAllActiveSessions(ctx context.Context) ([]string, error) {
	client := s.RedisClient

	list, err := client.LRange(ctx, "activeSessionKeys", 0, -1).Result()
	if err != nil || len(list) == 0 {
		return nil, err
	}

	data, err := client.HMGet(ctx, "activeSessions", list...).Result()
	if err != nil {
		return nil, err
	}

	//fmt.Println(data)
	values, err := utils.ToStringArr(data)
	if err != nil {
		return values, err
	}

	if values == nil {
		values = make([]string, 0)
	}
	return values, err

}

func (s LogStore) ServeLiveSessions(ws *websocket.Conn) {
	ctx := context.Background()
	sessions, err := s.getAllActiveSessions(ctx)
	if err != nil {
		logrus.Error(err)
		return
	}
	ws.WriteJSON(sessions)

	client := s.RedisClient

	var check = true

	psc := client.Subscribe(ctx, "updateActiveSession")
	//	defer psc.PSubscribe(ctx,"__key*__:*")

	var done chan bool
	done = make(chan bool)

	go receivePing(ws, done)

	go func() {
		<-done
		check = false
		psc.PUnsubscribe(ctx, "updateActiveSession")
		psc.Close()

	}()

	for check {
		//wait for change
		<-psc.Channel()
		sessions, _ := s.getAllActiveSessions(ctx)
		ws.WriteJSON(sessions)
	}

}

func receivePing(ws *websocket.Conn, done chan bool) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		read := make(chan string)

		go func(read chan string) {
			_, _, err := ws.ReadMessage()
			if err != nil {
				logrus.Debug(err)
				done <- false
				return
			}
			read <- "msg"
		}(read)

		select {
		case <-ticker.C:
			done <- false
			return
		case <-read:
			ticker.Stop()
			ticker = time.NewTicker(pingPeriod)
		}
	}
}
