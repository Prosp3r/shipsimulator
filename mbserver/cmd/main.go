/*
	Notes:
	Single word values(1 x uint16) values can be both uint16 and int16.

	Float (2 x uint16) can have endianess swapped at :
	- byte level within a word.
	- word level where each of the uint16's have swapped place.

	References :
	https://control.com/forums/threads/confused-modbus-tcp-vs-modbus-rtu-over-tcp.37377/
	https://www.simplymodbus.ca/TCP.htm
	https://modbus.org/docs/Modbus_Application_Protocol_V1_1b3.pdf
*/

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"

	"github.com/RaaLabs/shipsimulator/mbserver"
)

func main() {
	// Start a new server
	serv := mbserver.NewServer()
	err := serv.ListenRTUTCP(":5502")
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer serv.Close()

	// Create a new registry builder to specific
	// registry data like postions in slice etc.
	const startReg int = 100 // TODO: Replace this with parsed start value from input

	// // Create some input data to fill the register
	// registryData := []encoder{
	// 	float32LittleWordBigEndian{
	// 		Type:   "float32LittleWordBigEndian",
	// 		Number: 3.1415,
	// 		Size:   2,
	// 	},
	// 	float32BigWordBigEndian{
	// 		Type:   "float32BigWordBigEndian",
	// 		Number: 3.141516,
	// 		Size:   2,
	// 	},
	// }

	js := []byte(`[{
		"type":"float32LittleWordBigEndian",
		"number":3.1415,
		"size":2
	}, {
		"type":"float32BigWordBigEndian",
		"number":3.141516,
		"size":2
		}]`)

	objs := []map[string]interface{}{}
	err = json.Unmarshal(js, &objs)
	if err != nil {
		log.Fatal(err)
	}

	var registryData []encoder

	for _, obj := range objs {
		registryData = append(registryData, NewEncoder(obj))
	}

	// setRegister will set the values into the register
	setRegister(serv, registryData, startReg)

	// Wait for someone to press CTRL+C.
	fmt.Println("Press ctrl+c to stop")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Stopped")
}

// uint16ToLittleEndian will swap the byte order of the 'two
// 8 bit bytes that an uint16 is made up of.
func uint16ToLittleEndian(u uint16) uint16 {
	fmt.Printf("before: %b\n", u)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, u)
	v := binary.BigEndian.Uint16(b)
	fmt.Printf("after: %b\n", v)
	return v
}

// setRegister will set the values into the register.
// Will also take the starting position of the register.
func setRegister(serv *mbserver.Server, registryData []encoder, regPos int) error {
	for _, v := range registryData {

		serv.HoldingRegisters = append(serv.HoldingRegisters[:regPos], v.encode()...)

		valueSize := 2
		regPos = regPos + valueSize
	}

	return nil
}

// -----------------------------------Encoder's----------------------------------------

// encoder represent any value type that can be encoded
// into a []uint16 as a response back to the modbus request.
type encoder interface {
	encode() []uint16
	getSize() int
}

