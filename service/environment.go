package service

type Environment struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Value string `json:"value"`
}

func (t Environment) ID() string {
	return t.Name
}
