package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

type newsSource string

const (
	ynetSource              newsSource = "ynet"
	makoSource              newsSource = "mako"
	israelHayomSource       newsSource = "israelHayom"
	wallaSource             newsSource = "walla"
	maarivSource            newsSource = "maariv"
	globsSource             newsSource = "globs"
	jPostSource             newsSource = "jpost"
	theMarkerSource         newsSource = "theMarker"
	reutersMiddleEastSource newsSource = "reutersMiddleEast"
)

var (
	NewsSourceToGetter = map[newsSource]func(context.Context) (*gofeed.Feed, error){
		ynetSource:              getYnet,
		makoSource:              getMako,
		israelHayomSource:       getIseaelHayom,
		wallaSource:             getWalla,
		maarivSource:            getMaariv,
		globsSource:             getGlobs,
		jPostSource:             getJPost,
		theMarkerSource:         getTheMarker,
		reutersMiddleEastSource: getReutersMiddelEast,
	}
	NewsSourceToOrientation = map[newsSource]int8{
		ynetSource:              -5,
		makoSource:              -5,
		israelHayomSource:       5,
		wallaSource:             -3,
		maarivSource:            -5,
		globsSource:             0,
		jPostSource:             -3,
		theMarkerSource:         0,
		reutersMiddleEastSource: -7,
	}
)
