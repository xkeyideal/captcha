package pool

import (
	"bytes"

	"github.com/satori/go.uuid"
)

const (
	PngImage  = 1
	JpegImage = 2
)

type CaptchaBody struct {
	Id   string
	Data *bytes.Buffer
	Val  []byte
}

type CaptchaPool struct {
	captchaBuffer chan *CaptchaBody
	wordsBuffer   chan []byte
	wordLength    int
	width         int
	height        int
	imageType     int

	numWidth  int
	numHeight int
	dotSize   int

	//	ctx    context.Context
	//	cancel context.CancelFunc
	//	wg     sync.WaitGroup
}

func NewCaptchaPool(width, height, wordLength, poolsize, parallelNum, imageType int) *CaptchaPool {

	numWidth, numHeight, dotSize := calculateSizes(width, height, wordLength)

	pool := &CaptchaPool{
		wordLength:    wordLength,
		width:         width,
		height:        height,
		captchaBuffer: make(chan *CaptchaBody, poolsize),
		wordsBuffer:   make(chan []byte, poolsize),
		imageType:     imageType,

		numWidth:  numWidth,
		numHeight: numHeight,
		dotSize:   dotSize,

		//wg: sync.WaitGroup{},
	}

	//pool.ctx, pool.cancel = context.WithCancel(ctx)

	go pool.GenRandomWords()

	for i := 0; i < parallelNum; i++ {
		go pool.GenImage()
	}

	return pool
}

func calculateSizes(width, height, ncount int) (numWidth int, numHeight int, dotSize int) {
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
	dotSize = int(nh / fh)
	if dotSize < 1 {
		dotSize = 1
	}
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	numWidth = int(nw) - dotSize
	numHeight = int(nh)

	return numWidth, numHeight, dotSize
}

func (p *CaptchaPool) GenRandomWords() {
	//	defer p.wg.Done()

	//	p.wg.Add(1)

	for {
		words := randomWords(p.wordLength)

		//		select {
		//		default:
		//		case <-p.ctx.Done():
		//			//fmt.Println("Gen Random Words Stop")
		//			return
		//		}

		p.wordsBuffer <- words
	}
}

func (p *CaptchaPool) genImage() (*bytes.Buffer, []byte, error) {
	id, err := randomId()
	if err != nil {
		return nil, nil, err
	}

	words := <-p.wordsBuffer

	imgBuffer := new(bytes.Buffer)
	img := NewCImage(p.numWidth, p.numHeight, p.dotSize)
	img.ProductCImage(id, words, p.width, p.height)

	if p.imageType == PngImage {
		_, err = img.WriteToPng(imgBuffer)
	} else {
		_, err = img.WriteToJpeg(imgBuffer)
	}

	if err != nil {
		return nil, words, err
	}

	return imgBuffer, words, nil
}

func (p *CaptchaPool) GenImage() {
	//	defer p.wg.Done()

	//	p.wg.Add(1)

	for {

		imgBytes, words, err := p.genImage()

		//		select {
		//		default:
		//		case <-p.ctx.Done():
		//			//fmt.Println("GenImage Stop", num)
		//			return
		//		}

		if err == nil {
			captchaBody := &CaptchaBody{
				Id:   uuid.NewV4().String(),
				Data: imgBytes,
				Val:  words,
			}

			p.captchaBuffer <- captchaBody
		}
	}
}

func (p *CaptchaPool) GetImage() *CaptchaBody {
	return <-p.captchaBuffer
}

func (p *CaptchaPool) Stop() {
	//	if p.cancel != nil {
	//		p.cancel()
	//		p.cancel = nil
	//	}

	//	p.wg.Wait()

	close(p.captchaBuffer)
	close(p.wordsBuffer)
}
