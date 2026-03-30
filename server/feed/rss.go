package feed

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

func fetchRSS(url string) (*gofeed.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
