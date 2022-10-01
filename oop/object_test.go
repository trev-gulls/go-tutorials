package main

import "testing"

func TestObject_TotalValue(t *testing.T) {
	o := Object{Value: 0.3, Count: 5}
	want := 1.5
	got := o.TotalValue()
	if got != want {
		t.Errorf("TotalValue(): %f, want %f", got, want)
	}
}

func TestObject_AddOne(t *testing.T) {
	o := Object{}
	for _, v := range []int64{0, 1, 2, 3} {
		if o.Count != v {
			t.Errorf("Count: %d, want %d", o.Count, v)
		}
		o.AddOne()
	}
}

func TestObject_SetValue(t *testing.T) {
	o := Object{}
	o.SetValue(123.456)
	if o.Value != 123.456 {
		t.Errorf("Value: %f, want %f", o.Value, 123.456)
	}
}
