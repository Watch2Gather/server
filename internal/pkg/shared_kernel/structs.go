package sharedkernel

type Instruction struct {
	State  string `json:"state"`
	TimeCode  int    `json:"timecode"`
}
type Message struct {
	Instruction *Instruction
	Sender      string
}
