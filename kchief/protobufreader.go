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
//						|- []DataPoints
//							|- Timestamp
//							|- Value

// RunProtoReader will start the process, and read the protobuf binary
// file from disk, unmarshall the Message struct from the
// proto spesification, unzip the Data, then iterate through
// the structure described above to get the TAG, Timestamp
// and value fields.
func RunProtoReader(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error: readData, os.Open failed %v", err)
	}

	// unmarshall the Message struct from the protobuf raw data
	message := messagingpb.Message{}
	err = proto.Unmarshal(b, &message)
	if err != nil {
		return fmt.Errorf("error: failed Unmarshal Message: %v", err)
	}

	// Get the compressed content of the field Data,
	// and uncompress it.
	dataUncompressed, err := getDataFieldFromMessage(message)
	if err != nil {
		return err
	}

	// unmarshal Payloads from Data.
	payloads := messagingpb.Payloads{}
	err = proto.Unmarshal(dataUncompressed, &payloads)
	if err != nil {
		return fmt.Errorf("error: failed Unmarshal Message: %v", err)
	}

	getPayload(payloads)

	return nil
}

// getDataField uncompresses the compress field Data in
// message, and returns the result.
func getDataFieldFromMessage(message messagingpb.Message) ([]byte, error) {

	// The Data field in the Message struct contains
	// compress data that we need to decompress.
	messageDataCompressed := message.GetData()
	dataCompressedReader := bytes.NewReader(messageDataCompressed)

	gzipReader, err := gzip.NewReader(dataCompressedReader)
	if err != nil {
		return nil, fmt.Errorf("error: gzip newreader failed: %v", err)
	}

	dataUncompressed, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		log.Printf("error: readall gzipReader failed: %v\n", err)
	}

	return dataUncompressed, nil
}

// getPayload iterates over all the Payloads,
// then for each respective payload it finds
// the TAG name, and the corresponding TagdataPoint.
func getPayload(payloads messagingpb.Payloads) {
	for _, v := range payloads.Payloads {
		// fmt.Printf("index %v : %#v\n", i, v.GetTagdatapoints())
		vv := v.GetTagdatapoints()
		for _, v := range vv {
			fmt.Printf("*TAG* : type=%T, %v\n", v.Tag, v.Tag)
		}

		getDataPoint(v.GetTagdatapoints())
	}
}

func getDataPoint(tagDataPoints []*messagingpb.TagDataPoints) {
	for _, v := range tagDataPoints {
		dataPoint := v.GetDatapoints()
		for _, v := range dataPoint {
			fmt.Printf("*TIMESTAMP : type=%T, %v\n", v.Timestamp, v.Timestamp)
			fmt.Printf("*VALUE* : type=%T, %v\n", v.Value[0], v.Value[0])

		}
	}
}
