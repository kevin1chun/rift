package lang

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

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

// TODO: Can this work any better?
func GetSyntaxErrors(p *riftParser) string {
	var errors []string
	for _, err := range p.Error() {
		pos := translatePositions(p.Buffer, []int{int(err.begin), int(err.end)})[0]
		errors = append(errors, fmt.Sprintf("Line %d, character %d", pos.line, pos.symbol))
	}
	return strings.Join(errors, "\n")
}