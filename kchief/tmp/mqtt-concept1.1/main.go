package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewMQTTClient if successful will return a new client
// connection to the broker with options set
func NewMQTTClient(protocol string, address string, port string) (mqtt.Client, error) {
	// -- Set default connect options
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(protocol + "://" + address + ":" + port).SetClientID("btclient")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	// Create the connection to the broker
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

func subscribe(c mqtt.Client) error {
	// Since the Subscribe method uses a callback function
	// for what to do with the message, we declare such a
	// method to print out the messages we receive.
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("*** Got message *** [%s] %s\n", msg.Topic(), string(msg.Payload()))
	}

	// Start the consuming of the topic
	if token := c.Subscribe("mytopic", 0, messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func publish(c mqtt.Client, msgCh chan interface{}, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-done:
			wg.Done()
			return
		default:
			// Create a string with the data to publish to broker
			text := fmt.Sprintf("this is msg #%v!", <-msgCh)
			token := c.Publish("mytopic", 0, false, text)
			token.Wait()
		}
	}
}

// unSubscribe will unsubscribe the client from the topic
func unSubscribe(c mqtt.Client) error {
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func main() {
	var err error

	// Create new mqtt client
	client, err := NewMQTTClient("tcp", "10.0.0.26", "1883")
	if err != nil {
		log.Printf("error: newMQTTClient failed: %v\n", err)
		return
	}
	defer client.Disconnect(250)

	// Subscribe to topic,
	// subscribe will also print the result to console.
	err = subscribe(client)
	if err != nil {
		log.Printf("error: mqtt client subscribe failed: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	msgCh := make(chan interface{})
	doneCh := make(chan struct{})

	//start publish'er
	wg.Add(0)
	go publish(client, msgCh, doneCh, &wg)

	// send some data to the publisher channel
	for i := 0; i < 5; i++ {
		msgCh <- i
		time.Sleep(time.Second * 1)
	}

	// Wait for publisher go routine to finish
	wg.Wait()

	err = unSubscribe(client)
	if err != nil {
		log.Printf("error: mqtt client unSubscribe failed: %v\n", err)
		return
	}

}
