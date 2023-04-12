package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var participants []Participant
var S []int
var N int
var M int

func worker(S []int, N int) {
	for k := 0; k < N; k++ {
		// 1.随机生成i,j
		ran1, _ := rand.Int(rand.Reader, big.NewInt(int64(N)))
		i := (int)(ran1.Int64())
		ran2, _ := rand.Int(rand.Reader, big.NewInt(int64(N)))
		j := (int)(ran2.Int64())
		fmt.Printf("i = %d, j = %d\n", i, j)
		execSuccess := false
		for !execSuccess {
			//2.开启事务
			tran := &TranCoordinator{}
			tran.i0 = i
			tran.i1 = i + 1
			tran.i2 = i + 2
			tran.begin()

			//3.事务开始第一阶段prepare阶段
			execSuccess = tran.Prepare(participants[i], participants[(i+1)%N], participants[(i+2)%N], participants[j])
			//4.事务开始第二阶段，根据prepare返回结果决定是提交还是回退
			if execSuccess {
				//提交阶段，直接执行
				tmpSum := S[i] + S[(i+1)%N] + S[(i+2)%N]
				S[j] = tmpSum
				tempij := new(ij)
				tempij.i = i
				tempij.j = j
				recLock.Lock()
				ijRecord[idx] = *tempij
				idx += 1
				recLock.Unlock()
			} else {
				//如果不成功，回退不做操作
			}
			tran.Close(participants[i], participants[(i+1)%N], participants[(i+2)%N], participants[j])
		}

	}
}
