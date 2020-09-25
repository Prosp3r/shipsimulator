package modbustcp

import (
	"time"

	modbus "github.com/thinkgos/gomodbus/v2"
)

func Run() {
	srv := modbus.NewTCPServer()
	srv.LogMode(true)
	srv.AddNodes(
		//Holding register
		modbus.NewNodeRegister(1, 0, 10, 0, 10, 0, 10, 0, 10),
		//Input register
		modbus.NewNodeRegister(2, 0, 0, 0, 0, 11, 10, 0, 0),
	)

	l := srv.GetNodeList()

	go func() {
		for {
			l[0].WriteHoldings(0, []uint16{4444, 5555})
			// time.Sleep(time.Second * 2)
			// l[0].WriteHoldings(0, []uint16{2222, 3333})
			// time.Sleep(time.Second * 2)
			l[0].WriteHoldingsBytes(5, 1, []byte{7, 8})

			time.Sleep(time.Second * 2)
		}
	}()

	err := srv.ListenAndServe(":502")
	if err != nil {
		panic(err)
	}
}
