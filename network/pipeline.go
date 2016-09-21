package network

import (
	"container/list"
	"bytes"
	"fmt"
)

type pipelineItem struct {
	name    string
	handler interface{}
}
type Pipeline struct {
	handlers *list.List
}

func NewPipeline() *Pipeline {
	t := &Pipeline{
		handlers:list.New(),
	}
	return t
}

func (this *Pipeline) AddFirst(name string, handler interface{}) {
	item := &pipelineItem{
		name:name,
		handler:handler,
	}
	this.handlers.PushFront(item)
}
func (this *Pipeline) AddLast(name string, handler interface{}) {
	item := &pipelineItem{
		name:name,
		handler:handler,
	}
	this.handlers.PushBack(item)
}

func (this *Pipeline) Replace(name string, handler interface{}) {
	item := &pipelineItem{
		name:name,
		handler:handler,
	}
	for e := this.handlers.Front(); e != nil; e = e.Next() {
		if e.Value.(*pipelineItem).name == name {
			this.handlers.InsertAfter(item, e)
			this.handlers.Remove(e)
			break
		}
	}
}

func (this *Pipeline) Info() string {
	buff := new(bytes.Buffer)

	for e := this.handlers.Front(); e != nil; e = e.Next() {
		buff.WriteString(e.Value.(*pipelineItem).name)
		buff.WriteString(":")
		buff.WriteString(fmt.Sprintf("%v", &e.Value.(*pipelineItem).handler))
		buff.WriteString("|")
	}
	return buff.String()
}
func (this *Pipeline) FireConnect(chl *Channel) {
	for e := this.handlers.Front(); e != nil; e = e.Next() {
		if evt, ok := e.Value.(*pipelineItem).handler.(ConnectEvent); ok {
			evt.OnConnect(chl)
		}
	}
}
func (this *Pipeline) FireRead(chl *Channel, packet Packet) {
	t := packet
	for e := this.handlers.Front(); e != nil; e = e.Next() {
		if evt, ok := e.Value.(*pipelineItem).handler.(ReadEvent); ok {
			t = evt.OnRead(chl, t)
		}
	}
}
func (this *Pipeline) FireWrite(chl *Channel, packet Packet) ExchangePacket {
	t := packet
	for e := this.handlers.Back(); e != nil; e = e.Prev() {
		if evt, ok := e.Value.(*pipelineItem).handler.(WriteEvent); ok {
			t = evt.OnWrite(chl, t)
		}
	}
	return t
}
func (this *Pipeline) FireClose(chl *Channel) {
	for e := this.handlers.Front(); e != nil; e = e.Next() {
		if evt, ok := e.Value.(*pipelineItem).handler.(CloseEvent); ok {
			evt.OnClose(chl)
		}
	}
}
