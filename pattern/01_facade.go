package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Ответ:

Данный паттерн применяется, когда имеется сложная и раздробленная система, которую довольно проблематично использовать.
Суть заключается в том, что мы предоставляем унифицированный и простой в использовании интерфейс к данной системе.
Плюсы:
 - можно изолировать сложность системы от пользователя
Минусы:
 - фасад может разрастись до так называемого god object


Небольшая демонстрация на примере получения доступа к веб-странице:

*/

type DNSLookupServer struct{}

func (s *DNSLookupServer) LookupHost(url string) string {
	return "host"
}

type Connection struct{}

func NewConnection(host string) Connection {
	return Connection{}
}

type ResourceGetter struct{}

func (r *ResourceGetter) Get(url string, conn Connection) string {
	return "resource"
}

type Browser struct {
	dns DNSLookupServer
	rg  ResourceGetter
}

func (b *Browser) FetchPage(url string) string {
	host := b.dns.LookupHost(url)
	newConn := NewConnection(host)
	return b.rg.Get(url, newConn)
}

func main() {
	var b Browser
	fmt.Println(b.FetchPage("google.com"))
}
