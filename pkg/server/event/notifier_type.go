// go-enum 을 사용해서 열거형 데이터를 만들자
// - go-enum 설치 go-enum.install.sh 파일 실행
// - go generate 실행

//go:generate go-enum --file=notifier_type.go --names --nocase=true
package event

/* ENUM(
console
webhook
file
rabbitMQ
)
*/
type NotifierType int32
