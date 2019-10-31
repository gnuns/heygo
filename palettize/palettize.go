package palettize

import (
	"encoding/hex"
	"errors"
	"image/color"
	"math"
	"strings"
)

var ErrInvalidColorHexString = errors.New("invalid color hex string")

type Color color.RGBA

func ColorFromHex(hc string) (Color, error) {
	hc = strings.TrimPrefix(hc, "#")
	if len(hc) != 6 {
		return Color{}, ErrInvalidColorHexString
	}

	c, err := hex.DecodeString(hc)
	if err != nil {
		return Color{}, err
	}

	return Color{R: c[0], G: c[1], B: c[2]}, nil
}

func (c Color) Distance(c2 Color) float64 {
	rmean := (int64(c.R) + int64(c2.R)) / 2
	r := int64(c.R) - int64(c2.R)
	g := int64(c.G) - int64(c2.G)
	b := int64(c.B) - int64(c2.B)

	return math.Sqrt(float64((((512 + rmean) * r * r) >> 8) + 4*g*g + (((767 - rmean) * b * b) >> 8)))
}

func (c Color) Hex() string {
	hc := hex.EncodeToString([]byte{c.R, c.G, c.B})

	return "#" + strings.ToUpper(hc)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (c Color) IsGray() bool {
	maxDistance := max(
		abs(int(c.R)-int(c.G)),
		abs(int(c.R)-int(c.B)),
	)

	maxDistance = max(
		maxDistance,
		abs(int(c.B)-int(c.G)),
	)

	return maxDistance <= 5
}

type Palette []Color

func (p Palette) Index(c Color) int {
	ret, bestDist := 0, float64(-1)
	for i, v := range p {
		dist := v.Distance(c)
		if dist == 0 {
			return i
		}
		if (bestDist == -1) || (dist < bestDist) {
			ret, bestDist = i, dist
		}
	}
	return ret
}

func (p Palette) Convert(c Color) Color {
	if len(p) == 0 {
		return c
	}

	return p[p.Index(c)]
}
