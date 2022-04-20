// go-enum 을 사용해서 열거형 데이터를 만들자
// - go-enum 설치 go-enum.install.sh 파일 실행
// - go generate 실행
//go:generate go-enum --file=encryption_method.go --names --nocase
package enigma

/* ENUM(
	NONE
	AES
	DES
)
*/
type EncryptionMethod int

type NoneEncripter struct {
	// key []byte
}

func (encripter NoneEncripter) BlockSize() int {
	// return len(encripter.key)
	return 1
}

func (encripter NoneEncripter) Encrypt(dst, src []byte) {
	copy(dst, src)
}

func (encripter NoneEncripter) Decrypt(dst, src []byte) {
	copy(dst, src)
}
