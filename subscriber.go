package main

type Observeable struct {
	url       string
	observers map[string][]chan string
}

func (o *Observeable) AddObserver(e string, ch chan string) {
	if o.observers == nil {
		//create one if not exists
		o.observers = make(map[string][]chan string)
	}
	if _, ok := o.observers[e]; ok {
		o.observers[e] = append(o.observers[e], ch)
	} else {
		o.observers[e] = []chan string{ch}
	}
}

func (o *Observeable) RemoveObserver(e string, ch chan string) {
	if _, ok := o.observers[e]; ok {
		for i := range o.observers[e] {
			if o.observers[e][i] == ch {
				o.observers[e] = append(o.observers[e][:i], o.observers[e][i + 1:]...)
				break
			}
		}
	}
}

func (o *Observeable) Emit(e string, response string) {
	if _, ok := o.observers[e]; ok {
		for _, handler := range o.observers[e] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}