package enigma

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type ScopeSalt struct {
	Salt
	once     sync.Once //
	genValue []byte    //
}

func (salt *ScopeSalt) GenSalt() []byte {
	salt.once.Do(func() {
		salt.genValue = salt.Salt.GenSalt()
	})

	return salt.genValue
}

type Salt struct {
	value []byte //
	has   bool   //
	len   int    //
}

func (salt *Salt) SetValue(b []byte) *Salt {
	*salt = Salt{value: b, has: true}

	return salt
}

func (salt Salt) GenSalt() []byte {
	switch salt.Has() {
	case false:
		b, err := RandBytes(salt.len)
		if err != nil {
			panic(errors.Wrapf(err, "new salt by randbytes %v",
				logs.KVL(
					"salt_size", salt.len,
				)))
		}
		return b
	default:
		return safeSalt(salt.value, salt.len)
	}
}

func (salt Salt) Scope(fn func(*ScopeSalt) error) error {
	return fn(&ScopeSalt{Salt: salt})
}

func (salt Salt) Has() bool {
	return salt.has
}

func (salt Salt) Len() int {
	return salt.len
}

func (salt *Salt) SetLen(n int) *Salt {
	salt.len = n

	return salt
}

func RandBytes(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = io.ReadFull(rand.Reader, b)
	return
}

func safeSalt(salt []byte, n int) (salt_ []byte) {

	salt_ = make([]byte, n)
	copy(salt_, salt)

	return
}

func SaltEncodeRule(src []byte, salt []byte, has bool) (src_ []byte) {
	switch has {
	case false:
		src_ = make([]byte, len(src)+len(salt))
		copy(src_, append(salt, src...))
		// src_ = append(salt.Salt(), src...)
	default:
		src_ = make([]byte, len(src))
		copy(src_, src)
		// src_ = src
	}

	return
}

func SaltDecodeRule(src []byte, salt []byte, has bool) (src_, salt_ []byte) {
	switch has {
	case false:
		src_ = make([]byte, len(src)-len(salt))
		copy(src_, src[len(salt):])

		salt_ = src[:len(salt)]
		copy(salt_, src[:len(salt)])
		// src_ = src[salt.len:]
		// salt_ = src[:salt.len]
	default:
		src_ = make([]byte, len(src))
		copy(src_, src)

		salt_ = make([]byte, len(salt))
		copy(salt_, salt)
		// src_ = src
		// salt_ = salt.Salt()
	}

	return
}
