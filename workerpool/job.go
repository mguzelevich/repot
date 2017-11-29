package workerpool

import (
//	log "github.com/sirupsen/logrus"
)

type jobState string

const (
	jobStatePending   jobState = "pending"
	jobStateExecuting jobState = "executing"
	jobStateFailed    jobState = "failed"
	jobStateFinished  jobState = "finished"
)

type JobHandler func(uid string) error

type job struct {
	uid     string
	state   jobState
	handler JobHandler
}

func (j *job) start() {
	j.state = jobStateExecuting
	if err := j.handler(j.uid); err != nil {
		j.fail(err)
	} else {
		j.done(err)
	}
}

func (j *job) executed() bool {
	return j.state == jobStateFinished || j.state == jobStateFailed
}

func (j *job) done(err error) {
	j.state = jobStateFinished
}

func (j *job) fail(err error) {
	j.state = jobStateFailed
}
