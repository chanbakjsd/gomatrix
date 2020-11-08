package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/chanbakjsd/gotrix"
	"github.com/chanbakjsd/gotrix/event"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func handleMessage(c *gotrix.Client, m event.RoomMessageEvent) {
	// If it's a notice (another bot's message) or not "ping", ignore.
	if m.MsgType == event.RoomMessageNotice || m.Body != "ping" {
		return
	}
	// Otherwise, send pong!
	_, err := c.SendNotice(m.RoomID, "Pong!")
	panicIfErr(err)
}

func main() {
	// Ask for username and password.
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <username> <password>", os.Args[0])
		return
	}

	// Construct the client.
	cli, err := gotrix.New("http://localhost:8008")
	panicIfErr(err)

	// Login using provided creds.
	panicIfErr(cli.LoginPassword(os.Args[1], os.Args[2]))

	// Register the handler.
	cli.AddHandler(handleMessage)

	// Start the connection.
	panicIfErr(cli.Open())

	// Wait until interrupt happens.
	fmt.Println("Ctrl-C to terminate program.")
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate

	// Close the connection and logout as we're not persisting the session.
	panicIfErr(cli.Close())
	panicIfErr(cli.Logout())
}
