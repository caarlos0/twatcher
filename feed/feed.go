package feed

import (
	"log"
	"strings"

	"github.com/caarlos0/twatcher/torrent"
	rss "github.com/jteeuwen/go-pkg-rss"
)

// Feed type
type Feed struct {
	URL    string
	Filter string
	Names  []string
}

// NewFeed instance
func NewFeed(uri, filter string, names []string) *Feed {
	return &Feed{
		URL:    uri,
		Filter: strings.ToLower(filter),
		Names:  names,
	}
}

// Poll updates the feed and download it to ~/Downloads
func (f *Feed) Poll() {
	feed := rss.New(10, true, f.chanHandler, f.itemHandler)
	log.Println("Looking for new torrents...")
	if err := feed.Fetch(f.URL, nil); err != nil {
		log.Printf("Failed to fetch feed: %s: %s\n", f.URL, err)
	}
}

func (f *Feed) itemHandler(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
	for _, item := range items {
		for _, link := range item.Links {
			f.check(item, link)
		}
	}
}

func (f *Feed) check(item *rss.Item, link *rss.Link) {
	log.Println("Checking", link.Href)
	for _, name := range f.Names {
		if f.matches(link.Href, name) {
			log.Println("Matches", item.Title, link.Href)
			go torrent.NewTorrent(item.Title, link.Href).Download()
		}
	}
}

func (f *Feed) clean(s string) string {
	return strings.Replace(
		strings.Replace(strings.ToLower(s), ".", "", -1), " ", "", -1,
	)
}

func (f *Feed) matches(href, name string) bool {
	href = f.clean(href)
	return strings.Contains(href, f.clean(name)) && strings.Contains(href, f.Filter)

}

func (f *Feed) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	// does nothing
}
