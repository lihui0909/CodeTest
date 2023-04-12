package main

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"sync"
	"time"
)

//事务协调者,维护每个worker的全局事务状态
type TranCoordinator struct {
	id         string
	i0, i1, i2 int
	j          int
}

//begin 开启事务
func (tc *TranCoordinator) begin() string {
	//生成事务id
	rand, _ := rand.Int(rand.Reader, big.NewInt(int64(10000)))
	timeStamp := strconv.Itoa((int)(time.Now().Unix()))
	tc.id = rand.String() + timeStamp + strconv.Itoa(tc.i0) + strconv.Itoa(tc.i1) + strconv.Itoa(tc.i2) + strconv.Itoa(tc.j)
	DPrintf("begin transaction %s!\n", tc.id)
	return tc.id
}

func (tc *TranCoordinator) Prepare(p0 Participant, p1 Participant, p2 Participant, pj Participant) bool {
	res0 := p0.readLock(tc.id)
	res1 := p1.readLock(tc.id)
	res2 := p2.readLock(tc.id)
	resj := pj.writeLock(tc.id)
	if res0 && res1 && res2 && resj {
		return true
	}
	return false
}

func (tc *TranCoordinator) Close(p0 Participant, p1 Participant, p2 Participant, pj Participant) {
	DPrintf("close transaction %s!\n", tc.id)
	p0.releaseRead(tc.id)
	p1.releaseRead(tc.id)
	p2.releaseRead(tc.id)
	pj.releaseWrite(tc.id)
}

type Participant struct {
	state   int //参与者状态：0没有锁 1加了读锁 2加了写锁
	tranIds Set //如果有锁，记录加锁的事务id
	m       sync.RWMutex
}

func (p *Participant) readLock(id string) bool {
	p.m.Lock()
	defer p.m.Unlock()
	state := p.state
	set := p.tranIds
	if state == 0 || state == 1 || set.Contains(id) {
		p.state = 1
		p.tranIds.Add(id)
		return true
	}
	return false
}

func (p *Participant) writeLock(id string) bool {
	p.m.Lock()
	defer p.m.Unlock()
	state := p.state
	set := p.tranIds
	if state == 0 || (set.Size() == 1 && set.Contains(id)) {
		p.state = 2
		p.tranIds.Add(id)
		return true
	}
	return false
}

func (p *Participant) releaseRead(id string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.tranIds.Remove(id)
	if p.state == 1 && p.tranIds.Size() == 0 {
		p.state = 0
	}
}

func (p *Participant) releaseWrite(id string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.tranIds.Remove(id)
	p.state = 0
}
