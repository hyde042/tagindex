package tagindex

import (
	"sort"	"strings"
)


type TagInfo struct {
	Tag   string
	Count int
}

func (t *Index) Tags(prefix string, limit int) []TagInfo {
	var res []TagInfo // TODO: smart pre-alloc
	for tag, count := range t.tagCounts {
		if strings.HasPrefix(tag, prefix) {
			res = append(res, TagInfo{Tag: tag, Count: count})
		}
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Count > res[j].Count
	})
	if limit > 0 && len(res) > limit {
		res = res[:limit]
	}
	return res
}
