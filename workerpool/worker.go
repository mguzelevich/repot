package workerpool

import (
	// "fmt"

	log "github.com/sirupsen/logrus"
)

type worker struct {
	uid       string
	stateChan chan poolStateSignal
}

func (w *worker) loop(input <-chan *job, workerDone chan string) {
	w.stateChanged("start")
	for job := range input {
		w.stateChanged("job started")
		job.start(w.uid)
		w.stateChanged("job finished")
	}
	w.stateChanged("finish")
	workerDone <- w.uid
}

func (w *worker) stateChanged(state string) {
	log.WithFields(log.Fields{"worker": w.uid, "state": state}).Debug("worker state changed")
	select {
	case w.stateChan <- poolStateSignal{uid: w.uid, msg: state}:
		// fmt.Println("sent message", msg)
	default:
		//fmt.Println("no message sent")
	}
}
