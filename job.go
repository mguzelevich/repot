package repot

import (
//	log "github.com/sirupsen/logrus"
)

type jobStatus string

const (
	pending   jobStatus = "pending"
	executing jobStatus = "executing"
	failed    jobStatus = "failed"
	finished  jobStatus = "finished"
)

type job struct {
	uid      string
	handler  func(uid string) (string, error)
	executed bool

	status  jobStatus
	results interface{}
}

func (j *job) start() {
	j.status = executing
	if out, err := j.handler(j.uid); err != nil {
		j.fail(out, err)
	} else {
		j.done(out, err)
	}
	j.executed = true
}

func (j *job) done(out string, err error) {
	j.status = finished
}

func (j *job) fail(out string, err error) {
	j.status = failed
}
