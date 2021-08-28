package greeter

import "fmt"

type Greeter struct {
	salutation string
}

func NewGreeter(salutation string) *Greeter {
	return &Greeter{salutation: salutation}
}

func (g *Greeter) Greet(name string)  {
	fmt.Printf("%s, %s", g.salutation, name) //nolint:forbidigo
}
