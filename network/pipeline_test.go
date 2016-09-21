package network

import (
	"testing"
)

func testPipeline() *Pipeline {
	return NewPipeline()
}

func TestPipeline_All(t *testing.T) {
	h1 := i_am_a_simple_handler{}
	h2 := i_am_a_simple_handler{}
	h3 := i_am_a_simple_handler{}
	h4 := i_am_a_simple_handler{}
	p := testPipeline()
	p.AddFirst("handler2", h2)
	p.AddFirst("handler1", h1)
	t.Log(p.Info())
	p.AddFirst("handler3", h3)
	t.Log(p.Info())
	p.Replace("handler2", h4)
	t.Log(p.Info())
}
