package elk

import "fmt"

const (
	None Policy = iota
	Exclude
	Expose
)

type Policy uint

func (p Policy) Validate() error {
	if p > Expose {
		return fmt.Errorf("invalid policy %q", p)
	}
	return nil
}
