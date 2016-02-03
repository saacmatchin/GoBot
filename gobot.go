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
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var roomName = "#yourchannelhere"

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func getURL(site string) *http.Response {

	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9150")
	if err != nil {
		fatalf("Failed to parse proxy URL: %v\n", err)
	}

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fatalf("Failed to obtain proxy dialer: %v\n", err)
	}

	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	resp, err := client.Get(site)
	if err != nil {
		fatalf("Failed to issue GET request: %v\n", err)
	}

	fmt.Printf("GET returned: %v\n", resp.Status)
	return resp

}

func queryWikipedia(word string) string {
	word = strings.TrimSpace(word)
	website := "http://en.wikipedia.com/wiki/" + word
	site := getURL(website)
	contents, err := html.Parse(site.Body)

	if err != nil {
		fmt.Print("%s", err)
		panic(err)
		os.Exit(1)
	}
	intro, _ := scrape.Find(contents, scrape.ByTag(atom.P))
	var resp string = scrape.Text(intro)
	return resp
}

func resolveUrl(website string) string {
	site := getURL(website)

	contents, err := html.Parse(site.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
		panic(err)
	}
	title, _ := scrape.Find(contents, scrape.ByTag(atom.Title))
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
