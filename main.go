package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	messages    = "The world World! Hello KNTU"
	w           sync.WaitGroup
	flagNumber           = 0
	t           time.Time
	propagation int
	bandwidth   int
	fs          int
	flag=false
	tf float64
)

func server(data chan []byte, ack chan bool, done chan bool) {
	messageInByte := []byte(messages)
	counter := 0
	var a bool
	buffer := make([]byte, fs)
	for _, b := range messageInByte {
		// Make a frame to send
		if counter != fs-1 {
			flag = false
			buffer[counter] = b
			counter++
		} else {
			buffer[counter] = b
			t = time.Now()
			fmt.Println("Transmitter: ")
			fmt.Printf("Frame %d Sent:%s %s \n", flagNumber, buffer, time.Now())
			time.Sleep(time.Duration(propagation) * time.Microsecond)
			time.Sleep(time.Duration(tf)*time.Second)
			data <- buffer
			a = <-ack
			failed :=0
			if !a {
				fmt.Println("Timeout")
				for i:=0 ; i<3 ; i++{
					time.Sleep(1*time.Millisecond)
					t = time.Now()
					fmt.Println("Transmitter: ")
					fmt.Printf("Frame %d Sent:%s %s \n", flagNumber, buffer, time.Now())
					time.Sleep(time.Duration(propagation) * time.Microsecond)
					time.Sleep(time.Duration(tf)*time.Second)
					data <- buffer
					a = <- ack
					if !a {
						fmt.Println("Timeout")
						failed++
					}else{
						break
					}
				}
			}
			if failed == 3 {
				break
			}
			fmt.Println(time.Now().Sub(t))
			flagNumber++
			counter = 0
			flag = true
		}
	}
	if !flag {
		for counter  < len(buffer){
			buffer[counter]=0
			counter++
		}
		t = time.Now()
		fmt.Println("Transmitter: ")
		fmt.Printf("Frame %d Sent:%s %s \n", flagNumber, buffer, time.Now())
		time.Sleep(time.Duration(propagation) * time.Microsecond)
		time.Sleep(time.Duration(tf)*time.Second)
		data <- buffer
		a = <-ack
		fmt.Println(time.Now().Sub(t))
		flagNumber++
	}
	close(data)
	close(ack)
	done <- true
}
func client(data chan []byte, ack chan bool) {
	for data != nil {
		for b := range data {

			p := rand.Float64()
			fmt.Println(p)
			if p <= 0.8 {
				fmt.Println("Receiver: ")
				fmt.Println("Frame", flagNumber, " Received:", string(b[0:]), time.Now())
				time.Sleep(time.Duration(propagation) * time.Microsecond)
				ack <-true
			}else {
				ack <- false
			}
		}
	}

}

func main() {
	fmt.Println("Enter Propagation Time and Bandwidth: ")
	_, err := fmt.Scanf("%d %d", &propagation, &bandwidth)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	fmt.Println("Enter Frame Size")
	_, err = fmt.Scanf("%d", &fs)
	if err != nil {
		log.Fatalln("Fatal Error")
	}
	tf = float64( fs / bandwidth)
	data := make(chan []byte, fs)
	ack := make(chan bool, 1)
	done := make(chan bool)
	go server(data, ack, done)
	go client(data, ack)
	<-done
}
