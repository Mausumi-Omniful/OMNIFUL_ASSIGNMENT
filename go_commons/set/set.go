package set

/* ######### EXAMPLE ##############
	INTEGER:
itrInt := func(val int) { fmt.Println(val) }
set1 := set.New[int]()
set2 := set.New[int]()
set1.Add(1, 2, 3, 4, 9, 11, 2, 23, 3, 2, 4)
fmt.Println(set1.Has(9))  //returns true
set1.Delete(1, 2, 3)     //deletes from and returns nothing
set1.Iterate(itrInt)
set2.Add(3, 4, 5, 6, 7, 9)
set1.Union(set2).Iterate(itrInt)        //returns union set
set1.Intersection(set2).Iterate(itrInt) //returns intersection set
	STRING:
itrString := func(val string) { fmt.Println(val) }
set4 := set.New[string]()
letters := []string{"a", "b", "c", "d", "e", "f", "g", "a", "a"}
set4.AddArray(letters)
set4.Iterate(itrString)
#################################### */

type Set[T comparable] map[T]bool // Declaring new data type
// T is the generic type it can be any comparable data type

// New is a Constructor to create new set
// Example :-  New(int)() to create a int set
// New(string)() to create a string set
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Add : Adds values to set
func (s Set[T]) Add(values ...T) {
	for _, value := range values {
		s[value] = true
	}
}

// AddArray : Adds an array of type T to set
func (s Set[T]) AddArray(values []T) {
	for _, value := range values {
		s[value] = true
	}
}

// Delete : Deletes values from set and return nothing
func (s Set[T]) Delete(values ...T) {
	for _, value := range values {
		delete(s, value)
	}
}

// Len : Returns the length of the set
func (s Set[T]) Len() int {
	return len(s)
}

// Has : Returns the presence of value in the set
func (s Set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}

// Iterate : Iterate the set and perform the given function over the set
func (s Set[T]) Iterate(it func(T)) {
	for v := range s {
		it(v)
	}
}

// Values : Convert set to slice of values
func (s Set[T]) Values() []T {
	values := make([]T, 0)
	s.Iterate(func(value T) {
		values = append(values, value)
	})
	return values
}

// Clone : Returns a clone of the set
func (s Set[T]) Clone() Set[T] {
	set := make(Set[T])
	set.Add(s.Values()...)
	return set
}

// Union of 2 sets
func (s Set[T]) Union(other Set[T]) Set[T] {
	set := s.Clone()
	set.Add(other.Values()...)
	return set
}

// Intersection of 2 sets
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	set := make(Set[T])
	s.Iterate(func(value T) {
		if other.Has(value) {
			set.Add(value)
		}
	})
	return set
}
