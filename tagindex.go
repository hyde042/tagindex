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
	Entry     // TODO: avoid storing the full entity with all the tags etc.
	tagIDs    []uint32
	bloom     Bloom
	isDeleted bool
}

type Index struct {
	mu           sync.RWMutex
	data         []entry
	dataIndex    map[string]int
	tagIDs       map[string]uint32
	tagCounts    map[string]int // TODO: more optimal tag store (prefix map?)
	isDirty      bool
	orderCounter int64
}

func New() *Index {
	return &Index{
		dataIndex: make(map[string]int, 1<<8),
		tagIDs:    make(map[string]uint32, 1<<8),
		tagCounts: make(map[string]int, 1<<8),
	}
}

type QueryResult struct {
	Data       []Entry
	TotalCount int
}

func (t *Index) Query(tags []string, limit int) QueryResult {
	t.Commit()

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
		qBloom = MakeBloom(qTagIDs, bloomFilterK)
	)
	for _, me := range t.data {
		if me.isDeleted {
			continue
		}
		if !me.bloom.Contains(qBloom) {
			continue
		}
		if !SetContains(me.tagIDs, qTagIDs) {
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

		// TODO: compare to existing entry to avoid pointless commits

		var (
			tagIDs, _ = t.resolveTagIDs(e.Tags, true)
			me        = entry{
				Entry:  e,
				tagIDs: tagIDs,
				bloom:  MakeBloom(tagIDs, bloomFilterK),
			}
		)
		if i, ok := t.dataIndex[e.ID]; ok {
			if me.Order <= 0 {
				me.isDeleted = me.Order < 0
				me.Order = t.data[i].Order
			}

			// TODO: update tag counts

			t.data[i] = me
		} else {
			if me.Order == 0 {
				t.orderCounter++
				me.Order = -t.orderCounter
			}
			for _, tag := range me.Tags {
				t.tagCounts[tag]++
			}
			t.dataIndex[e.ID] = len(t.data)
			t.data = append(t.data, me)
		}
	}
	t.isDirty = true
}

func (t *Index) Commit() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isDirty {
		return
	}
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
	t.isDirty = false
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
			id = uint32(len(t.tagIDs) + 1)
			t.tagIDs[tag] = id
		}
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, true
}
