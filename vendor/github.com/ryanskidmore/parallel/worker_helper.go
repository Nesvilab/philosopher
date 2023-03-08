package parallel

import (
	"errors"
	"sync"
)

type WorkerHelper struct {
	worker *Worker
	wg     *sync.WaitGroup
}

func newWorkerHelper(w *Worker) *WorkerHelper {
	wg := &sync.WaitGroup{}
	return &WorkerHelper{worker: w, wg: wg}
}

func (wh *WorkerHelper) Done() {
	wh.wg.Done()
}

func (wh *WorkerHelper) PublishData(name string, data interface{}) error {
	if _, exists := wh.worker.p.dataChannels[name]; !exists {
		return errors.New("Data channel does not exist")
	}
	go func() {
		wh.worker.p.dataChannels[name] <- data
	}()
	return nil
}

func (wh *WorkerHelper) ConsumeData(name string) (interface{}, error) {
	if _, exists := wh.worker.p.dataChannels[name]; !exists {
		return nil, errors.New("Data channel does not exist")
	}
	data, open := <-wh.worker.p.dataChannels[name]
	if !open {
		return nil, errors.New("Data channel closed")
	}
	return data, nil
}

func (wh *WorkerHelper) ConsumeDataInBatches(name string, size int) ([]interface{}, error) {
	if _, exists := wh.worker.p.dataChannels[name]; !exists {
		return nil, errors.New("Data channel does not exist")
	}
	dataBatch := make([]interface{}, size, size)
	for i := 0; i < size; i++ {
		data, open := <-wh.worker.p.dataChannels[name]
		if !open {
			return dataBatch, errors.New("Data channel closed")
		}
		dataBatch[i] = data
	}
	return dataBatch, nil
}
