package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lichtwellenreiter/chsbticker/jsonstruct"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	setupRoutes()
	log.Println(http.ListenAndServe("0.0.0.0:8080", nil))
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func ticker(conn *websocket.Conn) {

	for {

		btc := getHitBtcData("https://api.hitbtc.com/api/2/public/ticker/BTCUSD")
		chsb := getHitBtcData("https://api.hitbtc.com/api/2/public/ticker/CHSBBTC")
		forex := getForexData("https://api.exchangeratesapi.io/latest?symbols=USD,CHF&base=USD")

		btcusd, _ := strconv.ParseFloat(btc.Last, 64)
		chsbbtc, _ := strconv.ParseFloat(chsb.Last, 64)

		chf := (btcusd * chsbbtc) * forex.Rates.CHF

		tickerdata := &jsonstruct.TickerData{
			Timestamp: time.Now().Format("01.02.2006 15:04:05"),
			Chsb:      chsb.Last,
			Chf:       fmt.Sprintf("%f", chf),
		}

		_, err := json.Marshal(tickerdata)
		if err != nil {

			log.Println(err)
		}

		err = conn.WriteJSON(tickerdata)
		if err != nil {
			log.Println("NO LONGER")
			log.Println(err)
			conn.Close()
			break

		}

		//log.Println(string(b))
		time.Sleep(1 * time.Second)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection

	ticker(ws)
	// reader(ws)
}

func setupRoutes() {
	fs := http.FileServer(http.Dir("/templates"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", wsEndpoint)
}

func getForexData(url string) jsonstruct.ForexExchange {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject jsonstruct.ForexExchange
	_ = json.Unmarshal(bodyBytes, &responseObject)
	// log.Printf("API Response as struct %+v\n", responseObject)
	return responseObject
}

func getHitBtcData(url string) jsonstruct.Ticker {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject jsonstruct.Ticker
	_ = json.Unmarshal(bodyBytes, &responseObject)
	// log.Printf("API Response as struct %+v\n", responseObject)
	return responseObject
}
