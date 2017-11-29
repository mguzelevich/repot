package workerpool

import (
	"sync"
	//	log "github.com/sirupsen/logrus"
)

type jobOutput []string

type simpleJobsOutputs struct {
	sync.Mutex
	outputs map[string]jobOutput
}

func (so *simpleJobsOutputs) Add(uid string, result jobOutput) {
	so.Lock()
	so.outputs[uid] = result
	so.Unlock()
}

func (so *simpleJobsOutputs) Get(uid string) jobOutput {
	return so.outputs[uid]
}

func NewSimpleJobsOutputs() *simpleJobsOutputs {
	return &simpleJobsOutputs{outputs: make(map[string]jobOutput)}
}
