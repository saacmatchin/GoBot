/* HISPAGATOS */
/* please read the license */

package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	//	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var roomName = "#anarchism"

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

	con.Loop()
}
