package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const israelHayomURL = "https://storage.googleapis.com/mako-sitemaps/rssWebSub.xml"

func getIseaelHayom(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, israelHayomURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
