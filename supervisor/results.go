package repot

import (
//	log "github.com/sirupsen/logrus"
)

type jobResult interface {
	Uid() string
	Err() error
}

type SimpleOutput struct {
	uid string
	out []string
	err error
}

func (s *SimpleOutput) Uid() {
	return s.uid
}

func (s *SimpleOutput) Err() {
	return s.err
}

type JobResults struct {
	results map[string][]string
	errors  map[string]error

	resultsChan chan jobResult
	resultsDone chan bool
}

func (j *JobResults) loop() {
	for res := range j.resultsChan {
		j.results[res.uid] = res.out
		j.errors[res.uid] = res.err
	}
	j.resultsDone <- true
}

func (j *JobResults) AddResult(uid string, out interface{}, err error) {
	j.resultsChan <- jobResult{
		uid: uid,
		out: out,
		err: err,
	}
}

func (j *JobResults) GetOut(uid string) []string {
	return j.results[uid]
}

func (j *JobResults) WaitDone() {
	close(j.resultsChan)
	<-j.resultsDone
}

func NewJobsResults() *JobResults {
	j := &JobResults{
		results:     make(map[string][]string),
		errors:      make(map[string]error),
		resultsChan: make(chan jobResult),
		resultsDone: make(chan bool),
	}
	go j.loop()
	return j
}
