package repot

import (
	"fmt"
	"os"
	"time"

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

type superVisor struct {
	ShowProgress bool

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

func (s *superVisor) status() string {
	status := ""
	for _, j := range s.jobs {
		jobStatus := "."
		switch job := s.queue[j]; job.status {
		case pending:
			jobStatus = "."
		case executing:
			jobStatus = "*"
		case failed:
			jobStatus = "E"
		case finished:
			jobStatus = "+"
		default:
			jobStatus = "!"
		}
		status += jobStatus
	}
	return status
}

func (s *superVisor) logStatus() {
	if !s.ShowProgress {
		return
	}
	status := s.status()
	fmt.Fprintf(os.Stderr, "JQ: %s\n", status)
	// log.WithFields(log.Fields{"status": status}).Info("JQ")
}

func (s *superVisor) addLoop() {
	log.Debug("add loop started")
	func() {
		for {
			select {
			case job, ok := <-s.addJobChan:
				if !ok {
					return
				}
				s.jobs = append(s.jobs, job.uid)
				s.queue[job.uid] = job
				log.WithFields(log.Fields{"uid": job.uid}).Debug("job added")
			default:
				continue
			}
		}
	}()
	log.Debug("all jobs added")
	s.addDoneChan <- true
}

func (s *superVisor) statusLoop() {
	log.Debug("status loop started")
	heartbeat := time.Tick(2 * time.Second)
	for {
		select {
		case <-heartbeat:
			s.logStatus()
		}
	}
}

func (s *superVisor) AddJob(uid string, handler func(uid string) (string, error)) error {
	j := &job{
		uid:     uid,
		handler: handler,
		status:  pending,
	}
	s.addJobChan <- j
	return nil
}

func (s *superVisor) checkDone() error {
	for {
		select {
		case uid := <-s.workerDoneChan:
			log.WithFields(log.Fields{"worker": uid}).Debug("worker loop finished")
			allDone := true
			for _, j := range s.queue {
				allDone = allDone && j.executed
			}
			if allDone {
				s.allDoneChan <- true
				return nil
			}
		case uid := <-s.jobDoneChan:
			log.WithFields(log.Fields{"job": uid}).Debug("job done msg received")
			//s.queue[uid].executed = true
		default:
			continue
		}
	}
}

func (s *superVisor) ExecJobs() error {
	close(s.addJobChan)
	<-s.addDoneChan
	s.logStatus()

	go s.checkDone()

	for _, j := range s.queue {
		s.jobExecChan <- j
	}
	log.Debug("all jobs sended")
	close(s.jobExecChan)
	s.logStatus()

	<-s.allDoneChan

	s.logStatus()
	return nil
}

func NewSuperVisor(poolSize int) *superVisor {
	sv := superVisor{
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
	go sv.statusLoop()

	log.Debug("supervisor init done")
	return &sv
}
