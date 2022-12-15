package sessions

type Status struct {
	Alive      bool `json:"alive"`
	Rebouncing bool `json:"rebouncing"`
}
