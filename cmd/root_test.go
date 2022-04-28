package cmd

import (
	"math"
	"testing"

	"github.com/gerow/go-color"
	"github.com/stretchr/testify/assert"
)

func TestHEXToHSL(t *testing.T) {
	r1, _ := HEXToHSL("#000000")
	assert.Equal(t, color.HSL{0.0, 0.0, 0.0}, r1)
	r2, _ := HEXToHSL("#ff0000")
	assert.Equal(t, color.HSL{0.0, 1.0, 0.5}, r2)
	_, err1 := HEXToHSL("#xf0000")
	assert.NotNil(t, err1)
}

func TestDistance(t *testing.T) {
	assert.Equal(t, 0.0, distance(color.HSL{H: 0, S: 0, L: 0}, color.HSL{H: 0, S: 0, L: 0}))
	assert.Equal(t, 0.0, distance(color.HSL{H: 1, S: 0, L: 0}, color.HSL{H: 0, S: 0, L: 0}))
	assert.Equal(t, 1.0, distance(color.HSL{H: 0, S: 0, L: 1}, color.HSL{H: 0, S: 0, L: 0}))
	assert.Equal(t, math.Sqrt(2.0), distance(color.HSL{H: 0, S: 1, L: 1}, color.HSL{H: 0, S: 0, L: 0}))
}

func TestGetClosestColor(t *testing.T) {
	c := &[]HEX_HSL{
		{HEX: "#ffffff", HSL: color.HSL{H: 0, S: 0, L: 1}},
		{HEX: "#000000", HSL: color.HSL{H: 0, S: 0, L: 0}},
	}

	r1, _ := getClosestColor("#030303", c, 0.5)
	assert.Equal(t, HEX("#000000"), r1)

	r2, _ := getClosestColor("#030303", c, 0.01)
	assert.Equal(t, HEX("#030303"), r2)
}

func TextGenColorList(t *testing.T) {
	r, _ := genColorList(&[]HEX{"#000000"})
	assert.Equal(t, 1, len(*r))
}

func TestGetColors(t *testing.T) {
	assert.Equal(t, HEX("#000000"), getColors("#000000")[0])
	assert.Equal(t, HEX("#000000"), getColors("color:#000000;")[0])
}
