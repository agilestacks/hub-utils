package cmd

import "time"

type Phase struct {
	Phase  string `json:"phase"`
	Status string `json:"status"`
}

type LatestOP struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Initiator string    `json:"initiator"`
	Timestamp time.Time `json:"timestamp"`
	Phases    []Phase   `json:"phases"`
}

type Component struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type StateLocation struct {
	Uri  string `json:"uri"`
	Kind string `json:"kind"`
}

type State struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Status        string        `json:"status"`
	LatestOP      LatestOP      `json:"latestOperation"`
	StateLocation StateLocation `json:"stateLocation"`
	Components    []Component   `json:"components"`
}
