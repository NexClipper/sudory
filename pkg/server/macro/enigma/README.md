# enigma

## config file

- enigma.yml

    ```yaml
    enigma:
    some-config:
        block-method: AES
        block-size: 128
        block-key: YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n # base64("brown fox jumps over the lazy dog")
        cipher-mode: GCM
        cipher-salt: c2FsdHk= # base64("salty")
        padding: NONE
        strconv: base64
    ```

## usage

```go
func main() {
    filename := "enigma.yml"

    config := enigma.Config{}
    if err := configor.Load(&config, filename); err != nil {
        panic(err)
    }
    if err := enigma.LoadConfig(config); err != nil {
        panic(err)
    }
}

func() {
    const some_enigma_config_name = "some-config"
    const example = "brown fox jumps over the lazy dog"
    
    encoded, _ := enigma.GetMachine(some_enigma_config_name).Encode([]byte(example))
    fmt.Println("encoded:", string(encoded)) 
    
    plain, _ := enigma.GetMachine(some_enigma_config_name).Decode(encoded)
    fmt.Println("plain:", string(plain)) 

    if strings.EqualFold(example, plain) {
        panic("diff")
    }
}
```

```text
output
encoded: nuZFj5ox3VvUUh12bgAnLdmo+Jg1cKZSonlY4nm7KBajBAz6w+yyR6VXDJdXfwxxwA==
plain: brown fox jumps over the lazy dog
```

## table of config

- block

    | block-method | block-size     | block-key      |
    |---           |---             |---             |
    | NONE         | default(1)     | base64(string) |
    | *AES         | 128, 192, 256  | base64(string) |
    | DES          | 64             | base64(string) |

- cipher

    | cipher-mode | block-method   |
    |---          |---             |
    | NONE        | NONE, AES, DES |
    | *GCM        | AES            |
    | CBC         | NONE, AES, DES |

  - cipher-salt: [base64(string) | (null)]

    지정되지 않은 경우 암호화 하면서 생성한 salt값을 암호화 결과의 앞에 붙여서 리턴 한다

    복호화 하면서 앞에 붙어있는 salt값을 분리하여 복호화에서 사용

- padding

    | padding | block-method+cipher-mode                          |
    |---      |---                                                |
    | *NONE   | AES+GCM                                           |
    | PKCS    | AES+NONE, AES+GCM, AES+CBC,<br/>DES+NONE, DES+CBC |

- strconv

    | strconv |
    |---      |
    | plain   |
    | *base64 |
    | hex     |

## 설정 순서

1. 암호화 블럭을 만든다 (block-method, block-size, block-key)

1. 암호화 cipher를 만든다 (cipher-mode)

    - GCM: 나중에 nonce를 생성하기 위해서 cipher.AEAD.NonceSize() 값을 저장

    - CBC: 나중에 iv를 생성하기 위해서 cipher.Block.BlockSize() 값을 저장

1. 암호화 블럭과 cipher를 이용하여 enigma.Machine에서 이용하는 Encoder, Decoder 함수를 생성하여 enigma.Machine 생성

## 암호화 순서

```go
func (machine *Machine) Encode(src []byte) ([]byte, error)
```

1. 암화와 블럭 사이즈 만큼 입력값에 패드 추가

1. Encoder 실행

1. salt encode rule 적용; salt 값이 null이면 암호화 결과에 salt를 앞에 붙이는 작업

1. strconv encode; 지정된 변환 설정에 따라 []byte 결과를 인코드 한다

## 복호화 순서 `(조립의 역순)`

```go
func (machine *Machine) Decode(src []byte) ([]byte, error)
```

1. strconv decode; 지정된 변환 설정에 따라 []byte 결과를 디코드 한다

1. salt decode rule 적용; salt 값이 null이면 암호화 결과에서 앞에 저장된 salt를 분리하는 작업

1. Decoder 실행

1. 암화와 블럭 사이즈 만큼 입력값에 패드 제거
