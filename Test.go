package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

var mutexs []sync.RWMutex

func main() {
	N := 10000
	S := make([]int, N)
	for i := 0; i < N; i++ {
		S[i] = 1
	}
	mutexs = make([]sync.RWMutex, N)
	M := 3
	for j := 0; j < M; j++ {
		go worker(S[:], N)
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("%+v", S)
}

// worker 中重复10000次操作1.随机生成i,j([0，10000）范围），2.更新S使得S(j) = S(i)+S(i+1)+S(i+2) （如果i+1,i+2大于N，取余）
func worker(S []int, N int) {
	for k := 0; k < N; k++ {
		// 1.随机生成i,j
		ran1, _ := rand.Int(rand.Reader, big.NewInt(int64(N)))
		i := (int)(ran1.Int64())
		ran2, _ := rand.Int(rand.Reader, big.NewInt(int64(N)))
		j := (int)(ran2.Int64())
		fmt.Printf("i = %d, j = %d\n", i, j)
		//2.更新S
		//2-1: 2PL第一阶段——扩展阶段
		mutexs[i].RLock()
		mutexs[(i+1)%N].RLock()
		mutexs[(i+2)%N].RLock()

		tmp := S[i] + S[(i+1)%N] + S[(i+2)%N]

		//2-2: 2PL第二阶段——收缩阶段
		mutexs[i].RUnlock()
		mutexs[(i+1)%N].RUnlock()
		mutexs[(i+2)%N].RUnlock()

		// 读写分开，避免j落在[i,i+2]中的情况单线程内出现死锁。同时也会避免出现线程之间的死锁，因为读锁是按顺序分配的，在读取数字个数小于N的情况下不会出现循环等待的情况
		mutexs[j].Lock()
		S[j] = tmp
		mutexs[j].Unlock()
	}
}
