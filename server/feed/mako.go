package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const makoURL = "https://storage.googleapis.com/mako-sitemaps/rssWebSub.xml"

func getMako(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, makoURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
