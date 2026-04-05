package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const wallaURL = "https://rss.walla.co.il/feed/1?type=main"

func getWalla(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, wallaURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
