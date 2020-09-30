/*
	Notes:
	Single word values(1 x uint16) values can be both uint16 and int16.

	Float (2 x uint16) can have endianess swapped at :
	- byte level within a word.
	- word level where each of the uint16's have swapped place.
*/

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"

	"github.com/RaaLabs/shipsimulator/mbserver"
)

func main() {
	serv := mbserver.NewServer()
	err := serv.ListenTCP(":5502")
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer serv.Close()

	// Create a new registry builder to specific
	// registry data like postions in slice etc.
	const startReg int = 0 // TODO: Replace this with parsed start value from input

	registryData := []encoder{
		float32LittleWordBigEndian{
			number: 3.1415,
			size:   2,
		},
		float32BigWordBigEndian{
			number: 3.141516,
			size:   2,
		},
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
	number float32
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	size int
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32LittleWordBigEndian) encode() []uint16 {
	v1 := uint16((math.Float32bits(f.number) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(f.number)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)
	return []uint16{v2, v1}
}

func (f float32LittleWordBigEndian) getSize() int {
	return f.size
}

// -------

type float32BigWordBigEndian struct {
	number float32
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	size int
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32BigWordBigEndian) encode() []uint16 {
	v1 := uint16((math.Float32bits(f.number) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(f.number)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)
	return []uint16{v1, v2}
}

func (f float32BigWordBigEndian) getSize() int {
	return f.size
}

// -------

type float32LittleWordLittleEndian struct {
	number float32
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	size int
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32LittleWordLittleEndian) encode() []uint16 {
	v1 := uint16((math.Float32bits(f.number) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(f.number)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)

	v1 = uint16ToLittleEndian(v1)
	v2 = uint16ToLittleEndian(v2)
	return []uint16{v2, v1}
}

func (f float32LittleWordLittleEndian) getSize() int {
	return f.size
}

// -------

type float32BigWordLittleEndian struct {
	number float32
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	size int
}

// encode will encode a float32 value into []uint16 where:
//	- The two 16 bits word are little endian
//	- The Byte order of each word a big endian
func (f float32BigWordLittleEndian) encode() []uint16 {
	v1 := uint16((math.Float32bits(f.number) >> 16) & 0xffff)
	v2 := uint16((math.Float32bits(f.number)) & 0xffff)
	fmt.Printf("*v1 = %v*\n", v1)

	v1 = uint16ToLittleEndian(v1)
	v2 = uint16ToLittleEndian(v2)
	return []uint16{v2, v1}
}

func (f float32BigWordLittleEndian) getSize() int {
	return f.size
}

// -------

type wordInt16BigEndian struct {
	number int16
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	size int
}

func (w wordInt16BigEndian) encode() []uint16 {
	v := uint16(w.number)

	return []uint16{v}
}

func (f wordInt16BigEndian) getSize() int {
	return f.size
}

// -------

type wordInt16LittleEndian struct {
	number int16
	// size in the measure of uint16's.
	// E.g. a float32 contains 2 x uint16's,
	// so the size will be 2.
	size int
}

func (w wordInt16LittleEndian) encode() []uint16 {
	v := uint16(w.number)
	v = uint16ToLittleEndian(v)

	return []uint16{v}
}

func (f wordInt16LittleEndian) getSize() int {
	return f.size
}

// -------------------------------------------------------------------------
