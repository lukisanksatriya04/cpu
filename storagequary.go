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

func str(port string, address byte, id int, filename string) int {
	// Ini struct untuk JSON
	// Storage : struct yang digunakan untuk membangun file JSON
	type storage struct {
		ID              int
		Time            string
		Voltagepack     float32     // tegangan total dari battery
		Voltagecell     [13]float32 // tegangan per sel
		Current         float32     // arus yang keluar dari storage
		Capacity        float32     // kapasitas per storage
		Cycle           float32     // cycle umut battery
		Temperaturepack float32     // temperature dari gateway
		Temperaturecell [13]float32 // temperature per sel
	}
	// mengisi data dengan nil
	bufstr := storage{
		ID:              id,
		Time:            time.Now().Format("15:04:05"),
		Voltagepack:     0,
		Voltagecell:     [13]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Current:         0,
		Capacity:        0,
		Cycle:           0,
		Temperaturepack: 0,
		Temperaturecell: [13]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	log.Println("quary process from storage begin ......")
	// cek apakah alamat yang dimaksud adalah alamat dari inverter dan converter
	if address < 3 {
		log.Println("this address has been used by inverter and converter, please use another address")
		time.Sleep(time.Second)
		return -1
	}

	// membangun koneksi ke storage melalui modbus RTU jalur serial
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

	t := time.Now()

	// mengirim permintaan untuk tegangan, arus, dan temperature dari gateway
	data, err = client.ReadHoldingRegisters(2, 3)
	if err != nil {
		log.Println(err) // log error ke konsol
		time.Sleep(time.Second)
		return -1 // return -1 sebagai tanda bahwa ini error
	}
	log.Println(data) // print data yang diterima dari gateway

	var voltagepackf, currentf, temperaturepackf float32
	var voltagecellf, temperaturecellf [13]float32

	var hi, lo, tot int

	// voltage total - gateway
	hi = int(data[0])
	lo = int(data[1])
	tot = (hi << 8) | lo
	voltagepackf = float32(tot) // masukan rumus untuk tegangan
	log.Printf("v : %.2f \n", voltagepackf)

	// current total - gateway
	hi = int(data[2])
	lo = int(data[3])
	tot = (hi << 8) | lo
	currentf = float32(tot / 10) // masukan rumus untuk arus
	log.Printf("c : %.2f \n", currentf)

	hi = int(data[4])
	lo = int(data[5])
	tot = (hi << 8) | lo
	temperaturepackf = float32(tot / 10) // masukan rumus untuk temperature
	log.Printf("t : %.2f \n", temperaturepackf)

	var i uint16

	for i = 0; i < 13; i++ {
		data, err = client.ReadHoldingRegisters(i*2+5, 2)
		if err != nil {
			log.Println(err) // log error ke konsol
			time.Sleep(time.Second)
			return -1 // return -1 sebagai tanda bahwa ini error
		}
		log.Println(data) // print data yang diterima dari gateway

		hi = int(data[0])
		lo = int(data[1])
		tot = (hi << 8) | lo
		voltagecellf[i] = float32(tot / 10) // masukan rumus untuk voltage per cell

		hi = int(data[2])
		lo = int(data[3])
		tot = (hi << 8) | lo
		temperaturecellf[i] = float32(tot / 10) // masukan rumus untuk temperature per cell

		log.Printf("t : %.2f  v : %.2f \n", temperaturecellf[i], voltagecellf[i])
	}

	bufstr = storage{
		ID:              id,
		Time:            t.Format("15:04:05"),
		Voltagepack:     voltagepackf,
		Voltagecell:     voltagecellf,
		Current:         currentf,
		Capacity:        capacity(),
		Cycle:           cycle(),
		Temperaturepack: temperaturepackf,
		Temperaturecell: temperaturecellf,
	}

	file, err := json.Marshal(bufstr)

	if err != nil {
		log.Println(err) // log error ke konsol
		time.Sleep(time.Second)
		return -1 // return -1 sebagai tanda bahwa ini error
	}

	err = ioutil.WriteFile(filename, file, 06644)

	log.Println("...... quary process from storage finish")
	return 0
}

func cycle() float32 {
	return 0
}

func capacity() float32 {
	return 0
}