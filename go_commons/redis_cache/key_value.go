package cache

type KVIn struct {
	Key string
	Val interface{}
}

type KVOut struct {
	Key    string
	Val    interface{} // must be pointer
	err    error
	exists bool
}

// SortedSetItem represents an item in a sorted set
type SortedSetItem struct {
	Member string
	Score  float64
}

func (ko *KVOut) reset() {
	ko.err = nil
	ko.exists = false
}

func (ko *KVOut) Err() error {
	return ko.err
}

func (ko *KVOut) Exists() bool {
	return ko.exists
}

func (ko *KVOut) OK() bool {
	return ko.exists && ko.err == nil
}

func makeKVByteOutArr(kArr []string) []*KVOut {
	kvOutArr := make([]*KVOut, 0)
	for _, key := range kArr {
		var temp interface{}
		kvOutArr = append(kvOutArr, &KVOut{Key: key, Val: &temp})
	}

	return kvOutArr
}
