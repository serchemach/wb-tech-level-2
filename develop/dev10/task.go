package main

import (
	// "errors"
	"flag"
	"io"
	// "sync"
	// "sync/atomic"

	// "sync/atomic"

	// "fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

var (
	lLog, lErr *log.Logger
)

func init() {
	lErr = log.New(os.Stderr, "[FATAL] ", 1)
	lLog = log.New(os.Stdout, "[LOG] ", 1)
}

func Serve(input io.Reader, output io.Writer, conn net.Conn) {
	buff := make([]byte, 2048)
	for {
		n, err := input.Read(buff)

		if err != nil {
			lLog.Println(err)
			conn.Close()
			break
		}

		if n > 0 {
			output.Write(buff[:n])
		}
	}
}

func main() {
	var (
		timeout     time.Duration
		timeoutText string
		err         error
	)

	flag.StringVar(&timeoutText, "timeout", "10s", "Длительность ожидания подключения, например, 10s или 1h")
	flag.Parse()

	timeout, err = time.ParseDuration(timeoutText)
	if err != nil {
		lErr.Fatal("Ошибка при обработке времени ожидания:", err)
	}

	if len(flag.Args()) < 1 {
		lErr.Fatal("Не предоставлены хост и порт подключения")
	}

	if len(flag.Args()) < 2 {
		lErr.Fatal("Не предоставлен порт подключения")
	}

	host := flag.Args()[0]
	port, err := strconv.Atoi(flag.Args()[1])
	if err != nil {
		lErr.Fatal("Неправильный формат порта:", err)
	}

	if port < 0 {
		lErr.Fatal("Порт не может быть отрицательным")
	}

	conn, err := net.DialTimeout("tcp", host+":"+flag.Args()[1], timeout)
	if err != nil {
		lErr.Fatal("Ошибка при открытии соединения:", err)
	}

	// Поскольку Read на stdin блочится, а закрыть stdin мы не можем, придётся просто смириться с тем,
	// что горутина, читающая из stdin будет умирать по завершению программы (в случае если сервер прислал EOF или соединение умерло)
	go Serve(os.Stdin, conn, conn)
	Serve(conn, os.Stdout, conn)
}
