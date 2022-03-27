package store

type Score struct {
	Score   float32 `json:"score"`
	Comment string  `json:"comment"`
}

type Skills struct {
	Total    *Score `json:"total,omitempty"`
	Forehand *Score `json:"forehand,omitempty"`
	Backhand *Score `json:"backhand,omitempty"`
}

type Review struct {
	Author *Coach  `json:"author"`
	Skills *Skills `json:"skills,omitempty"`
}

type Player struct {
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	Reviews    []*Review `json:"reviews,omitempty"`
	TotalScore *Score    `json:"total_score,omitempty"`
}

type Coach struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
