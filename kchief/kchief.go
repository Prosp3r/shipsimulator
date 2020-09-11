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

func readPBMessage() error {
	b, err := ioutil.ReadFile("./sample/sample2.bin")
	if err != nil {
		return fmt.Errorf("error: readData, os.Open failed %v", err)
	}

	// The main protobuf message seems to be in .Message
	m := messagingpb.Message{}
	err = proto.Unmarshal(b, &m)
	if err != nil {
		log.Printf("error: failed Unmarshal Message: %v\n", err)
	}

	// fmt.Printf("%#v", m)

	// Message for a field called Data which seems to be compress with gzip.
	mData := m.GetData()
	dataReader := bytes.NewReader(mData)

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

	tagDataPoints := []*messagingpb.TagDataPoints{}
	for i, v := range payloads.Payloads {
		fmt.Printf("index %v : %#v\n", i, v.GetTagdatapoints())
		vv := v.GetTagdatapoints()
		for i, v := range vv {
			fmt.Printf("inner index %v : %v\n", i, v.Tag)

		}
		tagDataPoints = append(tagDataPoints, v.GetTagdatapoints()...)
	}

	for _, v := range tagDataPoints {
		dataPoint := v.GetDatapoints()
		for i, v := range dataPoint {
			fmt.Printf("index %v : %#v\n", i, v.Timestamp)
			fmt.Printf("index %v : %v\n", i, v.Value[0])

		}
	}

	return nil
}

func Run() error {
	err := readPBMessage()
	if err != nil {
		return err
	}
	return nil
}
