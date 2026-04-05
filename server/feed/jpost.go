package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const jpostURL = "https://www.jpost.com/rss/rssfeedsfrontpage.aspx"

func getJPost(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, jpostURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
