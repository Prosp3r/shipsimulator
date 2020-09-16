package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewMQTTClient if successful will return a new client
// connection to the broker with options set
func NewMQTTClient(protocol string, address string, port string, clientID string) (mqtt.Client, error) {
	// -- Set default connect options
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(protocol + "://" + address + ":" + port).SetClientID(clientID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	// Create the connection to the broker.
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

// fileNR keeps track of the numbering of output files.
type fileNR struct {
	number int
}

// newMessageHandler is a wrapper to implement the numbering
// in the out filename created. It will setup a new MessageHandler
// and return the wrap'ed handler to the caller.
// The MessageHandler is used as a callback function to
//
// be called for each MQTT packet received.
// A new file will be created upon each call and the
// number in the filename is incremented by 1 on each call.
func newMessageHandler() func(client mqtt.Client, msg mqtt.Message) {
	fn := &fileNR{}
	fn.number = 1

	// Since the Subscribe method uses a callback function
	// for what to do with the message, we declare such a
	// function to print out the messages we receive.
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		// Create and open a file with next running number
		fileName := fmt.Sprintf("out%v.bin", fn.number)
		f, err := os.Create(fileName)
		if err != nil {
			log.Printf("error: open file for output: %vn", err)
			return
		}
		defer f.Close()

		// Write the Protobuf data to the file.
		n, err := f.Write(msg.Payload())
		if err != nil {
			log.Printf("error: writing to file: %v\n", err)
			log.Printf("info: characters written: %v\n", n)
		}

		fn.number++
	}

	return messageHandler
}

func subscribe(c mqtt.Client, topic string) error {
	// The Subscribe method uses a callback function
	// for what to do with the message.
	// Create a new message handler.
	messageHandler := newMessageHandler()

	// Start the consuming of the topic
	if token := c.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// unSubscribe will unsubscribe the client from the topic
func unSubscribe(c mqtt.Client, topic string) error {
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func main() {
	topic := flag.String("topic", "CloudBoundContainer", "The name of the MQTT topic")
	broker := flag.String("broker", "10.0.0.26", "The ip address of the MQTT broker")
	port := flag.String("port", "1883", "The port where the MQTT broker listens")
	protocol := flag.String("protocol", "tcp", "The protocol to use when connecting to the MQTT broker")
	clientID := flag.String("clientID", "btclient2", "The client ID to use with MQTT")

	flag.Parse()

	var err error

	// Create new mqtt client.
	client, err := NewMQTTClient(*protocol, *broker, *port, *clientID)
	if err != nil {
		log.Printf("error: newMQTTClient failed: %v\n", err)
		return
	}
	defer client.Disconnect(250)

	// Subscribe to topic,
	// subscribe will also print the result to console.
	err = subscribe(client, *topic)
	if err != nil {
		log.Printf("error: mqtt client subscribe failed: %v\n", err)
		return
	}

	// Wait for someone to press CTRL+C.
	fmt.Println("started subscriber")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("stopped subscriber")

	// Stop subscribing to MQTT messages.
	err = unSubscribe(client, *topic)
	if err != nil {
		log.Printf("error: mqtt client unSubscribe failed: %v\n", err)
		return
	}

}
