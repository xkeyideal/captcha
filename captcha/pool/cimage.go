package pool

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand"
)

const (
	// Maximum absolute skew factor of a single digit.
	maxSkew = 0.7
	// Number of background circles.
	circleCount = 20

	colorCount = 20
)

type CImage struct {
	*image.Paletted
	numWidth  int
	numHeight int
	dotSize   int
	rng       siprng
}

func (m *CImage) ProductCImage(id string, words []byte, width, height int) {

	// Initialize PRNG.
	m.rng.Seed(deriveSeed(imageSeedPurpose, id, words))

	//图像调色板
	m.Paletted = image.NewPaletted(image.Rect(0, 0, width, height), m.getRandomPalette())

	// Randomly position captcha inside the image.
	maxx := width - (m.numWidth+m.dotSize)*len(words) - m.dotSize
	maxy := height - m.numHeight - m.dotSize*2

	var border int
	if width > height {
		border = height / 5
	} else {
		border = width / 5
	}
	x := m.rng.Int(border, maxx-border)
	y := m.rng.Int(border, maxy-border)

	// Draw words.
	for _, i := range words {
		m.drawWords(font[i-48], x, y)
		x += m.numWidth + m.dotSize
	}

	// Draw strike-through line.
	m.strikeThrough()
	// Apply wave distortion.
	m.distort(m.rng.Float(5, 10), m.rng.Float(100, 200))
	// Fill image with random circles.
	//m.drawNoiseLine(5)
	m.fillWithCircles(circleCount, m.dotSize)
}

func NewCImage(numWidth, numHeight, dotSize int) *CImage {
	m := &CImage{
		numWidth:  numWidth,
		numHeight: numHeight,
		dotSize:   dotSize,
	}

	return m
}

// WriteTo writes captcha image in PNG format into the given writer.
func (m *CImage) WriteToPng(w io.Writer) (int64, error) {
	b, err := m.encodedPNG()
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

// encodedPNG encodes an image to PNG and returns
// the result as a byte slice.
func (m *CImage) encodedPNG() ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m.Paletted); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteTo writes captcha image in PNG format into the given writer.
func (m *CImage) WriteToJpeg(w io.Writer) (int64, error) {
	b, err := m.encodedJPEG()
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

// encodedJPEG encodes an image to JPEG and returns
// the result as a byte slice.
func (m *CImage) encodedJPEG() ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, m.Paletted, &jpeg.Options{10}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *CImage) fillWithCircles(n, maxradius int) {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		colorIdx := uint8(m.rng.Int(1, circleCount-1))
		r := m.rng.Int(1, maxradius)
		m.drawCircle(m.rng.Int(r, maxx-r), m.rng.Int(r, maxy-r), r, colorIdx)
	}
}

func (m *CImage) distort(amplude float64, period float64) {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y

	oldm := m.Paletted
	newm := image.NewPaletted(image.Rect(0, 0, w, h), oldm.Palette)

	dx := 2.0 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			newm.SetColorIndex(x, y, oldm.ColorIndexAt(x+int(xo), y+int(yo)))
		}
	}
	m.Paletted = newm
}

func (m *CImage) strikeThrough() {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	y := m.rng.Int(maxy/3, maxy-maxy/3)

	amplitude := m.rng.Float(5, 20)

	period := m.rng.Float(80, 180)
	dx := 2.0 * math.Pi / period
	for x := 0; x < maxx; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		//for yn := 0; yn < 10; yn++ {
		for yn := 0; yn < m.dotSize; yn++ {
			r := m.rng.Int(0, m.dotSize)
			m.drawCircle(x+int(xo), y+int(yo)+(yn*m.dotSize), r/2, 1)
			//m.drawCircle(x+int(xo), y+int(yo)+(yn*11), r/2, 1)
		}
	}
}

