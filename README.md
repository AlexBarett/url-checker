## URL-checker
CLI-утилита для обработки URl из файла

## How to use

### С использованием makefile

#### запуск

```shell
    # default args
    # input="urls.txt" output="output.txt" timeout=2000 retries=3 connectionLimit={CPU count}
    make run input="urls.txt" output="output.txt" timeout=30
```

#### сборка

```shell
    make build
```

### Без makefile

#### запуск

```shell
    go run cmd/app/main.go  -input="urls.txt" -output="output.txt" -timeout=30 
```

#### сборка

```shell
    go build cmd/app/main.go
```
