package main

import (
	"fmt"
	"sync"
	"time"
)

var ijRecord []ij
var idx int
var recLock sync.Mutex

type ij struct {
	i int
	j int
}

func main() {
	N = 10000
	S = make([]int, N)
	for i := 0; i < N; i++ {
		S[i] = 1
	}
	participants = make([]Participant, N)
	for i := 0; i < N; i++ {
		participants[i] = *new(Participant)
		participants[i].tranIds = *New()
	}

	ijRecord = make([]ij, 3*N)
	idx = 0

	M = 3
	for j := 0; j < M; j++ {
		go worker(S[:], N)
	}
	time.Sleep(12 * time.Second)

	DPrintf(" i j 's record is  ： %+v", ijRecord)
	S1 := make([]int, N)
	for i := 0; i < N; i++ {
		S1[i] = 1
	}
	for _, tempij := range ijRecord {
		S1[tempij.j] = S1[tempij.i] + S1[(tempij.i+1)%N] + S1[(tempij.i+2)%N]
	}

	res := true
	for k := 0; k < N; k++ {
		if S[k] != S1[k] {
			res = false
		}
	}
	DPrintf("Array S is : %+v\n", S)
	DPrintf("Array S1 is : %+v\n", S1)
	fmt.Printf("the result is : %t\n", res)
}
