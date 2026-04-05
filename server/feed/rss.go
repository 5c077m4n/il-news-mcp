package feed

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

func fetchRSS(ctx context.Context, url string) (*gofeed.Feed, error) {
	parserCtx, parserCancel := context.WithTimeout(ctx, 10*time.Second)
	defer parserCancel()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, parserCtx)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
