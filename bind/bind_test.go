package bind

import (
	"reflect"
	"testing"
	"time"
)

func TestBind(t *testing.T) {
	type Anonymous struct {
		A1 int    `json:"a1"`
		A2 string `json:"a2"`
	}
	type T struct {
		F1 int     `json:"f1"`
		F2 string  `json:"f2"`
		F3 int     `json:"f3"`
		F4 float64 `json:"f4"`
		F5 bool    `json:"f5"`
		F6 int32   `json:"f6"`
		F7 float32 `json:"f7"`
		Anonymous
	}
	vals := []struct {
		Query  string
		Origin T
		Target T
	}{{
		Query:  "f1=123&f2=abcd&f3=456&f4=1.89&f5=true&f6=12&f7=34.5&f8=",
		Target: T{F1: 123, F2: "abcd", F3: 456, F4: 1.89, F5: true, F6: 12, F7: 34.5},
	}, {
		Query:  "f1=23&f2=a45d&f3=83&f4=12.9&f5=false&f6=66&f7=55.5&f8=",
		Target: T{F1: 23, F2: "a45d", F3: 83, F4: 12.9, F5: false, F6: 66, F7: 55.5},
	}, {
		Query:  "f1=123&a1=10",
		Target: T{F1: 123, Anonymous: Anonymous{A1: 10}},
	}}

	start := time.Now()
	for i := range vals {
		v := &vals[i]
		_ = Bind(v.Query, &v.Origin)
		if !reflect.DeepEqual(v.Origin, v.Target) {
			t.Error(v.Origin, v.Target)
		}
	}

	t.Log(time.Since(start))
}
