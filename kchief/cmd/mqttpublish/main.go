package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

func publish(c mqtt.Client, topic string, msgCh chan interface{}) {
	for {
		// Create a string with the data to publish to broker
		text := fmt.Sprintf("this is msg #%v!", <-msgCh)
		token := c.Publish(topic, 0, false, text)
		token.Wait()
	}
}

// unSubscribe will unsubscribe the client from the topic
func unSubscribe(c mqtt.Client, topic string) error {
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type inFileData struct {
	bytes []byte
}

func main() {
	topic := flag.String("topic", "CloudBoundContainer", "The name of the MQTT topic")
	broker := flag.String("broker", "10.0.0.26", "The ip address of the MQTT broker")
	port := flag.String("port", "1883", "The port where the MQTT broker listens")
	protocol := flag.String("protocol", "tcp", "The protocol to use when connecting to the MQTT broker")
	clientID := flag.String("clientID", "btclient", "The client ID to use with MQTT")
	var inFile arrayFlags
	flag.Var(&inFile, "inFile", "specify the files to use as input comma separated")
	repetitions := flag.Int("repetitions", 1, "specify how many repetitions to run")
	delay := flag.Int("delay", 300, "The number of milliseconds to wait between each mqtt publish")

	flag.Parse()

	if len(inFile) == 0 {
		fmt.Printf("no input files, use the --inFile flag to specify one or mores files to use as input.\n")
		return
	}

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

	msgCh := make(chan interface{})

	//start publish'er
	go publish(client, *topic, msgCh)

	//get the binary data from the files
	inFilesData := []inFileData{}

	for _, v := range inFile {
		fmt.Printf("v contains : %v\n", v)
		fh, err := os.Open(v)
		if err != nil {
			log.Printf("error: failed to open protoFile: %v\n", err)
			return
		}

		b, err := ioutil.ReadAll(fh)
		if err != nil {
			log.Printf("error: file ReadAll failed: %v\n", err)
		}

		d := inFileData{
			bytes: b,
		}

		inFilesData = append(inFilesData, d)
		fh.Close()

	}

	fmt.Printf("*** %#v\n", len(inFilesData))

	// send some data to the publisher channel
	for i := 0; i < *repetitions; i++ {
		fmt.Printf("protoFiles contains : %#v\n", inFile)
		for _, v := range inFilesData {

			msgCh <- v.bytes
			time.Sleep(time.Millisecond * time.Duration(*delay))
		}
	}

	err = unSubscribe(client, *topic)
	if err != nil {
		log.Printf("error: mqtt client unSubscribe failed: %v\n", err)
		return
	}

}
