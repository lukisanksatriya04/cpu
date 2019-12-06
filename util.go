package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/tarm/serial"
)

// this used for turn off logging in console
func turnoff() {
	log.SetOutput(ioutil.Discard)
}

// this function used for scanning available serial port in computer
func scanserial() string {
	log.Println("scanning available serial port")

	// linux serial port name
	var buffer [256]string
	for i := 0; i < 256; i++ {
		buffer[i] = "/dev/ttyUSB" + strconv.Itoa(i)
	}

	for i := 0; i < 256; i++ {
		c := &serial.Config{Name: buffer[i], Baud: 9600}
		_, err := serial.OpenPort(c)
		if err == nil {
			log.Println("found serial port in " + buffer[i])
			return buffer[i]
		}
	}

	log.Println("serial port not found")
	return ""
}

// this function used for create a directory
func createdir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}
}
