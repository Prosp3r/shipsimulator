package kchief

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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
//
// INFO: Each tag that are represented in the protobuf are represented as an
// individual Payload in the the Payloads slice.
//
// The structure of the Payloads
//
// payloads := messagingpb.Payloads{
// 	Payloads: []*messagingpb.Payload{
// 		{
// 			Tagdatapoints: []*messagingpb.TagDataPoints{
// 				{
// 					Tag: "apekatt",
// 					Datapoints: []*messagingpb.DataPoint{
// 						{
// 							Timestamp: ptypes.TimestampNow(),
// 							Value:     []float64{3.14},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }

// RunProtoEncode will loop over the TagDataPoints specified in the JSON file
// given as input, and create TagDataPoints for each of them, put them in a
// Payload struct, and add each one to the Payloads slice.
func RunProtoEncode(fileName string, inJsonFile string) error {
	inDataPoints, err := readJsonInFile(inJsonFile)
	if err != nil {
		log.Printf("error: read inJson failed: %v\n", err)
	}

	// ----------------------marshal payloads---------------------------

	payloadsSlice := []*messagingpb.Payload{}

	for _, v := range inDataPoints.DataPoints {
		tdp := []*messagingpb.TagDataPoints{
			{
				Tag: v.TagName,
				Datapoints: []*messagingpb.DataPoint{
					{
						Timestamp: ptypes.TimestampNow(),
						Value:     []float64{v.Value},
					},
				},
			},
		}
		pl := &messagingpb.Payload{}
		pl.Tagdatapoints = append(pl.Tagdatapoints, tdp...)
		payloadsSlice = append(payloadsSlice, pl)
	}

	payloadsStruct := messagingpb.Payloads{
		Payloads: payloadsSlice,
	}

	// marshall the inner Payloads structure to protobuf format.
	b, err := proto.Marshal(&payloadsStruct)
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

type inDataPoints struct {
	DataPoints []inDataPoint `json:"dataPoints"`
}
type inDataPoint struct {
	TagName string  `json:"tagName"`
	Value   float64 `json:"value"`
}

// readJsonInFile will read the json file,
// unmarshal the data read, and return the structure.
//
// NB: The timestamps have been ommited from the json
// structure, and are generated on the fly while doing
// the protobuf structure later
//
// The struture of the JSON file
//
// {
//     "dataPoints": [{
//             "tagName": "tag nr 1",
//             "value": 1.1
//         },
//         {
//             "tagName": "tag nr 2",
//             "value": 2.2
//         }
//     ]
// }
func readJsonInFile(fileName string) (inDataPoints, error) {
	dps := inDataPoints{}

	fh, err := os.Open(fileName)
	if err != nil {
		return inDataPoints{}, fmt.Errorf("error: failed to open JSON file for reading: %v", err)
	}

	jsonDecoder := json.NewDecoder(fh)
	err = jsonDecoder.Decode(&dps)
	if err != nil {
		return inDataPoints{}, fmt.Errorf("error: decoding json failed: %v", err)
	}

	return dps, nil
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
			log.Printf("info: successfully wrote the data to file=%v :error=%v\n", fileName, err)
			break
		}

		if err != nil {
			log.Printf("error: failed writing file: %v\n", err)
		}
	}

	return nil
}