type float32LittleWordBigEndian struct {
	Type   string
	Number float64
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	Size float64
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32LittleWordBigEndian) encode() []uint16 {
	n := float32(f.Number)
	v1 := uint16((math.Float32bits(n) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(n)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)
	return []uint16{v2, v1}
}

func (f float32LittleWordBigEndian) getSize() int {
	n := int(f.Size)
	return n
}

// -------

type float32BigWordBigEndian struct {
	Type   string
	Number float64
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	Size float64
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32BigWordBigEndian) encode() []uint16 {
	n := float32(f.Number)
	v1 := uint16((math.Float32bits(n) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(n)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)
	return []uint16{v1, v2}
}

func (f float32BigWordBigEndian) getSize() int {
	n := int(f.Size)
	return n
}

// -------

type float32LittleWordLittleEndian struct {
	Type   string
	Number float64
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	Size float64
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32LittleWordLittleEndian) encode() []uint16 {
	n := float32(f.Number)
	v1 := uint16((math.Float32bits(n) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(n)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)

	v1 = uint16ToLittleEndian(v1)
	v2 = uint16ToLittleEndian(v2)
	return []uint16{v2, v1}
}

func (f float32LittleWordLittleEndian) getSize() int {
	n := int(f.Size)
	return n
}

// -------

type float32BigWordLittleEndian struct {
	Type   string
	Number float64
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	Size float64
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32BigWordLittleEndian) encode() []uint16 {
	n := float32(f.Number)
	v1 := uint16((math.Float32bits(n) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(n)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)

	v1 = uint16ToLittleEndian(v1)
	v2 = uint16ToLittleEndian(v2)
	return []uint16{v2, v1}
}

func (f float32BigWordLittleEndian) getSize() int {
	n := int(f.Size)
	return n
}

// -------

type wordInt16BigEndian struct {
	Type   string
	Number float64
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	Size float64
}

func (w wordInt16BigEndian) encode() []uint16 {
	v := uint16(w.Number)

	return []uint16{v}
}

func (f wordInt16BigEndian) getSize() int {
	return int(f.Size)
}

// -------

type wordInt16LittleEndian struct {
	Type   string
	Number float64
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	Size float64
}

func (w wordInt16LittleEndian) encode() []uint16 {
	v := uint16(w.Number)
	v = uint16ToLittleEndian(v)

	return []uint16{v}
}

func (f wordInt16LittleEndian) getSize() int {
	return int(f.Size)
}

// -------------------------------------------------------------------------

func NewEncoder(m map[string]interface{}) encoder {
	switch m["type"].(string) {
	case "float32LittleWordBigEndian":
		return NewFloat32LittleWordBigEndian(m)
	case "float32BigWordBigEndian":
		return NewFloat32BigWordBigEndian(m)
	case "float32LittleWordLittleEndian":
		return NewFloat32LittleWordLittleEndian(m)
	case "float32BigWordLittleEndian":
		return NewFloat32BigWordLittleEndian(m)
	case "wordInt16BigEndian":
		return NewWordInt16BigEndian(m)
	case "wordInt16LittleEndian":
		return NewWordInt16LittleEndian(m)
	}
	return nil
}

func NewFloat32LittleWordBigEndian(m map[string]interface{}) *float32LittleWordBigEndian {
	return &float32LittleWordBigEndian{
		Type:   m["type"].(string),
		Number: m["number"].(float64),
		Size:   m["size"].(float64),
	}
}

func NewFloat32BigWordBigEndian(m map[string]interface{}) *float32BigWordBigEndian {
	return &float32BigWordBigEndian{
		Type:   m["type"].(string),
		Number: m["number"].(float64),
		Size:   m["size"].(float64),
	}
}

func NewFloat32LittleWordLittleEndian(m map[string]interface{}) *float32LittleWordLittleEndian {
	return &float32LittleWordLittleEndian{
		Type:   m["type"].(string),
		Number: m["number"].(float64),
		Size:   m["size"].(float64),
	}
}

func NewFloat32BigWordLittleEndian(m map[string]interface{}) *float32BigWordLittleEndian {
	return &float32BigWordLittleEndian{
		Type:   m["type"].(string),
		Number: m["number"].(float64),
		Size:   m["size"].(float64),
	}
}

func NewWordInt16BigEndian(m map[string]interface{}) *wordInt16BigEndian {
	return &wordInt16BigEndian{
		Type:   m["type"].(string),
		Number: m["number"].(float64),
		Size:   m["size"].(float64),
	}
}

func NewWordInt16LittleEndian(m map[string]interface{}) *wordInt16LittleEndian {
	return &wordInt16LittleEndian{
		Type:   m["type"].(string),
		Number: m["number"].(float64),
		Size:   m["size"].(float64),
	}
}
