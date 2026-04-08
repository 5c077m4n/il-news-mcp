package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

type NewsSource string

const (
	ynetSource              NewsSource = "ynet"
	makoSource              NewsSource = "mako"
	israelHayomSource       NewsSource = "israelHayom"
	wallaSource             NewsSource = "walla"
	maarivSource            NewsSource = "maariv"
	globsSource             NewsSource = "globs"
	jPostSource             NewsSource = "jpost"
	theMarkerSource         NewsSource = "theMarker"
	reutersMiddleEastSource NewsSource = "reutersMiddleEast"
)

var (
	NewsSourceToGetter = map[NewsSource]func(context.Context) (*gofeed.Feed, error){
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
	NewsSourceToOrientation = map[NewsSource]int8{
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
