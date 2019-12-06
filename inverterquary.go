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

func inv(port string, address byte, total int, filename string) int {

	// struct yangS digunakan untuk membangun file JSON
	type inverter struct {
		Time        string
		Voltage     float32
		Current     float32
		Power       float32
		Powerfactor float32
		Frequency   float32
		Temperature float32
	}

	// mengisi data dengan nil
	bufstr := inverter{
		Time:        time.Now().Format("15:04:05"),
		Voltage:     0.0,
		Current:     0.0,
		Power:       0.0,
		Powerfactor: 0.0,
		Frequency:   0.0,
		Temperature: 0.0,
	}

	log.Println("quary process from inverter begin ......")
	// cek apakah alamat yang dimaksud adalah alamat dari inverter dan converter
	if address != 2 {
		log.Println("this address is not inverter address")
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
	data, err = client.ReadHoldingRegisters(0, 10)

	if err != nil {
		log.Println(err) // log error ke konsol
		time.Sleep(time.Second)
		return -1 // return -1 sebagai tanda bahwa ini error
	}

	log.Println(data) // print data yang diterima dari gateway

	var v, c, p, pf, f, t float32

	var hi, lo, tot, hi2, lo2 int

	// voltage total - gateway
	hi = int(data[0])
	lo = int(data[1])
	tot = (hi << 8) + lo
	v = float32(tot) // masukan rumus untuk tegangan

	// current total - gateway
	hi = int(data[2])
	lo = int(data[3])
	hi2 = int(data[4])
	lo2 = int(data[5])
	tot = (hi << 24) | (lo << 16) | (hi2 << 8) | lo2
	c = float32(tot) // masukan rumus untuk arus

	// power total - gateway
	hi = int(data[6])
	lo = int(data[7])
	hi2 = int(data[8])
	lo2 = int(data[9])
	tot = (hi << 24) | (lo << 16) | (hi2 << 8) | lo2
	p = float32(tot) // masukan rumus untuk power

	// frequency total - gateway
	hi = int(data[14])
	lo = int(data[15])
	tot = (hi << 8) | lo
	f = float32(tot) // masukan rumus untuk frequency

	// power factor total - gateway
	hi = int(data[16])
	lo = int(data[17])
	tot = (hi << 8) | lo
	pf = float32(tot) // masukan rumus untuk power factor

	// temperature - gateway
	hi = int(data[18])
	lo = int(data[19])
	tot = (hi << 8) | lo
	t = float32(tot) // masukan rumus untuk temperature

	bufstr = inverter{
		Time:        ti.Format("15:04:05"),
		Voltage:     v,
		Current:     c,
		Power:       p,
		Powerfactor: pf,
		Frequency:   f,
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
