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

// publish will start the publishing to the MQTT queue,
// and it will take the input to publish via the msgCh.
func publish(c mqtt.Client, topic string, msgCh chan interface{}) {
	for {
		// Create a string with the data to publish to broker
		// text := fmt.Sprintf("this is msg #%v!", <-msgCh)
		token := c.Publish(topic, 0, false, <-msgCh)
		token.Wait()
	}
}

// inFileFlags allows us to use several --inFile flags
// when starting the program.
type inFileFlags []string

func (i *inFileFlags) String() string {
	return "my string representation"
}

func (i *inFileFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// inFileData represents the data of one of the files
// to read from giv'en as input.
type inFileData struct {
	bytes []byte
}

// getInFileData loops over the inFileFlags, reads the data
// of each individual file, appends the data in byte format
// to a slice, and returns it to the caller.
func getInFileData(fileNames inFileFlags) ([]inFileData, error) {
	inFilesData := []inFileData{}
	for _, v := range fileNames {
		fh, err := os.Open(v)
		if err != nil {
			return []inFileData{}, fmt.Errorf("error: failed to open protoFile: %v", err)
		}

		b, err := ioutil.ReadAll(fh)
		if err != nil {
			return []inFileData{}, fmt.Errorf("error: file ReadAll failed: %v", err)
		}

		d := inFileData{
			bytes: b,
		}

		inFilesData = append(inFilesData, d)
		fh.Close()

	}

	return inFilesData, nil
}

func main() {
	topic := flag.String("topic", "CloudBoundContainer", "The name of the MQTT topic")
	broker := flag.String("broker", "10.0.0.26", "The ip address of the MQTT broker")
	port := flag.String("port", "1883", "The port where the MQTT broker listens")
	protocol := flag.String("protocol", "tcp", "The protocol to use when connecting to the MQTT broker")
	clientID := flag.String("clientID", "btclient1", "The client ID to use with MQTT")
	var inFile inFileFlags
	flag.Var(&inFile, "inFile", "specify the files to use as input comma separated")
	repetitions := flag.Int("repetitions", 1, "specify how many repetitions to run")
	delay := flag.Int("delay", 300, "The number of milliseconds to wait between each mqtt publish")

	flag.Parse()

	// Check if inFiles for reading have been specified.
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

	// Create a channel to put the messages we want to publish.
	msgCh := make(chan interface{})

	// Start publish'er
	go publish(client, *topic, msgCh)

	inFilesData, err := getInFileData(inFile)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	// send some data to the publisher channel
	for i := 0; i < *repetitions; i++ {
		for _, v := range inFilesData {
			msgCh <- v.bytes
			time.Sleep(time.Millisecond * time.Duration(*delay))
		}
	}

}
