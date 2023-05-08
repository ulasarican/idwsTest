package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	numClients = 500
)

func randomDuration(min, max int, r *rand.Rand) time.Duration {
	return time.Duration(r.Intn(max-min+1)+min) * time.Millisecond
}

func handleIncomingMessages(c *websocket.Conn) {
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Received message: %s", message)
	}
}

func simulateClient(wg *sync.WaitGroup, baseURL, roomID string, isAdmin bool, r *rand.Rand) {
	defer wg.Done()

	u := url.URL{Scheme: "ws", Host: baseURL, Path: "/ws"}
	_ = u
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.identify.com.tr", nil)

	if err != nil {
		log.Printf("Error connecting to WebSocket server: %v", err)
		return
	}
	defer c.Close()

	err = c.WriteJSON(map[string]interface{}{
		"action":   "subscribe",
		"room":     roomID,
		"is_admin": isAdmin,
	})

	if err != nil {
		log.Printf("Error subscribing to room: %v", err)
		return
	} else {
		log.Println("Subscribed to " + roomID)
	}

	go handleIncomingMessages(c)

	for i := 0; i < 10; i++ {
		time.Sleep(randomDuration(500, 2000, r))

		err = c.WriteJSON(map[string]interface{}{
			"action": "uploadIdFront",
			"result": true,
			"room":   roomID,
			"data":   fmt.Sprintf("Hello from %s (%s)", roomID, "admin"),
		})

		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		} else {
			log.Println(map[string]interface{}{
				"action": "uploadIdFront",
				"result": true,
				"room":   roomID,
				"data":   fmt.Sprintf("Hello from %s (%s)", roomID, "admin"),
			})
		}
	}

	time.Sleep(randomDuration(500, 2000, r))
	c.Close()
}

func main() {
	baseURL := "ws.identify.com.tr:443"

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	wg := &sync.WaitGroup{}
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		roomID := fmt.Sprintf("room-%v", uuid.NewString())
		go simulateClient(wg, baseURL, roomID, i%2 == 0, r)
		time.Sleep(500 * time.Millisecond) // Increase the delay between starting clients

	}

	wg.Wait()
}
