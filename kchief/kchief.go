package kchief

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/RaaLabs/shipsimulator/kchief/messagingpb"
	"github.com/golang/protobuf/proto"
)

// The structure seems to be as follows
//
// File, containing protobuf
//	|- Message
//		|- Data, which is in binary format compressed with gzip
//			|- Payloads, which is an array unzipped from Data above
//				|- []Payload
//					|- []TagDataPoints
//						|- TAG, the actual string representation of the tag name
//						|- []DataPoint
//							|- Timestamp
//							|- Value

func readPBMessage() error {
	b, err := ioutil.ReadFile("./sample/sample2.bin")
	if err != nil {
		return fmt.Errorf("error: readData, os.Open failed %v", err)
	}

	// The main protobuf message seems to be in .Message
	message := messagingpb.Message{}
	err = proto.Unmarshal(b, &message)
	if err != nil {
		log.Printf("error: failed Unmarshal Message: %v\n", err)
	}

	// Message for a field called Data which seems to be compress with gzip.
	messageData := message.GetData()
	dataReader := bytes.NewReader(messageData)

	gzipReader, err := gzip.NewReader(dataReader)
	if err != nil {
		log.Printf("error: gzip newreader failed: %v\n", err)
	}

	dataPayloads, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		log.Printf("error: readall gzipReader failed: %v\n", err)
	}

	// fmt.Println(dataPayloads)

	payloads := messagingpb.Payloads{}

	err = proto.Unmarshal(dataPayloads, &payloads)
	if err != nil {
		log.Printf("error: failed Unmarshal Message: %v\n", err)
	}

	getPayloads(payloads)

	return nil
}

func getPayloads(payloads messagingpb.Payloads) {
	for _, v := range payloads.Payloads {
		// fmt.Printf("index %v : %#v\n", i, v.GetTagdatapoints())
		vv := v.GetTagdatapoints()
		for _, v := range vv {
			fmt.Printf("*TAG* : %v\n", v.Tag)
		}

		getDataPoint(v.GetTagdatapoints())
	}
}

func getDataPoint(tagDataPoints []*messagingpb.TagDataPoints) {
	for _, v := range tagDataPoints {
		dataPoint := v.GetDatapoints()
		for _, v := range dataPoint {
			fmt.Printf("*TIMESTAMP : %v\n", v.Timestamp)
			fmt.Printf("*VALUE* : %v\n", v.Value[0])

		}
	}
}

func Run() error {
	err := readPBMessage()
	if err != nil {
		return err
	}
	return nil
}
