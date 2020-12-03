/*
Copyright © 2020 MasseElch <info@masseelch.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package elk

const (
	OrderASC  = "asc"
	OrderDESC = "desc"
)

type (
	// Used on a schema to pass options to the handler generator.
	HandlerAnnotation struct {
		Skip             bool
		SkipCreate       bool
		SkipRead         bool
		SkipUpdate       bool
		SkipDelete       bool
		SkipList         bool
		CreateGroups     []string
		ReadGroups       []string
		UpdateGroups     []string
		ListGroups       []string
		DefaultListOrder []Order
	}
	// Used on fields pass options to the handler generator.
	FieldAnnotation struct {
		SkipCreate          bool
		CreateValidationTag string
		SkipUpdate          bool
		UpdateValidationTag string
	}
	// Used on (list) edges to specify a default order.
	EdgeAnnotation struct {
		DefaultOrder []Order
	}
	Order struct {
		Order string
		Field string
	}
)

func (HandlerAnnotation) Name() string {
	return "HandlerGen"
}

func (FieldAnnotation) Name() string {
	return "FieldGen"
}

func (EdgeAnnotation) Name() string {
	return "EdgeGen"
}
