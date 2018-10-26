package gopter

import "testing"

// Properties is a collection of properties that should be checked in a test
type Properties struct {
	parameters *TestParameters
	props      map[string]Prop
	propNames  []string
}

// NewProperties create new Properties with given test parameters.
// If parameters is nil default test parameters will be used
func NewProperties(parameters *TestParameters) *Properties {
	if parameters == nil {
		parameters = DefaultTestParameters()
	}
	return &Properties{
		parameters: parameters,
		props:      make(map[string]Prop, 0),
		propNames:  make([]string, 0),
	}
}

// Property add/defines a property in a test.
func (p *Properties) Property(name string, prop Prop) {
	p.propNames = append(p.propNames, name)
	p.props[name] = prop
}

// Run checks all definied properties
// Useful for passing directly to testing.T#Run
func (p *Properties) Run(t *testing.T) {
	for _, propName := range p.propNames {
		prop := p.props[propName]

		t.Run(propName, func(t *testing.T) { prop.Check(t, p.parameters) })
	}
}
