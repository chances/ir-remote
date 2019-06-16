package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	messages := make(chan int)

	go uiServer(messages)

	// Gobot
	r := raspi.NewAdaptor()
	ir := gpio.NewDirectPinDriver(r, "7")

	work := func() {
		oldVal := 0

		gobot.Every(1*time.Millisecond, func() {
			irVal, err := ir.DigitalRead()
			if err != nil {
				return
			}

			if oldVal != irVal {
				log.Print("IR val: ")
				log.Println(irVal)

				messages <- irVal

				oldVal = irVal
			}
		})
	}

	robot := gobot.NewRobot("irRemoteBot",
		[]gobot.Connection{r},
		[]gobot.Device{ir},
		work,
	)

	log.Fatal(robot.Start())
}

func uiServer(messages chan int) {
	irVal := 0

	// UI Web Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, world (from go!)\n"))
		w.Write([]byte("ir value: "))
		w.Write([]byte(string(irVal)))
		w.Write([]byte("\n"))
	})

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  256,
		WriteBufferSize: 256,
	}

	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			irVal = <-messages
			data := []byte(string(irVal))
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				return
			}
		}
	})

	addr := ":80"
	fmt.Println("Example app listening on port ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
