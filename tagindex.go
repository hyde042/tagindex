package tagindex

import (
	"sort"
	"sync"
)

const bloomFilterK = 9

type Entry struct {
	ID    string
	Order int64
	Tags  []string
}

type entry struct {
	Entry
	tagIDs []uint32
	bloom  bloom
}

type Index struct {
	mu         sync.RWMutex
	data       []entry
	dataIndex  map[string]int
	tagIDs     map[string]uint32
	tagIDCount uint32
}

func New() *Index {
	return &Index{
		dataIndex: make(map[string]int, 1<<8),
		tagIDs:    make(map[string]uint32, 1<<8),
	}
}

type QueryResult struct {
	Data       []Entry
	TotalCount int
}

func (t *Index) Query(tags []string, limit int) QueryResult {
	t.mu.RLock()
	defer t.mu.RUnlock()

	qTagIDs, ok := t.resolveTagIDs(tags, false)
	if !ok {
		return QueryResult{Data: []Entry{}}
	}
	preAlloc := limit
	if preAlloc > 1<<10 {
		preAlloc = 1 << 10
	}
	var (
		qRes   = QueryResult{Data: make([]Entry, 0, preAlloc)}
		qBloom = makeBloom(qTagIDs, bloomFilterK)
	)
	for _, me := range t.data {
		if me.Order <= 0 {
			continue
		}
		if !me.bloom.contains(qBloom) {
			continue
		}
		if !setContains(me.tagIDs, qTagIDs) {
			continue
		}
		qRes.TotalCount++
		if limit > 0 && len(qRes.Data) == limit {
			continue
		}
		qRes.Data = append(qRes.Data, me.Entry)
	}
	return qRes
}

func (t *Index) Put(e ...Entry) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, e := range e {
		var (
			tagIDs, _ = t.resolveTagIDs(e.Tags, true)
			me        = entry{
				Entry:  e,
				tagIDs: tagIDs,
				bloom:  makeBloom(tagIDs, bloomFilterK),
			}
		)
		if i, ok := t.dataIndex[e.ID]; ok {
			t.data[i] = me
		} else {
			t.data = append(t.data, me)
		}
	}
	t.commit()
}

func (t *Index) commit() {
	sort.Slice(t.data, func(i, j int) bool {
		if t.data[i].Order == t.data[j].Order {
			return t.data[i].ID < t.data[j].ID
		}
		return t.data[i].Order > t.data[j].Order
	})
	if len(t.dataIndex) > len(t.data)*2 {
		t.dataIndex = make(map[string]int, len(t.data))
	}
	for i, e := range t.data {
		t.dataIndex[e.ID] = i
	}
}

func (t *Index) resolveTagIDs(tags []string, create bool) ([]uint32, bool) {
	ids := make([]uint32, 0, len(tags))
	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}
		id, ok := t.tagIDs[tag]
		if !ok {
			if !create {
				return nil, false
			}
			t.tagIDCount++
			id = t.tagIDCount
			t.tagIDs[tag] = id
		}
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, true
}
