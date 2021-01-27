package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"time"
)

var (
	messages     = "Hello World! Hello KNTU, the best place that exists in the world :|"
	frameNumber  = 0
	t            time.Time
	propagation  int
	bandwidth    int
	fs           int
	flag         = false
	tf           float64
	prompt       = color.New(color.FgGreen)
	serverColor  = color.New(color.FgHiMagenta)
	clientColor  = color.New(color.FgHiYellow)
	ackColor     = color.New(color.FgHiGreen)
	timeoutColor = color.New(color.FgHiRed)
)

func server(data chan []byte, ack chan bool, done chan bool) {
	messageInByte := []byte(messages)
	counter := 0
	var a bool
	buffer := make([]byte, fs)
	for _, b := range messageInByte {
		//Start Timer
		total := time.Now()

		// Make a frame to send
		if counter != fs-1 {
			flag = false
			buffer[counter] = b
			counter++
		} else {
			buffer[counter] = b
			t = time.Now()
			// Print Send Log
			clientSendLog(buffer, t)
			//Simulate Propagation and Transmission Time
			clientSendTimeSimulation()
			//Send Frame
			data <- buffer
			// Listen for Ack
			t = time.Now()
			a = <-ack

			// Implement Timeout mechanism
			failed := 0
			if !a {
				timeoutMechanism(data, ack, buffer, &failed)
			}
			//  Break connection if after 3 time No ack was received
			if failed == 3 {
				break
			}
			_, _ = ackColor.Printf("Ack Time: %s, Total Time: %s\n", time.Now().Sub(t), time.Now().Sub(total))
			frameNumber++
			counter = 0
			flag = true
		}
	}

	if !flag {
		total := time.Now()
		for counter < len(buffer) {
			buffer[counter] = 0
			counter++
		}
		t = time.Now()
		clientSendLog(buffer, t)
		clientSendTimeSimulation()
		data <- buffer
		t = time.Now()
		a = <-ack
		frameNumber++
		failed := 0
		if !a {
			timeoutMechanism(data, ack, buffer, &failed)
		}
		if failed != 3 {
			_, _ = ackColor.Printf("Ack Time: %s, Total Time: %s\n", time.Now().Sub(t), time.Now().Sub(total))
		}
	}
	close(data)
	close(ack)
	done <- true
}

func client(data chan []byte, ack chan bool) {
	for data != nil {
		for b := range data {

			p := rand.Float64()
			if p <= 0.8 {
				_, _ = serverColor.Println("Receiver: ")
				fmt.Println("\tFrame", frameNumber, " Received:", string(b[0:]), time.Now())
				time.Sleep(time.Duration(propagation) * time.Microsecond)
				ack <- true
			} else {
				ack <- false
			}
		}
	}

}

func main() {
	_, _ = prompt.Println("Enter Propagation Time and Bandwidth: ")
	_, err := fmt.Scanf("%d %d", &propagation, &bandwidth)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	_, _ = prompt.Println("Enter Frame Size")
	_, err = fmt.Scanf("%d", &fs)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	tf = float64(fs / bandwidth)
	data := make(chan []byte, fs)
	ack := make(chan bool, 1)
	done := make(chan bool)
	go server(data, ack, done)
	go client(data, ack)
	<-done
}
func clientSendLog(buffer []byte, t time.Time) {
	_, _ = clientColor.Println("Transmitter: ")
	_, _ = clientColor.Printf("\tFrame %d Sent:%s %s \n", frameNumber, buffer, t)
}
func clientSendTimeSimulation() {
	time.Sleep(time.Duration(propagation) * time.Microsecond)
	time.Sleep(time.Duration(tf) * time.Second)
}

func timeoutMechanism(data chan []byte, ack chan bool, buffer []byte, failed *int) {

	fmt.Println("Timeout")
	var a bool
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Millisecond)
		t = time.Now()
		_, _ = clientColor.Println("Transmitter: ")
		fmt.Printf("\tFrame %d Sent:%s %s \n", frameNumber, buffer, time.Now())
		clientSendTimeSimulation()
		data <- buffer
		a = <-ack
		if !a {
			_, _ = timeoutColor.Println("Timeout")
			*failed++
		} else {
			break
		}
	}
}
