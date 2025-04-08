package model

type Status int

const (
	Successful Status = iota
	Failed
	Queued
	Running
	Stopped
)

func (s Status) String() string {
    return [...]string{"Successful", "Failed", "Queued", "Running", "Stopped"}[s]
}

type History struct {
	Id string `json:"id"`
	Status Status `json:"status"`
	Logs string `json:"logs"`
}
