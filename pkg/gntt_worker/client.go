package gntt_worker

import (
	"gntt/pkg/gntt_math"
	"gntt/pkg/gntt_optional"
	"sync"

	log "github.com/sirupsen/logrus"
)

type ClientWorkerJob interface {
	NumJobs() gntt_optional.Optional[int]
	MaxClients() int
	Perform()
	CancelAll()
	JobsFinished()
}

func ClientWorker(job ClientWorkerJob) chan bool {

	var jobCount int = 0

	numj := job.NumJobs()
	isDone := func() bool {
		return (numj.HasValue() && (numj.Get() < jobCount))
	}

	maxTokens := job.MaxClients()
	if numj.HasValue() {
		maxTokens = gntt_math.Max(maxTokens, numj.Get())
	}

	connTokCh := make(chan bool, maxTokens)

	var wjobs sync.WaitGroup

	worker := func(id int) {
		// do work
		job.Perform()
		wjobs.Done()
		connTokCh <- true
	}

	endWork := make(chan bool)

	// work dispatcher
	go func() {

		// queue work tokens
		for i := 0; i < maxTokens; i++ {
			connTokCh <- true
		}
		for {
			select {
			case <-connTokCh:
				// dequeue a token
				jobCount++
				if !isDone() {
					wjobs.Add(1)
					go worker(jobCount)
				} else {
					connTokCh = nil
					job.JobsFinished()
				}
			case <-endWork:
				log.Tracef("end work")
				goto done
			}
		}

	done:
		job.CancelAll()
		wjobs.Wait()
		endWork <- true
	}()

	return endWork
}
