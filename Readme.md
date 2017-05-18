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

	NewCaptchaPool(ctx, width, height, wordLength, poolsize, parallelNum, imageType int)

Creates a new captcha pool

1. ctx context.Background() 
2. width, height: image's width and height
3. wordLength: generate words' length
4. poolsize: buffer size
5. parallelNum: goroutine number
6. imageType: PNG or JPEG

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
    
    CaptchaPool = pool.NewCaptchaPool(context.Background(), 240, 80, 6, 10, 1, 2)
    
    captchaBody := CaptchaPool.GetImage()
	
	CaptchaPool.Stop()
    
```
See detail in file [captcha.go](https://github.com/xkeyideal/captcha/blob/master/captcha/captcha.go)
