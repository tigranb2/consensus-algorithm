package main

import (
	"consensus-algorithm/config"
	"consensus-algorithm/message"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	faulty bool
	isFailed = false
	connections  = make(map[string]*net.UDPAddr) //contains all connections
	pc			 net.PacketConn
	R            = make([]int, config.NodeCount)
	mu           sync.Mutex
	state        message.Message
	p            = 0
	quit         = make(chan bool)
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

	R[id-1] = 1
	start := time.Now()
	go server(id, CONNECT) //creates a server
	go broadcast(int(delay))

	for {
		if p >= 100 {
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
			variation := rand.Intn((delay/3) + (delay/3)) - (delay/3) //generates number betweeen [-delay/3, delay/3]
			time.Sleep(time.Duration(delay + variation) * time.Millisecond) //sleep for delay + variation milliseconds
			for _, CONNECT := range config.NodesCONNECT {
				unicast(CONNECT, s) //sends averaged states to all nodes
			}
		} 
	}
}

func server(id int, CONNECT string) {
	msg := make(chan message.Message)
	port := ":" + strings.Split(CONNECT, ":")[1]
	pc, _ = net.ListenPacket("udp4", port)

	//goroutine handles multiple incoming connections
	go func() {
		buf := make([]byte, 1000)
		rs := message.Message{}
		for {
			length, _, _ := pc.ReadFrom(buf)
			buffer := bytes.NewBuffer(buf[:length])
			gob.NewDecoder(buffer).Decode(&rs)
			msg <- rs
		}
	}()

	for {
		select {
		//recieves messages from incoming connection
		case rs := <-msg:
			if rs.P > state.P { //copy state and jump to future phase
				mu.Lock()
				if faulty{faultTest()}
				state.Value = rs.Value
				state.P = rs.P
				p = rs.P
				mu.Unlock()
				reset(id)
			} else if rs.P == state.P && R[rs.Source-1] == 0 {
				R[rs.Source-1] = 1
				mu.Lock()
				update(rs.Value) //averages states
				mu.Unlock()
				if sum(R) >= (config.NodeCount - config.FaultCount) {
					mu.Lock()
					if faulty{faultTest()}
					state.P++
					p++
					mu.Unlock()
					reset(id)
				}
			}

		case <-quit:
			return
		}

	}
}

func dial(destination string){
	if pc == nil {
		return
	}

	addr,err := net.ResolveUDPAddr("udp4", destination)
	if err != nil {
		return
	}
	connections[destination] = addr
}

func unicast(CONNECT string, state message.Message) {
	if _, ok := connections[CONNECT]; !ok {
		dial(CONNECT)
		return
	}

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(state)
	pc.WriteTo(buf.Bytes(), connections[CONNECT])
}

func faultTest(){
	value := rand.Intn(100) //generates [0, 100)
	isFailed = (value == 1) // 1% chance of failure
}

func sum(R []int) int {
	total := 0
	for _, i := range R {
		total += i
	}

	return total
}

func reset(i int) {
	R = make([]int, config.NodeCount)
	R[i-1] = 1
}

func update(newState float32) {
	total := (float32(sum(R)-1) * state.Value) + newState
	state.Value = total / float32(sum(R))
}