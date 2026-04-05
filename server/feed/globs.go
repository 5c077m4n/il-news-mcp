package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const globsURL = "https://www.globes.co.il/webservice/rss/rssfeeder.asmx/FeederNode?iID=1725"

func getGlobs(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, globsURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
