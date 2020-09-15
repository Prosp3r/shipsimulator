package kchief

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"

	"github.com/RaaLabs/shipsimulator/kchief/messagingpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
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

func RunProtoEncode(fileName string) error {

	// ----------------------marshal payloads---------------------------

	payloads := messagingpb.Payloads{
		Payloads: []*messagingpb.Payload{
			{
				Tagdatapoints: []*messagingpb.TagDataPoints{
					{
						Tag: "apekatt",
						Datapoints: []*messagingpb.DataPoint{
							{
								Timestamp: ptypes.TimestampNow(),
								Value:     []float64{3.14},
							},
						},
					},
				},
			},
		},
	}

	// marshall the inner Payloads structure to protobuf format.
	b, err := proto.Marshal(&payloads)
	if err != nil {
		log.Printf("error: marshaling proto: %v\n", err)
	}

	fmt.Println(string(b))

	// Create an in-memory buffer to store the gzip'ed data.
	buf := []byte{}
	gzipBuffer := bytes.NewBuffer(buf)

	// write the data b in gzip format into the buffer
	gzipWriter := gzip.NewWriter(gzipBuffer)
	n, err := gzipWriter.Write(b)
	if err != nil {
		log.Printf("error: gzip writer failed: %v\n", err)
	}
	fmt.Printf("info: wrote %v bytes\n", n)

	err = gzipWriter.Flush()
	if err != nil {
		log.Printf("error: gzip write flush failed: %v\n", err)
	}

	err = gzipWriter.Close()
	if err != nil {
		log.Printf("error: gzip writer close failed: %v\n", err)
	}

	// ---------------------------marshal message-------------------------

	// Create a variable containing the upper Message structure
	// of the protobuf structure
	msg := messagingpb.Message{
		Data: gzipBuffer.Bytes(),
	}

	messageProto, err := proto.Marshal(&msg)
	if err != nil {
		log.Printf("error: proto marshal of Message failed: %v\n", err)
	}

	fmt.Printf("message : %v\n", messageProto)

	// Write the protobuf data to file for now...
	err = writeProtoFile(messageProto, fileName)
	if err != nil {
		return fmt.Errorf("error: writeProtofile failed: %v", err)
	}

	return nil
}

// writeProtoFile will write the protobuf data to disk
func writeProtoFile(messageProto []byte, fileName string) error {
	fh, err := os.Create(fileName)
	if err != nil {
		log.Printf("error: os create failed: %v\n", err)
	}
	defer fh.Close()

	for {
		n, err := fh.Write(messageProto)
		fmt.Println(n, err)
		if err == nil {
			log.Printf("info: successfully wrote the data to file: %v\n", err)
			break
		}

		if err != nil {
			log.Printf("error: failed writing file: %v\n", err)
		}
	}

	return nil
}
