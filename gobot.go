/* HISPAGATOS */
/* By ReK2 */
/* please read the license */

package main

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"github.com/thoj/go-ircevent"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/proxy"
	"log"
	"mvdan.cc/xurls"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	roomName      = "#hispagatos"
	serverName    = "10.8.0.1:6668"
	torServerName = "socks5://10.8.0.1:9050"
	i2pServerName = "socks5://10.8.0.1:4447"
)

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func getURL(site string) *http.Response {

	tbProxyURL, err := url.Parse(torServerName)
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

func getI2pURL(site string) *http.Response {

	tbProxyURL, err := url.Parse(i2pServerName)
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
	resp := scrape.Text(intro)
	return resp
}

func resolveURL(website string) string {
	var site *http.Response

	if strings.Contains(website, ".i2p") {
		site = getI2pURL(website)

	} else {
		site = getURL(website)

	}

	contents, err := html.Parse(site.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
		panic(err)
	}
	title, _ := scrape.Find(contents, scrape.ByTag(atom.Title))
	titulo := scrape.Text(title)
	return titulo

}

func gobot() {

	con := irc.IRC("GoBot", "goBot")
	err := con.Connect(serverName)

	// For Debug
	//con.VerboseCallbackHandler = true
	//con.Debug = true

	if err != nil {
		fmt.Println("Failed connecting")
		fmt.Printf("Err %s", err)
		return
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(roomName)
	})

	con.AddCallback("366", func(e *irc.Event) {})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "!help") {
			output := "Hello Im a Bot my commands are !wiki, !help and I resolve URL's info on channel my owner is ReK2"
			con.Privmsg(roomName, output)
		}
	})

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if strings.Contains(e.Message(), "http") {
			fixed := xurls.Relaxed().FindString(e.Message())
			output := resolveURL(fixed) + " >===> " + fixed
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

func main() {
	cntxt := &daemon.Context{
		PidFileName: "GoBot.pid",
		PidFilePerm: 0644,
		LogFileName: "GoBot.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[GoBot]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("GoBot Started")

	gobot()
}
