/* HISPAGATOS */
/* please read the license */

package main

import (
	"cgt.name/pkg/go-mwclient"
	"fmt"
	"github.com/thoj/go-ircevent"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"os"
	"strings"
)

var roomName = "#anarchism"

func queryWikipedia(word string) string {

	w, err := mwclient.New("https://en.wikipedia.org/w/api.php", "GoBot.go wikibot")
	if err != nil {
		panic(err)
	}

	resp, time, err := w.GetPageByName(word)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp, time)
	result := time + resp
	return result
}

func resolveUrl(website string) string {
	println(website)
	resp, err := http.Get(website)
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	contents, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
		panic(err)
	}
	title, ok := scrape.Find(contents, scrape.ByTag(atom.Title))
	fmt.Println(scrape.Text(title))
	fmt.Println("%b", ok)
	var titulo string = scrape.Text(title)
	return titulo

}

func main() {
	con := irc.IRC("GoBot", "goBot")
	err := con.Connect("10.8.0.1:6668")
	if err != nil {
		fmt.Println("Failed connecting")
		return
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(roomName)
	})

	con.AddCallback("JOIN", func(e *irc.Event) {
		con.Privmsg(roomName, "Hello! I am prototype of a bot wrote in heavy development")
	})

	/* con.AddCallback("PRIVMSG", func(e *irc.Event) {
		con.Privmsg(roomName, e.Message())
	}) */

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "http") {
			output := resolveUrl(e.Message()) + " => " + e.Message()
			con.Privmsg(roomName, output)
		}
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "!wiki") {
			output := queryWikipedia(e.Message())
			con.Privmsg(roomName, output)
		}
	})
	con.Loop()
}
