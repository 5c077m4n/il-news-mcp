package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const YnetURL = "https://www.ynet.co.il/Integration/StoryRss2.xml"

type (
	YNetRSSItem struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        string `xml:"link"`
		PubDate     string `xml:"pubDate"`
		Guid        string `xml:"guid"`
		Tags        string `xml:"tags"`
	}
	YnetRSSChannel struct {
		Title         string        `xml:"title"`
		Link          string        `xml:"link"`
		Description   string        `xml:"description"`
		Copyright     string        `xml:"copyright"`
		Language      string        `xml:"language"`
		PubDate       string        `xml:"pubDate"`
		LastBuildDate string        `xml:"lastBuildDate"`
		Items         []YNetRSSItem `xml:"item"`
	}
	YnetRSS struct {
		Channel YnetRSSChannel `xml:"channel"`
	}
)

func newYnetRSS(feed *gofeed.Feed) YnetRSS {
	ynet := YnetRSS{
		Channel: YnetRSSChannel{
			Title:         feed.Title,
			Link:          feed.Link,
			Description:   feed.Description,
			Copyright:     feed.Copyright,
			Language:      feed.Language,
			PubDate:       feed.Published,
			LastBuildDate: feed.Updated,
		},
	}

	for _, item := range feed.Items {
		tags := ""
		if len(item.Categories) > 0 {
			tags = item.Categories[0] // Simple approach, could be improved
		}

		ynet.Channel.Items = append(ynet.Channel.Items, YNetRSSItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     item.Published,
			Guid:        item.GUID,
			Tags:        tags,
		})
	}

	return ynet
}

func GetYnet(ctx context.Context) (*YnetRSS, error) {
	feed, err := fetchRSS(ctx, YnetURL)
	if err != nil {
		return nil, err
	}

	ynetFeed := newYnetRSS(feed)
	return &ynetFeed, nil
}
