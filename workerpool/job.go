package workerpool

import (
	log "github.com/sirupsen/logrus"
)

type jobState string

const (
	jobStatePending   jobState = "pending"
	jobStateExecuting jobState = "executing"
	jobStateFailed    jobState = "failed"
	jobStateFinished  jobState = "finished"
)

type JobHandler func(uid string) error

type jobStateSignal struct {
	uid   string
	state jobState
}

type job struct {
	uid     string
	state   jobState
	handler JobHandler

	workerUid string
	stateChan chan jobStateSignal
}

func (j *job) stateChanged(state jobState) {
	j.state = state
	log.WithFields(log.Fields{"job": j.uid, "state": state, "wuid": j.workerUid}).Debug("job state changed")
	select {
	case j.stateChan <- jobStateSignal{uid: j.uid, state: state}:
		// fmt.Println("sent message", msg)
	default:
		//fmt.Println("no message sent")
	}
}

func (j *job) start(workerUid string) {
	j.workerUid = workerUid
	j.stateChanged(jobStateExecuting)
	if err := j.handler(j.uid); err != nil {
		j.stateChanged(jobStateFailed)
		// j.fail(err)
	} else {
		j.stateChanged(jobStateFinished)
		// j.done(err)
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

func newJob(uid string, stateChan chan jobStateSignal, handler JobHandler) *job {
	j := &job{
		uid:     uid,
		handler: handler,

		workerUid: "<supervisor>",
		stateChan: stateChan,
	}
	j.stateChanged(jobStatePending)
	return j
}
