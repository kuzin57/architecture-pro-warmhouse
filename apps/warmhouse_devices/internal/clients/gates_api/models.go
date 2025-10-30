package gatesapi

type GateState string

const (
	GateStateOpen   GateState = "open"
	GateStateClosed GateState = "closed"
)

func (g GateState) String() string {
	return string(g)
}

type StatusResponse struct {
	IsActive bool      `json:"is_active"`
	State    GateState `json:"state"`
}

type ChangeStateRequest struct {
	State GateState `json:"state"`
}
