/* HISPAGATOS */
/* please read the license */

package main

import (
	"fmt"
	"github.com/mvdan/xurls"
	"github.com/thoj/go-ircevent"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"os"
	"strings"
)

var roomName = "#leftsec"

func queryWikipedia(word string) string {
	word = strings.TrimSpace(word)
	website := "http://en.wikipedia.com/wiki/" + word
	//println(website)

	site, err := http.Get(website)
	if err != nil {
		fmt.Println("%s", err)
		panic(err)
		os.Exit(1)
	}
	defer site.Body.Close()
	contents, err := html.Parse(site.Body)
	if err != nil {
		fmt.Print("%s", err)
		panic(err)
		os.Exit(1)
	}
	intro, _ := scrape.Find(contents, scrape.ByTag(atom.P))
	//fmt.Println(scrape.Text(intro))
	//fmt.Println("%b", ok)
	var resp string = scrape.Text(intro)
	return resp
}

func resolveUrl(website string) string {
	//println(website)
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
	title, _ := scrape.Find(contents, scrape.ByTag(atom.Title))
	//fmt.Println(scrape.Text(title))
	//fmt.Println("%b", ok)
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

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "!help") {
			output := "Hello Im a Bot, my commands are !wiki, !help and I resolve URL's info on channel my owner is NetAnarchist"
			con.Privmsg(roomName, output)
		}
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "http") {
			fixed := xurls.Relaxed.FindString(e.Message())
			output := resolveUrl(fixed) + " >===> " + fixed
			con.Privmsg(roomName, output)
		}
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "!wiki") {
			fixed := strings.Replace(e.Message(), "!wiki", "", -1)
			output := queryWikipedia(fixed)
			con.Privmsgf(roomName, output)
		}
	})
	con.Loop()
}
