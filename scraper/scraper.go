package scraper

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/mucanyu/eksisozluk-go/model"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	mainURL   = "https://www.eksisozluk.com/"
	gundemURL = mainURL + "basliklar/gundem"

	grey        = color.New(color.FgBlack, color.FgWhite)
	baslikColor = color.New(color.Bold, color.FgHiGreen)
	hiYellow    = color.New(color.FgHiYellow)
	hiWhite     = color.New(color.FgHiWhite)

	// declare matchers
	entryListMatcher, dateMatcher, entryMatcher, authorMatcher func(n *html.Node) bool
	maxGundemNumMatcher, topicMatcher, isTopicFoundMatcher     func(n *html.Node) bool
)

func init() {
	topicMatcher = func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent, "class") == "topic-list"
		}
		return false
	}

	isTopicFoundMatcher = func(n *html.Node) bool {
		if scrape.Attr(n, "id") == "topic" {
			for _, atr := range n.Attr {
				if atr.Key == "data-not-found" {
					return false
				}
			}
			return true
		}
		return false
	}

	entryMatcher = scrape.ByClass("content")
	dateMatcher = scrape.ByClass("entry-date")
	authorMatcher = scrape.ByClass("entry-author")
	maxGundemNumMatcher = scrape.ByClass("topic-list-description")
	entryListMatcher = scrape.ById("entry-item-list")
}

// PrintGundem prints popular topics
func PrintGundem(limit, pageVal int) error {
	t, err := popularTopics(limit, pageVal)
	if err != nil {
		return err
	}

	hiWhite.Printf("%3s\n", "#")
	for i, topic := range t {
		hiYellow.Printf("%3d", i+1)
		hiWhite.Printf(" %s ", topic.Title)
		grey.Printf("(%s)\n", topic.NewEntryCount)
	}
	return nil
}

// popularTopics scrapes the link and then
// returns a []models.Topic which contains title, new entry count and ID.
func popularTopics(lim, pageVal int) ([]model.Topic, error) {
	topicList := make([]model.Topic, 0)

	for {
		resp, err := http.Get(gundemURL + "?p=" + strconv.Itoa(pageVal))
		if err != nil {
			return nil, errors.New("ERROR: internet bağlantınızı kontrol edin")
		}

		root, err := html.Parse(resp.Body)
		if err != nil {
			return nil, errors.New("ERROR: An error occured while parsing body")
		}
		defer resp.Body.Close()

		topics := scrape.FindAll(root, topicMatcher)
		if len(topics) == 0 {
			return nil, errors.New("ERROR: Lütfen parametre değerlerinizi kontrol edin")
		}
		for _, topic := range topics {
			t := model.Topic{}

			t.Title = scrape.Text(topic.FirstChild)
			t.NewEntryCount = scrape.Text(topic.LastChild)

			topicList = append(topicList, t)

			if len(topicList) == lim {
				return topicList, nil
			}
		}
		pageVal++
	}
}

// PrintTopic prints entries in a topic
func PrintTopic(params *model.BaslikParams) error {
	var redirectedURL string

	// Get the redirect url instead of page content
	req, err := http.NewRequest("GET", mainURL+"?q="+url.QueryEscape(params.Topic), nil)
	if err != nil {
		return errors.New("ERROR: internet bağlantınızı kontrol edin")
	}
	client := new(http.Client)

	// It is called before following any redirection
	// If this function returns an error, then httpClient.Do(...) will not follow the redirect and instead will return an error
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) == 2 {
			redirectedURL = req.URL.String()
			return errors.New("stopping redirect")
		}
		return nil
	}

	response, err := client.Do(req)
	if _, err := response.Location(); err != nil {
		return errors.New("ERROR: Böyle bir başlık yok")
	}
	response.Body.Close()

	if params.Sukela {
		redirectedURL += "?a=nice"
	}

	if strings.Contains(redirectedURL, "?") {
		redirectedURL += "&p="
	} else {
		redirectedURL += "?p="
	}

	entryList, err := getEntries(redirectedURL, params.Page, params.Limit)
	if err != nil {
		return err
	}

	baslikColor.Printf("%s ", params.Topic)
	grey.Printf("[%d. sayfa]\n\n", params.Page)
	for i, entry := range entryList {
		hiYellow.Printf("%d. ", i+1)
		hiWhite.Printf("%s \n", entry.Text)
		color.Green("[%s | %s]", entry.Author, entry.Date)
	}
	return nil
}

// getEnries returns either specific count of entries or all of them in page(s)
func getEntries(URL string, page, limit int) ([]model.Entry, error) {
	entries := make([]model.Entry, 0)

	for {
		resp, err := http.Get(URL + strconv.Itoa(page))
		if err != nil {
			return nil, errors.New("ERROR: internet bağlantınızı kontrol edin")
		}

		root, err := html.Parse(resp.Body)
		if err != nil {
			return nil, errors.New("ERROR: An error occured while parsing body")
		}
		defer resp.Body.Close()

		if _, ok := scrape.Find(root, isTopicFoundMatcher); !ok {
			return nil, errors.New("ERROR: " + strconv.Itoa(page) + ". sayfa yok")
		}

		entryList, ok := scrape.Find(root, entryListMatcher)
		if !ok {
			return nil, errors.New("ERROR: Site layoutu değişmiş olabilir | Report: github.com/mucanyu/eksisozluk-go/issues")
		}

		entryNodeList := scrape.FindAll(entryList, entryMatcher)

		for _, enode := range entryNodeList {
			entry := model.Entry{}
			entry.Text = scrape.Text(enode)

			autNode, ok := scrape.Find(enode.Parent, authorMatcher)
			if ok {
				entry.Author = scrape.Text(autNode)
			}

			dateNode, ok := scrape.Find(enode.Parent, dateMatcher)
			if ok {
				entry.Date = scrape.Text(dateNode)
			}

			entries = append(entries, entry)
			limit--
			if limit == 0 {
				return entries, nil
			}
		}
		if len(entries)%10 != 0 {
			return entries, nil
		}
		page++
	}
}
