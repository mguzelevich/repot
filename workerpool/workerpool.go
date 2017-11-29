package workerpool

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type worker struct {
	uid string
}

func (w *worker) loop(input <-chan *job, jobDone chan<- string, workerDone chan string) {
	log.WithFields(log.Fields{"worker": w.uid}).Debug("worker loop started")
	for job := range input {
		log.WithFields(log.Fields{"worker": w.uid, "uid": job.uid}).Debug("job started")
		job.start()
		log.WithFields(log.Fields{"worker": w.uid, "uid": job.uid}).Debug("job finished")
		jobDone <- job.uid
	}
	workerDone <- w.uid
}

type WorkerPool struct {
	workersPool []worker

	jobs  []string
	queue map[string]*job

	addJobChan  chan *job
	addDoneChan chan bool

	jobExecChan chan *job
	jobDoneChan chan string

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
		uid:     uid,
		state:   jobStatePending,
		handler: handler,
	}
	wp.addJobChan <- j
	return nil
}

func (wp *WorkerPool) checkDone() error {
	for {
		select {
		case uid := <-wp.workerDoneChan:
			log.WithFields(log.Fields{"worker": uid}).Debug("worker loop finished")
			allDone := true
			for _, j := range wp.queue {
				allDone = allDone && j.executed()
			}
			if allDone {
				wp.allDoneChan <- true
				return nil
			}
		case uid := <-wp.jobDoneChan:
			log.WithFields(log.Fields{"job": uid}).Debug("job done msg received")
			//wp.queue[uid].executed = true
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
	sv := WorkerPool{
		queue: make(map[string]*job),

		addJobChan:  make(chan *job),
		addDoneChan: make(chan bool),

		jobExecChan: make(chan *job),
		jobDoneChan: make(chan string),

		workerDoneChan: make(chan string),

		allDoneChan: make(chan bool),
	}

	for i := 0; i < poolSize; i++ {
		w := worker{uid: fmt.Sprintf("%d", i)}
		sv.workersPool = append(sv.workersPool, w)
		go w.loop(sv.jobExecChan, sv.jobDoneChan, sv.workerDoneChan)
	}

	go sv.addLoop()

	log.Debug("WorkerPool init done")
	return &sv
}
