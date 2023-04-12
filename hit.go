package main

type Hit struct {
	Dot
	duration int
}

func (h *Hit) update() {
	h.y -= 1
	h.duration -= 1
}
