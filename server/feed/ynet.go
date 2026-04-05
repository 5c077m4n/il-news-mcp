package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const ynetURL = "https://www.ynet.co.il/Integration/StoryRss2.xml"

func getYnet(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, ynetURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
