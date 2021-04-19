package main

import (
	"consensus-algorithm/config"
	"consensus-algorithm/message"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	faulty      bool
	isFailed    = false
	connections = make(map[string]*net.UDPConn) //contains all connections
	R           = make([]float32, config.NodeCount)
	gotValue    = make([]bool, config.NodeCount)
	mu          sync.Mutex
	state       message.Message
	p           = 0
	quit        = make(chan bool)
)

func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("Please provide node id, delay")
		return
	}

	id, _ := strconv.Atoi(arguments[1])
	delay, _ := strconv.ParseInt(arguments[2], 10, 64)

	faulty = (id > 1 && id <= (1+config.FaultCount)) //config.FaultCount nodes after 1 will be faulty

	CONNECT := config.NodesCONNECT[id]
	delete(config.NodesCONNECT, id) // so that node will not message itself

	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(2) // generates state, either 0 or 1
	fmt.Printf("Initial state: %v\n", value)
	state = message.Message{Source: id, Value: float32(value), P: p}
	R[id-1] = state.Value
	gotValue[id-1] = true

	start := time.Now()
	go server(id, CONNECT) //creates a server
	go broadcast(int(delay))

	for {
		time.Sleep(500 * time.Millisecond)
		if p == 100 {
			timeElapsed := time.Since(start)
			defer fmt.Printf("Time taken: %v\n", timeElapsed)
			quit <- true
			time.Sleep(1500 * time.Millisecond) //broadcast for 1.5 seconds before terminating
			return
		}
	}
}

func broadcast(delay int) {
	for {
		if !isFailed {
			mu.Lock()
			s := state
			mu.Unlock()
			time.Sleep(time.Duration(delay) * time.Millisecond) //sleep for [delay] milliseconds
			for _, CONNECT := range config.NodesCONNECT {
				unicast(CONNECT, s) //sends averaged states to all nodes
			}
		}
	}
}

func server(id int, CONNECT string) {
	msg := make(chan message.Message)
	port := ":" + strings.Split(CONNECT, ":")[1]

	//creates server on ip:port CONNECT
	addr, err := net.ResolveUDPAddr("udp4", port)
	ln, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Printf("Error opening server: \n%v\n", err)
		return
	}

	//goroutine handles multiple incoming connections
	go func() {
		buf := make([]byte, 1000)
		rs := message.Message{}
		for {
			length, _, _ := ln.ReadFromUDP(buf)
			buffer := bytes.NewBuffer(buf[:length])
			gob.NewDecoder(buffer).Decode(&rs)
			msg <- rs
		}
	}()

	for {
		select {
		//recieves messages from incoming connection
		case rs := <-msg:
			if rs.P >= state.P && !gotValue[rs.Source-1] {
				R[rs.Source-1] = rs.Value
				gotValue[rs.Source-1] = true
			}
			if totalRecieved() >= (config.NodeCount - config.FaultCount) {
				R[id-1] = reduce(id)
				mu.Lock()
				if faulty {
					faultTest()
				}
				state.P++
				p++
				mu.Unlock()
				reset(id)
			}

		case <-quit:
			return
		}

	}
}

func dial(destination string) {
	addr, err := net.ResolveUDPAddr("udp4", destination)
	c, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return
	}

	connections[destination] = c
}

func unicast(CONNECT string, state message.Message) {
	for {
		if _, ok := connections[CONNECT]; !ok { //if connection doesn't exist, creates it
			dial(CONNECT)
		} else {
			break
		}
	}

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(state)
	connections[CONNECT].Write(buf.Bytes())
}

func faultTest() {
	value := rand.Intn(100) //generates [0, 100)
	isFailed = (value == 1) // 1% chance of failure
}

func totalRecieved() int {
	sum := 0
	for _, v := range gotValue {
		if v {
			sum++
		}
	}
	return sum
}

func reset(id int) {
	gotValue = make([]bool, config.NodeCount)
}

func reduce(id int) float32{
	rSorted := []float32{}
	for i, v := range R {
		if gotValue[i] && (len(rSorted) < config.NodeCount-config.FaultCount) {
			rSorted = append(rSorted, v)
		}
	}

	sort.Slice(rSorted, func(i, j int) bool { return rSorted[i] < rSorted[j] }) //sort non ⊥ elements
	state.Value = (rSorted[config.FaultCount] + rSorted[(config.NodeCount-(2*config.FaultCount))-1]) / 2
	return state.Value
}
