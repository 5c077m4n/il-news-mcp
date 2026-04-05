package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const reutersMiddelEastURL = "https://www.reuters.com/world/middle-east/rss"

func getReutersMiddelEast(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, reutersMiddelEastURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
