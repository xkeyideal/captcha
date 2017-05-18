package pool

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"time"
)

// idLen is a length of captcha id string.
// (20 bytes of 62-letter alphabet give ~119 bits.)
const idLen = 20

// idChars are characters allowed in captcha id.
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// rngKey is a secret key used to deterministically derive seeds for
// PRNGs used in image and audio. Generated once during initialization.
var rngKey [32]byte

const (
	imageSeedPurpose = 0x01
)

func init() {
	if _, err := io.ReadFull(rand.Reader, rngKey[:]); err != nil {
		fmt.Println("captcha: error reading random source: " + err.Error())
		os.Exit(-1)
	}
}

func deriveSeed(purpose byte, id string, digits []byte) (out [16]byte) {
	var buf [sha256.Size]byte
	h := hmac.New(sha256.New, rngKey[:])
	h.Write([]byte{purpose})
	io.WriteString(h, id)
	h.Write([]byte{0})
	h.Write(digits)
	sum := h.Sum(buf[:0])
	copy(out[:], sum)
	return
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte, err error) {
	b = make([]byte, length)
	_, err = io.ReadFull(rand.Reader, b)
	if err != nil {
		err = errors.New("captcha: error reading random source: " + err.Error())
	}
	return
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte, err error) {
	if length == 0 {
		return nil, errors.New("length is zero")
	}
	if mod == 0 {
		return nil, errors.New("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r, err := randomBytes(length + (length / 4))
		if err != nil {
			return nil, err
		}
		for _, c := range r {
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return b, nil
			}
		}
	}

}

// randomId returns a new random id string.
func randomId() (string, error) {
	b, err := randomBytesMod(idLen, byte(len(idChars)))
	if err != nil {
		return "", err
	}
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b), nil
}

var StdChars = []byte("0123456789")

func randomWords(length int) []byte {
	result := make([]byte, length)
	seed := time.Now().UnixNano() + mrand.Int63n(10000)
	mrand.Seed(seed)
	sl := len(StdChars)
	for i := 0; i < length; i++ {
		result[i] = StdChars[mrand.Intn(sl)]
	}
	return result
}
