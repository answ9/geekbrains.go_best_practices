# gb_go_best_practices Homework 1

В консоли запустить, задав настройки при помощи аргументов

```shell
go run cmd/main.go --startUrl https://www.w3.org/Consortium/ --maxDepth 1 --maxErrors 4 --timeOut 15
```

В отдельном окне через консоль отправить команду (убрав лишний слэш перед символом доллара)

```shell
kill -SIGUSR1 $(pgrep main)
```