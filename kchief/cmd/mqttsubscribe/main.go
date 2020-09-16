package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	// Create the connection to the broker
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

func subscribe(c mqtt.Client, topic string) error {
	// Since the Subscribe method uses a callback function
	// for what to do with the message, we declare such a
	// method to print out the messages we receive.
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("*** Got message *** [%s] %s\n", msg.Topic(), string(msg.Payload()))
	}

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
	clientID := flag.String("clientID", "btclient", "The client ID to use with MQTT")

	flag.Parse()

	var err error

	// Create new mqtt client
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

	err = unSubscribe(client, *topic)
	if err != nil {
		log.Printf("error: mqtt client unSubscribe failed: %v\n", err)
		return
	}

}
