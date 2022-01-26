// go-enum 을 사용해서 열거형 데이터를 만들자
// - go-enum 설치 go-enum.install.sh 파일 실행
// - go generate 실행

//go:generate go-enum --file=listener_type.go --names --nocase=true
package events

/* ENUM(
webhook
file
)
*/
type ListenerType int32
