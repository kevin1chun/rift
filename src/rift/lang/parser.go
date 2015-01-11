package lang

import (
	"io"
	"io/ioutil"
)

// TODO: Use a reader instead
func Parse(source io.Reader) (*riftParser, error) {
	readSource, sourceErr := ioutil.ReadAll(source)
	if sourceErr != nil {
		return nil, sourceErr
	}

	parser := &riftParser{Buffer: string(readSource[:])}
	parser.Init()
	err := parser.Parse()
	if err != nil {
		return parser, err
	}
	parser.Execute()

	return parser, nil
}