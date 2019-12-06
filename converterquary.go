package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

/*
 * fungsi ini digunakan untuk menggali data dari bms, data yang diambil
 * diantaranya adalah data tegangan, arus, dan temperature, data ini
 * diakses menggunakan komunikasi modbus
 */

func cnv(port string, address byte, total int, filename string) int {

	// struct yang digunakan untuk membangun file JSON
	type inverter struct {
		Time        string
		Voltage     float32
		Current     float32
		Temperature float32
	}

	// mengisi data dengan nil
	bufstr := inverter{
		Time:        time.Now().Format("15:04:05"),
		Voltage:     0.0,
		Current:     0.0,
		Temperature: 0.0,
	}

	log.Println("quary process from converter begin ......")
	// cek apakah alamat yang dimaksud adalah alamat dari inverter dan converter
	if address != 1 {
		log.Println("this address is not converter address")
		time.Sleep(time.Second)
		return -1
	}

	// membangun koneksi ke inverter melalui modbus RTU jalur serial
	handler := modbus.NewRTUClientHandler(port)

	// setting untuk koneksi serial
	handler.BaudRate = 2400           // baudrate
	handler.DataBits = 8              // databits
	handler.Parity = "N"              // parity
	handler.StopBits = 1              // stopbits
	handler.Timeout = 5 * time.Second // timeout untuk handler

	// alamat yang digunakan storage
	handler.SlaveId = address // alamat untuk modbus slave yang digunakan

	// melakukan kontak dan melakukan error handling saat tidak ada yang merespon
	err := handler.Connect() // connect ke modbus
	if err != nil {
		log.Println(err) // log error ke konsol
		time.Sleep(time.Second)
		return -1 // return -1 sebagai tanda bahwa ini error
	}
	defer handler.Close() // mendereference handler untuk modbus

	// prosedur dari library untuk membangun koneksi dengan modbus
	client := modbus.NewClient(handler)

	// char array penampung untuk data dari modbus
	var data []byte

	ti := time.Now()

	// mengirim permintaan untuk tegangan, arus, dan temperature dari gateway
	data, err = client.ReadHoldingRegisters(0, 3)

	if err != nil {
		log.Println(err) // log error ke konsol
		time.Sleep(time.Second)
		return -1 // return -1 sebagai tanda bahwa ini error
	}

	log.Println(data) // print data yang diterima dari gateway

	var v, c, t float32

	var hi, lo, tot int

	// voltage total - gateway
	hi = int(data[0])
	lo = int(data[1])
	tot = (hi << 8) + lo
	v = float32(tot) // masukan rumus untuk tegangan

	// current total - gateway
	hi = int(data[2])
	lo = int(data[3])
	tot = (hi << 8) | lo
	c = float32(tot) // masukan rumus untuk arus

	// temperature - gateway
	hi = int(data[4])
	lo = int(data[5])
	tot = (hi << 8) | lo
	t = float32(tot) // masukan rumus untuk temperature

	bufstr = inverter{
		Time:        ti.Format("15:04:05"),
		Voltage:     v,
		Current:     c,
		Temperature: t,
	}

	file, err := json.Marshal(bufstr)

	if err != nil {
		log.Println(err) // log error ke konsol
		time.Sleep(time.Second)
		return -1 // return -1 sebagai tanda bahwa ini error
	}

	err = ioutil.WriteFile(filename, file, 06644)

	log.Println("...... quary process from inverter finish")
	return 0
}
