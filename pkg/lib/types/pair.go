package types

import "fmt"

type Pair struct {
	Key   string
	Value interface{}
}

func (p Pair) String() string {
	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}
