package workerpool

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type poolStateSignal struct {
	uid string
	msg string
}

type WorkerPool struct {
	workersPool []worker

	jobs  []string
	queue map[string]*job

	PoolStateChan chan poolStateSignal
	JobsStateChan chan jobStateSignal

	addJobChan  chan *job
	addDoneChan chan bool

	jobExecChan chan *job

	workerDoneChan chan string
	allDoneChan    chan bool
}

func (wp *WorkerPool) JobState(uid string) jobState {
	return wp.queue[uid].state
}

func (wp *WorkerPool) JobsStatus() []jobState {
	status := []jobState{}
	for _, j := range wp.jobs {
		job := wp.queue[j]
		status = append(status, job.state)
	}
	return status
}

func (wp *WorkerPool) JobsStatusString() string {
	statuses := map[jobState]string{
		jobStatePending:   " ",
		jobStateExecuting: ".",
		jobStateFailed:    "E",
		jobStateFinished:  "+",
	}

	status := ""
	for _, state := range wp.JobsStatus() {
		jobStatus, ok := statuses[state]
		if !ok {
			status += "!"
		}
		status += jobStatus
	}
	return status
}

func (wp *WorkerPool) addLoop() {
	log.Debug("add loop started")
	func() {
		for {
			select {
			case job, ok := <-wp.addJobChan:
				if !ok {
					return
				}
				wp.jobs = append(wp.jobs, job.uid)
				wp.queue[job.uid] = job
				log.WithFields(log.Fields{"uid": job.uid}).Debug("job added")
			default:
				continue
			}
		}
	}()
	log.Debug("all jobs added")
	wp.addDoneChan <- true
}

func (wp *WorkerPool) AddJob(uid string, handler JobHandler) error {
	j := &job{
		uid:       uid,
		state:     jobStatePending,
		stateChan: wp.JobsStateChan,
		handler:   handler,
	}
	wp.addJobChan <- j
	return nil
}

func (wp *WorkerPool) checkDone() error {
	for {
		select {
		//case uid := <-wp.workerDoneChan:
		case <-wp.workerDoneChan:
			allDone := true
			for _, j := range wp.queue {
				allDone = allDone && j.executed()
			}
			if allDone {
				close(wp.PoolStateChan)
				close(wp.JobsStateChan)
				close(wp.allDoneChan)
				return nil
			}
		default:
			continue
		}
	}
}

func (wp *WorkerPool) ExecJobs() error {
	close(wp.addJobChan)
	<-wp.addDoneChan

	go wp.checkDone()

	for _, j := range wp.queue {
		wp.jobExecChan <- j
	}
	log.Debug("all jobs sended")
	close(wp.jobExecChan)

	<-wp.allDoneChan
	return nil
}

func NewWP(poolSize int) *WorkerPool {
	wp := WorkerPool{
		queue: make(map[string]*job),

		PoolStateChan: make(chan poolStateSignal),
		JobsStateChan: make(chan jobStateSignal),

		addJobChan:  make(chan *job),
		addDoneChan: make(chan bool),

		jobExecChan: make(chan *job),

		workerDoneChan: make(chan string),

		allDoneChan: make(chan bool),
	}

	for i := 0; i < poolSize; i++ {
		w := worker{uid: fmt.Sprintf("%d", i)}
		wp.workersPool = append(wp.workersPool, w)
		go w.loop(wp.jobExecChan, wp.workerDoneChan)
	}

	go wp.addLoop()

	log.Debug("WorkerPool init done")
	return &wp
}
