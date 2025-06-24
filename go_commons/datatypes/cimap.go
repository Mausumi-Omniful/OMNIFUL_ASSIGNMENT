package datatypes

import "golang.org/x/text/cases"

// CIMap is a case-insensitive map with generic value type.
type CIMap[V any] struct {
	m map[string]V
}

// NewCIMap creates a new case-insensitive map.
func NewCIMap[V any]() CIMap[V] {
	return CIMap[V]{m: make(map[string]V)}
}

// Set adds or updates a value in the map.
func (cim CIMap[V]) Set(key string, value V) {
	cim.m[cases.Fold().String(key)] = value
}

// Get retrieves a value from the map.
func (cim CIMap[V]) Get(key string) (value V, ok bool) {
	value, ok = cim.m[cases.Fold().String(key)]

	return
}

// Del deletes a key from the map.
func (cim CIMap[V]) Del(key string) {
	delete(cim.m, cases.Fold().String(key))
}

// Keys returns a slice of all keys in the map.
func (cim CIMap[V]) Keys() []string {
	keys := make([]string, 0, len(cim.m))

	for key := range cim.m {
		keys = append(keys, key)
	}

	return keys
}

// Values returns a slice of all values in the map.
func (cim CIMap[V]) Values() []V {
	values := make([]V, 0, len(cim.m))

	for _, value := range cim.m {
		values = append(values, value)
	}

	return values
}

// Length returns a length of the map.
func (cim CIMap[V]) Length() int64 {
	return int64(len(cim.m))
}
