package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const maarivURL = "https://www.maariv.co.il/rss/rsschadashot"

func getMaariv(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, maarivURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
