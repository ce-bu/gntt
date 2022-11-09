package gntt_worker

import "sync"

type ServerWorkerJob[T any] interface {
	Setup() error
	MaxConcurrentJobs() int
	Accept() (T, error)
	Perform(job T)
	Teardown()
	CancelAll()
}

func ServerWorker[T any](job ServerWorkerJob[T]) (chan bool, error) {
	err := job.Setup()
	if err != nil {
		return nil, err
	}

	jobChan := make(chan T, job.MaxConcurrentJobs())

	endAcc := func(jobChan chan T) chan bool {
		endAcc := make(chan bool)
		go func() {
			for {
				newj, err := job.Accept()
				if err != nil {
					jobChan = nil
				}
				select {
				case jobChan <- newj:
				case <-endAcc:
					endAcc <- true
					return
				}
			}
		}()
		return endAcc
	}(jobChan)

	var wjobs sync.WaitGroup

	endJobs := make(chan bool)

	go func() {
		for {
			select {
			case newj := <-jobChan:
				wjobs.Add(1)
				go func() {
					job.Perform(newj)
					wjobs.Done()
				}()
			case <-endJobs:
				job.Teardown()
				endAcc <- true
				goto end
			}
		}
	end:
		wjobs.Wait()
		<-endAcc
		endJobs <- true
	}()
	return endJobs, nil
}
