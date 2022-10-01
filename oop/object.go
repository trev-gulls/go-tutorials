package main

type IObject interface {
	TotalValue() float64
	AddOne()
	SetValue(v float64)
}

type Object struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Count int64   `json:"count"`
	Value float64 `json:"value"`
}

func (o Object) TotalValue() float64 {
	return float64(o.Count) * o.Value
}

func (o Object) AddOne() {
	o.Count += 1
}

func (o Object) SetValue(v float64) {
	o.Value = v
}
