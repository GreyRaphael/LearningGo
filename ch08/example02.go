package main

import (
	"context"
	"fmt"
	"time"
)

type Input struct {
	A string
	B string
	C string
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

type processor struct {
	outA chan AOut
	outB chan BOut
	errs chan error
}

func callWebServiceA(ctx context.Context, a string) (AOut, error) {
	// Implement the logic to call web service A
	// time.Sleep(time.Second)
	return AOut{name: "grey", age: 20}, nil
}

func callWebServiceB(ctx context.Context, b string) (BOut, error) {
	// Implement the logic to call web service B
	// time.Sleep(time.Second)
	return BOut{subject: "Physics", score: 96.5}, nil
}

func callWebServiceC(ctx context.Context, c CIn) (COut, error) {
	resultChan := make(chan COut, 1)

	go func() {
		// Implement the logic to call web service C
		joinedData := COut{name: c.A.name, subject: c.B.subject, score: c.B.score}
		// time.Sleep(time.Second) // Replace with the actual logic

		resultChan <- joinedData
	}()

	select {
	case <-ctx.Done():
		return COut{}, ctx.Err() // Context canceled or timed out
	case result := <-resultChan:
		return result, nil
	}
}

func (p *processor) launch(ctx context.Context, data Input) {
	go func() {
		aOut, err := callWebServiceA(ctx, data.A)
		if err != nil {
			p.errs <- err
			return
		}
		p.outA <- aOut
	}()
	go func() {
		bOut, err := callWebServiceB(ctx, data.B)
		if err != nil {
			p.errs <- err
			return
		}
		p.outB <- bOut
	}()
}

func (p *processor) waitForAB(ctx context.Context) (CIn, error) {
	var inputC CIn
	count := 0
	for count < 2 {
		select {
		case a := <-p.outA:
			inputC.A = a
			count++
		case b := <-p.outB:
			inputC.B = b
			count++
		case err := <-p.errs:
			return CIn{}, err
		case <-ctx.Done():
			return CIn{}, ctx.Err()
		}
	}
	return inputC, nil
}

func GatherAndProcess(ctx context.Context, data Input) (COut, error) {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	p := processor{
		outA: make(chan AOut, 1),
		outB: make(chan BOut, 1),
		errs: make(chan error, 2),
	}
	p.launch(ctx, data)
	inputC, err := p.waitForAB(ctx)
	if err != nil {
		return COut{}, err
	}
	out, err := callWebServiceC(ctx, inputC)
	return out, err
}

func main() {
	ctx := context.Background()
	data := Input{} // provide actual input data
	result, err := GatherAndProcess(ctx, data)
	if err != nil {
		// Handle error
		fmt.Println(err)
	} else {
		// Use the result
		fmt.Println(result)
	}
}
