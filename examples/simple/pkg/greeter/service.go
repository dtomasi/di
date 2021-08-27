package greeter

import "fmt"

// Greeter This is that one greeter example that you know from everywhere
type Greeter struct {
	baseGreeting string
}

func NewGreeter(baseGreeting string) *Greeter {
	return &Greeter{
		baseGreeting: baseGreeting,
	}
}

func (g *Greeter) Greet(name string) {
	fmt.Printf("%s %s\n", g.baseGreeting, name)
}