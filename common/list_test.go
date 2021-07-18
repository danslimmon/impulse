package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
func TestTextToBlipLine(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testCase struct {
		exp  blipLine
		text []byte
	}

	testCases := []testCase{
		testCase{
			exp:  blipLine{Indent: 0, Text: "hello"},
			text: []byte("hello"),
		},
		testCase{
			exp:  blipLine{Indent: 2, Text: "hello"},
			text: []byte("\t\thello"),
		},
		testCase{
			exp:  blipLine{Indent: 2, Text: "hello\tgoodbye\t\t"},
			text: []byte("\t\thello\tgoodbye\t\t"),
		},
	}

	for i, tc := range testCases {
		t.Logf("test case %d", i)
		rslt, err := textToBlipLine(tc.text)
		assert.Nil(err)
		assert.Equal(tc.exp.Indent, rslt.Indent)
		assert.Equal(tc.exp.Text, rslt.Text)
	}
}
*/

/*
func TestUnmarshalBlopList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rslt, err := UnmarshalBlopList("foo", []byte("hello\n"))
	assert.Nil(err)
	assert.NotNil(rslt)
	assert.Equal("foo", rslt.Name)
	assert.Equal(1, len(rslt.Blops))
	assert.Equal("hello", rslt.Blops[0].Name())
	assert.Equal(0, len(rslt.Blops[0].Root.Children))
}
*/
