package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

const theMarkerURL = "https://www.themarker.com/srv/tm-news"

func getTheMarker(ctx context.Context) (*gofeed.Feed, error) {
	feed, err := fetchRSS(ctx, theMarkerURL)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
