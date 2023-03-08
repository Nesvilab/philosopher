package parallel

import "errors"

type Worker struct {
	p       *Parallel
	Name    string
	Config  *WorkerConfig
	execute func(wh *WorkerHelper, args interface{})
	helper  *WorkerHelper
}

type WorkerConfig struct {
	Parallelism int
}

func (p *Parallel) NewWorker(name string, cfg *WorkerConfig) (*Worker, error) {
	if _, exists := p.workers[name]; exists {
		return nil, errors.New("worker already exists")
	}
	if cfg.Parallelism < 1 {
		return nil, errors.New("parallelism must be 1 or higher")
	}
	w := &Worker{
		p:      p,
		Name:   name,
		Config: cfg,
	}
	p.workers[name] = w
	return w, nil
}

func (p *Parallel) Worker(name string) *Worker {
	if _, exists := p.workers[name]; !exists {
		return nil
	}
	return p.workers[name]
}

func (w *Worker) SetExecution(exec func(wh *WorkerHelper, args interface{})) {
	w.execute = exec
}

func (w *Worker) Start(args interface{}) {
	wh := newWorkerHelper(w)
	w.helper = wh
	for i := 0; i < w.Config.Parallelism; i++ {
		w.helper.wg.Add(1)
		go w.execute(wh, args)
	}
}

func (w *Worker) Wait() {
	w.helper.wg.Wait()
}

func (w *Worker) SetParallelism(p int) {
	w.Config.Parallelism = p
}
