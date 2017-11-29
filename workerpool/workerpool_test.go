package workerpool

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var testJobOk = func(uid string) error {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return nil
}

var testJobFailed = func(uid string) error {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return errors.New("failed")
}

func TestWorkerPool_addJobs(t *testing.T) {
	sv := NewWP(3)

	sv.AddJob(fmt.Sprintf("uid %v", 1), testJobOk)
	sv.AddJob(fmt.Sprintf("uid %v", 2), testJobOk)

	// if len(sv.queue) != 2 {
	// 	t.Error("Incorrect task queue size")
	// }

	// stats := sv.status()
	// for _, s := range []jobStatus{finished, failed, executing} {
	// 	if stats[s] != 0 {
	// 		t.Error("Incorrect ", s, " task counter")
	// 	}
	// }
	// if stats[pending] != 2 {
	// 	t.Error("Incorrect ", pending, " task counter")
	// }
}

func TestWorkerPool_addExecJobs(t *testing.T) {
	sv := NewWP(3)
	sv.AddJob(fmt.Sprintf("uid %v", 1), testJobOk)
	sv.ExecJobs()
	// if err := sv.AddJob(fmt.Sprintf("uid %v", 2), testJobOk); err == nil {
	// 	t.Error("Add after exec")
	// }
}

func TestSupervisor_execJobs(t *testing.T) {
	sv := NewWP(3)
	for i := 0; i < 5; i++ {
		sv.AddJob(fmt.Sprintf("uid %v", i*2), testJobOk)
		sv.AddJob(fmt.Sprintf("uid %v", i*2+1), testJobFailed)
	}
	sv.ExecJobs()

	// stats := sv.status()
	// if stats[finished] != 1 {
	// 	t.Error("Incorrect finished task counter")
	// }
	// if stats[failed] != 1 {
	// 	t.Error("Incorrect failed task counter")
	// }
	// if stats[executing] != 0 {
	// 	t.Error("Incorrect executing task counter")
	// }
	// if stats[pending] != 8 {
	// 	t.Error("Incorrect pending task counter")
	// }
	// if len(sv.queue) != 10 {
	// 	t.Error("Incorrect task queue size")
	// }
}

// supervisor := superVisor{}

// for i, repository := range manifest.Repositories {
// 	idx := i + 1
// 	if i%Jobs == 0 {
// 		supervisor.wgWait()
// 	}

// 	supervisor.taskStarted(idx, repository)
// 	go func(idx int, repository *Repository) {
// 		defer supervisor.taskDone(idx, repository)

// 		clonePath := filepath.Join(rootPath, repository.Path, repository.Name)
// 		git.Clone(repository.Repository, clonePath)
// 	}(idx, repository)
// 	//		}
// }
// // fmt.Fprintf(os.Stderr, "WG: final wait [%s]\n", processedIdx)
// supervisor.wg.Wait()
