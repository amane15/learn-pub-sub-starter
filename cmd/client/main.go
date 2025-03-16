package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	// "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	// "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Println("Error getting username")
		return
	}

	gameState := gamelogic.NewGameState(username)

	gamelogic.PrintClientHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		}
	}

	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey, 2)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
		return
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