func (m *CImage) drawWords(digit []byte, x, y int) {
	skf := m.rng.Float(-maxSkew, maxSkew)
	xs := float64(x)
	r := m.dotSize / 2
	y += m.rng.Int(-r, r)
	for yo := 0; yo < fontHeight; yo++ {
		for xo := 0; xo < fontWidth; xo++ {
			if digit[yo*fontWidth+xo] != blackChar {
				continue
			}

			m.drawCircle(x+xo*m.dotSize, y+yo*m.dotSize, r, 1)
		}
		xs += skf
		x = int(xs)
	}
}

func (m *CImage) drawHorizLine(fromX, toX, y int, colorIdx uint8) {
	for x := fromX; x <= toX; x++ {
		m.SetColorIndex(x, y, colorIdx)
	}
}

func (m *CImage) drawNoiseLine(num int) {
	for i := 0; i < num; i++ {
		x := rnd(0, m.Bounds().Max.X)
		m.drawHorizLine(
			int(float32(x)/1.5),
			x,
			rnd(0, m.Bounds().Max.Y),
			uint8(rnd(1, colorCount)),
		)
	}
}

// return a number bettwen from and to
func rnd(from, to int) int {
	return rand.Intn(to+1-from) + from
}

func (m *CImage) drawCircle(x, y, radius int, colorIdx uint8) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xo := 0
	yo := radius

	m.SetColorIndex(x, y+radius, colorIdx)
	m.SetColorIndex(x, y-radius, colorIdx)
	m.drawHorizLine(x-radius, x+radius, y, colorIdx)

	for xo < yo {
		if f >= 0 {
			yo--
			dfy += 2
			f += dfy
		}
		xo++
		dfx += 2
		f += dfx
		m.drawHorizLine(x-xo, x+xo, y+yo, colorIdx)
		m.drawHorizLine(x-xo, x+xo, y-yo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y+xo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y-xo, colorIdx)
	}
}

func (m *CImage) calculateSizes(width, height, ncount int) {
	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 4
	} else {
		border = width / 4
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)
	h := float64(height - border*2)
	// fw takes into account 1-dot spacing between digits.
	fw := float64(fontWidth + 1)
	fh := float64(fontHeight)
	nc := float64(ncount)
	// Calculate the width of a single digit taking into account only the
	// width of the image.
	nw := w / nc
	// Calculate the height of a digit from this width.
	nh := nw * fh / fw
	// Digit too high?
	if nh > h {
		// Fit digits based on height.
		nh = h
		nw = fw / fh * nh
	}
	// Calculate dot size.
	m.dotSize = int(nh / fh)
	if m.dotSize < 1 {
		m.dotSize = 1
	}
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	m.numWidth = int(nw) - m.dotSize
	m.numHeight = int(nh)
	fmt.Println("DotSize: ", m.dotSize)
}

func (m *CImage) getRandomPalette() color.Palette {
	p := make([]color.Color, circleCount+1)
	// Transparent color.
	p[0] = color.RGBA{0xFF, 0xFF, 0xFF, 0xff}
	// Primary color.
	prim := color.RGBA{
		uint8(m.rng.Intn(255)),
		uint8(m.rng.Intn(255)),
		uint8(m.rng.Intn(255)),
		0xFF,
	}
	p[1] = prim

	// Circle colors.
	for i := 2; i <= circleCount; i++ {
		p[i] = m.randomBrightness(prim, 255)
	}
	return p
}

func (m *CImage) randomBrightness(c color.RGBA, max uint8) color.RGBA {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return c
	}
	n := m.rng.Intn(int(max-maxc)) - int(minc)
	return color.RGBA{
		uint8(int(c.R) + n),
		uint8(int(c.G) + n),
		uint8(int(c.B) + n),
		uint8(c.A),
	}
}

func min3(x, y, z uint8) (m uint8) {
	m = x
	if y < m {
		m = y
	}
	if z < m {
		m = z
	}
	return
}

func max3(x, y, z uint8) (m uint8) {
	m = x
	if y > m {
		m = y
	}
	if z > m {
		m = z
	}
	return
}
