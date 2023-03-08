package parallel

import "errors"

func (p *Parallel) NewDataChannel(name string) error {
	if _, exists := p.dataChannels[name]; exists {
		return errors.New("Data channel already exists")
	}
	p.dataChannels[name] = make(chan interface{})
	return nil
}

func (p *Parallel) CloseDataChannel(name string) error {
	if _, exists := p.dataChannels[name]; !exists {
		return errors.New("Data channel doesn't exists")
	}
	close(p.dataChannels[name])
	delete(p.dataChannels, name)
	return nil
}
