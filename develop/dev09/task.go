package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	lLog, lErr *log.Logger
	seenUrls   map[string]struct{}
)

func init() {
	lErr = log.New(os.Stderr, "[FATAL] ", 1)
	lLog = log.New(os.Stdout, "[LOG] ", 1)
	seenUrls = make(map[string]struct{})
}

func SavePage(link *url.URL, page []byte) {
	filename := link.Host + link.Path
	if filename[len(filename)-1] == '/' {
		filename += "index.html"
	}

	var err error
	ind := strings.LastIndex(filename, "/")
	if ind != -1 {
		err = os.MkdirAll(filename[:ind+1], 0700)
		if err != nil {
			lLog.Println("Ошибка при создании папки для", filename, err)
			return
		}
	}

	err = os.WriteFile(filename, page, 0644)
	if err != nil {
		lLog.Println("Ошибка при сохранении страницы", link.String(), err)
		return
	}

	lLog.Println("Сохранена страница", link.String(), "по пути", filename)
}

func GetNewSiteLinks(baseLink *url.URL, page []byte) []*url.URL {
	z := html.NewTokenizer(bytes.NewReader(page))
	links := make([]*url.URL, 0)

	for {
		if z.Next() == html.ErrorToken {
			break
		}

		curToken := z.Token()

		if curToken.Type != html.StartTagToken || curToken.Data != "a" {
			continue
		}

		for _, attr := range curToken.Attr {
			possibleUrl, err := url.Parse(attr.Val)
			if err != nil {
				continue
			}

			*possibleUrl = url.URL{
				Host:   possibleUrl.Host,
				Scheme: possibleUrl.Scheme,
				Path:   possibleUrl.Path,
			}

			// Фиксим относительные ссылки
			if possibleUrl.Host == "" {
				possibleUrl = baseLink.ResolveReference(possibleUrl)
			}
			// lLog.Println(possibleUrl)

			if _, ok := seenUrls[possibleUrl.String()]; ok {
				continue
			}
			seenUrls[possibleUrl.String()] = struct{}{}

			if possibleUrl.Host == baseLink.Host && strings.HasPrefix(possibleUrl.Path, baseLink.Path) {
				if len(links) == cap(links) {
					links = append(links, make([]*url.URL, len(links))...)[:len(links)]
				}
				links = append(links, possibleUrl)
			}
		}

	}

	return links
}

func WalkDown(depth int, link *url.URL, baseLink *url.URL, client *http.Client) {
	if depth == 0 {
		return
	}

	*link = url.URL{
		Host:   link.Host,
		Scheme: link.Scheme,
		Path:   link.Path,
	}

	if link.Scheme == "" {
		link.Scheme = "http"
	}

	resp, err := client.Get(link.String())
	if err != nil {
		lLog.Println("Не получилось скачать", link.String(), err)
		return
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		lLog.Println("Не получилось прочитать данные с", link.String(), err)
		return
	}

	SavePage(link, respData)

	for _, newLink := range GetNewSiteLinks(baseLink, respData) {
		WalkDown(depth-1, newLink, baseLink, client)
	}
}

func main() {
	var walkDepth int
	flag.IntVar(&walkDepth, "depth", 1, "Глубина скачивания сайта. При значении -1 пытается скачать все возможные подстраницы.")
	flag.Parse()

	if len(flag.Args()) < 1 {
		lErr.Fatal("Не предоставленно ссылок.")
	}

	client := &http.Client{}

	for _, link := range flag.Args() {
		url, err := url.Parse(link)
		if err != nil {
			lLog.Println("Некорректная ссылка", url, ", пропускается")
			continue
		}

		WalkDown(walkDepth, url, url, client)
		seenUrls = make(map[string]struct{})
	}

}
