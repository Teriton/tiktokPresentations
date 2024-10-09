package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

func main() {
	downloadFiles(
		"Здравствуй, наконец-то мы тебя нашли! Это администрация порно сайта который ты неоднократно посещаешь.Мы решили напомнить, что срок твоей годовой подписки истекает сегодня и ее необходимо продлить.Также по твоим многочисленным просьбам мы обновили раздел Гей-Порно с кучей увлекательных сюжетов, которые обязательно придутся тебе по вкусу. Не забудь поделиться этой информацией со своей семьей, друзьями и знакомыми.Ждем Всех в следующий раз!",
		5,
	)
	// InitGtk()
}

func downloadFiles(str string, deviedeFactor int) {
	strArr := parseString(str, deviedeFactor)

	for i, j := range strArr {
		print(j, "(", i, ")")
	}

	for numOfStr, word := range strArr {
		links, err := getListOfLinks(url.QueryEscape(word))

		if err != nil {
			log.Fatal("GG")
			return
		}

		for _, urlOfImage := range links {
			err = downloadFile(urlOfImage, fmt.Sprintf("./buf/%d.jpg", numOfStr))
			if err != nil {
				continue
			} else {
				break
			}
		}

	}
}

func parseString(str string, groupFactor int) []string {
	patt := regexp.MustCompile(`[-,_+,',–]`)
	str = patt.ReplaceAllString(str, "")
	pattern := regexp.MustCompile(fmt.Sprintf("((\\w|[\\p{Cyrillic}])+.){0,%d}", groupFactor))
	allSubstringMatches := pattern.FindAllString(str, -1)

	return allSubstringMatches
}

// func InitGtk() {
// 	gtk.Init(&os.Args)
// 	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
// 	window.SetDefaultSize(700, 300)
// 	vbox := gtk.NewVBox(false, 1)
// 	image := gtk.NewImageFromFile("file.jpg")
// 	vbox.Add(image)
// 	window.Add(vbox)
// 	window.ShowAll()
// 	gtk.Main()
// }

func getListOfLinks(query string) (links []string, err error) {
	c := colly.NewCollector()

	// OnHTML callback
	c.OnHTML("img", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("src"))
	})
	c.Visit(fmt.Sprintf("https://www.google.com/search?q=%s&udm=2", query))
	// c.Visit(fmt.Sprintf("https://yandex.by/images/search?from=tabbar&text=%s", query))
	return
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getHtmlPage(webPage string) (string, error) {

	resp, err := http.Get(webPage)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {

		return "", err
	}

	return string(body), nil
}
