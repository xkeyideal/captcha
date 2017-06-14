Package captcha
=====================

Package captcha implements generation and verification of image
CAPTCHAs.

A captcha solution is the sequence of digits 0-9 with the defined length.

An image representation is a PNG-encoded or JPEG-encoded  image with the solution printed on
it in such a way that makes it hard for computers to solve it using OCR.

This package doesn't require external files or libraries to generate captcha
representations; it is self-contained.

Package code refer from [dchest/captcha](https://github.com/dchest/captcha)

Advantages:

1. High-Performance：Generation captcha use goroutine + channel,  get captcha ready in advance by channel
2. Change panic to error, avoid runtime panic
3. Not inline store interface, can use any store method such as Redis, Memcache, Memory and so on after get captcha image
4. Use uuid instead of original random id avoid conflict
5. Add Context to control generate captcha goroutine, can stop generate programming active


Examples
--------

![Image](https://github.com/xkeyideal/captcha/raw/master/image/exampleimage.png)

Functions
---------

### func NewCaptchaPool

	NewCaptchaPool(width, height, wordLength, poolsize, parallelNum, imageType int)

Creates a new captcha pool

1. width, height: image's width and height
2. wordLength: generate words' length
3. poolsize: buffer size
4. parallelNum: goroutine number
5. imageType: PNG or JPEG

### func Stop

	func (p *CaptchaPool) Stop()

Stop CaptchaPool active

Usage
--------
```go
    type CaptchaBody struct {
    	Id   string
    	Data *bytes.Buffer
    	Val  []byte
    }
    
    CaptchaPool = pool.NewCaptchaPool(240, 80, 6, 10, 1, 2)
    
    captchaBody := CaptchaPool.GetImage()
	
	CaptchaPool.Stop()
    
```
See detail in file [captcha.go](https://github.com/xkeyideal/captcha/blob/master/captcha.go)

Golang context can't control goroutine channel will deadlock by using sync.WaitGroup to wait all goroutine return and close channels.
