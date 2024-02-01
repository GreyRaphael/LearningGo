package main

import (
	"errors"
	"fmt"
	"time"
)

type Input struct {
	A string
	B string
}

type AOut struct {
	name string
	age  int
}

type BOut struct {
	subject string
	score   float64
}

type CIn struct {
	A AOut
	B BOut
}

type COut struct {
	name    string
	subject string
	score   float64
}

func callWebServiceA(a string, out chan<- AOut) {
	// Implement the logic to call web service A
	// time.Sleep(time.Second)
	out <- AOut{name: "grey", age: 20}
}

func callWebServiceB(b string, out chan<- BOut) {
	// Implement the logic to call web service B
	// time.Sleep(time.Second)
	out <- BOut{subject: "Physics", score: 96.5}
}

func callWebServiceC(in <-chan CIn, out chan<- COut) {
	// Implement the logic to call web service C
	inputC := <-in
	// time.Sleep(time.Second)
	out <- COut{inputC.A.name, inputC.B.subject, inputC.B.score}
}

func waitForAB(outA <-chan AOut, outB <-chan BOut) (CIn, error) {
	var inputC CIn
	count := 0
	for count < 2 {
		select {
		case a := <-outA:
			inputC.A = a
			count++
		case b := <-outB:
			inputC.B = b
			count++
		case <-time.After(50 * time.Millisecond):
			return CIn{}, errors.New("timeout exceed in waitForAB")
		}
	}
	return inputC, nil
}

func waitForC(outC <-chan COut) (COut, error) {
	select {
	case c := <-outC:
		return c, nil
	case <-time.After(50 * time.Millisecond):
		return COut{}, errors.New("timeout exceed in waitForC")
	}
}

func gatherAndProcess(data Input) (COut, error) {
	outA := make(chan AOut, 1)
	outB := make(chan BOut, 1)
	inC := make(chan CIn, 1)
	outC := make(chan COut, 1)

	go callWebServiceA(data.A, outA)
	go callWebServiceB(data.B, outB)

	inputC, err := waitForAB(outA, outB)
	if err != nil {
		return COut{}, err
	}

	inC <- inputC
	go callWebServiceC(inC, outC)

	result, err := waitForC(outC)
	return result, err
}

func main() {
	data := Input{}
	result, err := gatherAndProcess(data)
	if err != nil {
		// Handle error
		fmt.Println(err)
	} else {
		// Use the result
		fmt.Println(result)
	}
}
