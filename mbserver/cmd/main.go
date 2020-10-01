/*
	Notes:
	Single word values(1 x uint16) values can be both uint16 and int16.

	Float (2 x uint16) can have endianess swapped at :
	- byte level within a word.
	- word level where each of the uint16's have swapped place.

	Function codes:
	1, read coils
	2, read descrete inputs
	3, read holding registers
	4, read input registers

	Snip from https://www.csimn.com/CSI_pages/Modbus101.html
	--------------------------------------------------------
	Valid address ranges as originally defined for Modbus were 0 to 9999 for each of the above register types. Valid ranges allowed in the current specification are 0 to 65,535. The address range originally supported by Babel Buster gateways was 0 to 9999. The extended range addressing was later added to all new Babel Buster products.
	The address range applies to each type of register, and one needs to look at the function code in the Modbus message packet to determine what register type is being referenced. The Modicon convention uses the first digit of a register reference to identify the register type.
	Register types and reference ranges recognized with Modicon notation are as follows:
	0x = Coil = 00001-09999
	1x = Discrete Input = 10001-19999
	3x = Input Register = 30001-39999
	4x = Holding Register = 40001-49999
	On occasion, it is necessary to access more than 10,000 of a register type. Based on the original convention, there is another de facto standard that looks very similar. Additional register types and reference ranges recognized with Modicon notation are as follows:
	0x = Coil = 000001-065535
	1x = Discrete Input = 100001-165535
	3x = Input Register = 300001-365535
	4x = Holding Register = 400001-465535
	When using the extended register referencing, it is mandatory that all register references be exactly six digits. This is the only way Babel Buster will know the difference between holding register 40001 and coil 40001. If coil 40001 is the target, it must appear as 040001.

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
	"strings"

	"github.com/RaaLabs/shipsimulator/mbserver"
)

type functionRegisterFiles struct {
	name string
	fh   *os.File
}

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

	// The configuration is split in 4 files, 1 for each register
	fileNames := []string{"coil.json", "descrete.json", "input.json", "holding.json"}
	var fhs []functionRegisterFiles

	// Iterate over all the filenames specified, and create a holding
	// strcuture to keep all the file handles in with info about each
	// register.
	for _, v := range fileNames {
		fh, err := os.Open(v)
		if err != nil {
			log.Printf("error: failed to open config file for %v: %v\n", v, err)
		}

		fName := strings.Split(v, ".")
		f := functionRegisterFiles{
			name: fName[0],
			fh:   fh,
		}
		fhs = append(fhs, f)
		defer fh.Close()
	}
	objs := []map[string]interface{}{}
	js := json.NewDecoder(fhs[3].fh)
	err = js.Decode(&objs)
	//err = json.Unmarshal(js, &objs)
	if err != nil {
		log.Printf("error: decoding json: %v\n", err)
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
