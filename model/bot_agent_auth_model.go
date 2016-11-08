package model

type Prefix string

func (p Prefix) String() string {
	return string(p)
}
