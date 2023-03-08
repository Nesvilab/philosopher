package parallel

type Parallel struct {
	workers      map[string]*Worker
	dataChannels map[string]chan interface{}
}

func New() *Parallel {
	workers := make(map[string]*Worker)
	dataChannels := make(map[string]chan interface{})
	return &Parallel{workers: workers, dataChannels: dataChannels}
}
