# go-enum auto-generator

## install go-enum

- go install github.com/abice/go-enum@v0.4.3

## go-enum 을 사용한 코드 생성

- go generate

    ```text
    //go:generate go run github.com/abice/go-enum --file=[FILE_NAME].go --names --nocase=true
    ```

- cli

    ```sh
    go-enum --file=[FILE_NAME].go --names --nocase=true
    ```
