package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	l := log.New(os.Stderr, "", 1)
	var ntpUrl string
	flag.StringVar(&ntpUrl, "url", "0.ru.pool.ntp.org", "The url of the remote NTP server")
	flag.Parse()

	time, err := ntp.Time(ntpUrl)
	if err != nil {
		l.Fatal("Error while querying the time from the remote server: ", err)
	}
	fmt.Println("The current time is: ", time)
}
