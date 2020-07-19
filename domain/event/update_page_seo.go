package event

import (
	"github.com/reecerussell/distro-blog/domain/datamodel"
)

type UpdatePageSEO struct {
	PageID string
	SEO *datamodel.SEO
}
