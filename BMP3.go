package main

import (
	"consensus-algorithm/config"
	"consensus-algorithm/message"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
	"sort"
	"strconv"
)

var (
	connections = make(map[string]*gob.Encoder) //contains all connections
	faulty      bool
	isFailed 	= false
	lossRate    int
	r           = 0
	start       time.Time
	timeout		int
)

//sends and recieves state
func main() {
	arguments := os.Args

	if len(arguments) < 4 {
		fmt.Println("Please provide node id, loss rate, and timeout")
		return
	}

	id, _ := strconv.Atoi(arguments[1])
	timeout, _ = strconv.Atoi(arguments[2])
	lossRate, _ = strconv.Atoi(arguments[3])

	faulty = (id > 1 && id <= (1+config.FaultCount)) //config.FaultCount nodes after 1 will be faulty

	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(2) // generates state, either 0 or 1
	fmt.Printf("Initial state: %v\n", value)
	initialState := message.Message{Source: id, Value: float32(value), R: r}

	start = time.Now()
	server(id, initialState) //creates a server
}

func server(id int, initialState message.Message) {
	localStates := make(map[int][]message.Message) //key is the source, value is slice of all states recieved
	conn := make(chan net.Conn)
	msg := make(chan message.Message)
	state := message.Message{}
	CONNECT := config.NodesCONNECT[id]

	//creates server on ip:port CONNECT
	ln, err := net.Listen("tcp", CONNECT)
	if err != nil {
		log.Printf("Error opening server: \n%v\n", err)
	}

	//goroutine handles multiple incoming connections
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}
			conn <- c
		}
	}()

	broadcast(initialState)

	for {
		select {
		//recieves messages from incoming connection
		case c := <-conn:
			go func(c net.Conn) {
				decoder := gob.NewDecoder(c)
				rs := message.Message{}
				for {
					decoder.Decode(&rs)
					msg <- rs
				}
			}(c)

		case rs := <-msg:
			localStates[rs.Source] = append(localStates[rs.Source], rs)
		}

		currentStates := make(map[int]float32) //stores all states from current round
		for source := 1; source <= config.NodeCount; source++ { //localStates[souce] is each individual node's slice of states
			for i := 0; i < len(localStates[source]); i++ { //localStates[source][i] is one of state received from source
				if localStates[source][i].R == r {
					currentStates[localStates[source][i].Source] = localStates[source][i].Value
				} else if localStates[source][i].R < r {
					localStates[source] = append(localStates[source][:i], localStates[source][i+1:]...) //deletes old value
				}
			}
		}

		if len(currentStates) >= (config.NodeCount - config.FaultCount) { //checks if sufficent messages have been recieved
			newValue := reduce(currentStates)
			r++
			state = message.Message{Source: id, Value: newValue, R: r}

			isFailed = faulty && test(1) //1% chance of failure if node if faulty
			if !isFailed {
				broadcast(state)
			}

			if r == 100 { //on round 100, broadcast for 1.5 seconds and terminate
				timeElapsed := time.Since(start)
				defer fmt.Printf("Time taken: %v\n", timeElapsed)
				go func() {
					for {
						time.Sleep(100 * time.Millisecond) //send message every 100ms
						broadcast(state)
					}
				}()
				time.Sleep(1500 * time.Millisecond)
				return
			}
		}
	}
}

func reduce(currentStates map[int]float32) float32 {
	sortedSlice := []float32{}
	for _, state := range currentStates {
		sortedSlice = append(sortedSlice, state)
	}
	
	sort.Slice(sortedSlice, func(i, j int) bool {return sortedSlice[i] < sortedSlice[j]})
	return (sortedSlice[config.FaultCount] + sortedSlice[config.NodeCount - (2 * config.FaultCount)-1])/2.0
}

func test(p int) bool {
	value := rand.Intn(100)
	if value < p { // p% chance of returning true
		return true
	}

	return false
}

func dial(destination string) {
	c, err := net.Dial("tcp", destination) //dial specified TCP server
	if err != nil {
		return //returns if server could not be established
	}
	connections[destination] = gob.NewEncoder(c)
}

func broadcast(state message.Message) {
	for _, CONNECT := range config.NodesCONNECT {
		for _, ok := connections[CONNECT]; ok == false; _, ok = connections[CONNECT] { //retries dial while connection doesn't exist
			dial(CONNECT)
		}

		for {
			if !test(lossRate) { //message isn't considered "lost" (lostRate% chance to be lost)
				fmt.Printf("Sending %v to %v\n", state, CONNECT)
				connections[CONNECT].Encode(state)
				break
			} else {
				time.Sleep(time.Duration(timeout) * time.Millisecond)
			}
		}
	}
}
