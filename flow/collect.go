package flow // import "github.com/GPKyte/banter/flow"

// I envision a lambda style dotnotation usage such as
// bars := CreateCollection(somethings...)
// func whichAreGood(thing) bool {...}
// goodBars := Any(bars)(whichAreGood)
//
// Why is this useful? Not sure but I think of one of these all the time
// But I get stuck on typing, usually there is a logical flaw
// Where the type ambiguity is at fault.
// Often it boils down to the interface{} and type checking would happen in the whichAreGood function


type Collection interface {
	Len() int
	Get(this int) collectable
}

type collectable interface {}

func NewGenericCollection(ofThese... []interface{}) GenericCollection {
	var cap int = 0 // no cap, this is a cap, no cap
	for each := range ofThese {
		for eachCollectable := range each {
			cap += 1
		}
	}

	var collection = make([]collectable, 0, cap)
	for each := range ofThese {
		for eachCollectable := range each {
			collection = append(collection, collectable(eachCollectable))
		}
	}

	return GeneralCollection{collection}
}

// Select of the matching collectable will be returned in one Collection
func Select(c *Collection) (func(MapReduceFunc) Collection) {
	return ImmutableMapReduce
}

// ImmutableMapReduce is a variant of a wellknown method
// in effect, given a selective filter, return that which matches 'true'
func (c Collection) ImmutableMapReduce(mrfun FilterFunc) Collection {
	var reduction Collection
	var countCollectables int = c.Len()
	var thisHere int = 0

	for thisHere < countCollectables; thisHere++ {
		var this collectable = c.Get(thisHere)
		if mrfun() {
			buffer = append(buffer, this)
		}
	}

	return c.(type){buffer}
}

type GenericCollection []collectable

// Len over a GenericCollection returns length of underlying slice
func (GC GenericCollection) Len() int {
	return len(([]interface{})GC)
}

// Get this element from the GeneralCollection. Note that default collectable is returned in OOB cases
func (GC GenericCollection) Get(this int) collectable {
	if this < 0 || this >= GC.Len() {
		// Logging this transgression would be okay behavior.
		// Allow out-of-bound runtime error is not okay
		// Default value returned
		return collectable{}
	}
	return GC[this]
}

// Must be provided with runtime safety implemented by user of library for given types
type FilterFunc func(this interface{}) bool

type Even(integer interface{}) bool {
	var whetherEven = func(isUnknown int) bool {
		// bitmasking is by far the fastest method with unsigned integers, but we want to account for more
		// Since I care more about the example's simplicity and less so about efficiency, we will use modulus operation instead
		return (isUnknown % 2 == 0)
	}
	switch integer.(type) {
	case int, int8, int16, int32, int64:
		return whetherEven(int(integer))
	default:
		return false
	}
}

type Odd(integer interface{}) bool {
	return false // Implement later
}

type Short(stick interface{}) bool {
	var whetherShort bool
	const shortie int = 5

	switch stick.(type) {
	case string:
		whetherShort = len(stick) <= shortie
	default:
		whetherShort = false
	}

	return whetherShort
}

var nrrTracer int = 0
type NotReallyRandom(aldsjk interface{}) bool {
	var trace = func() {
		nrrTracer += 1 + nrrTracer / 23
	}

	if nrrTracer * 9 % 10 == 7 || nrrTracer * 2 % 10 < 3 {
		return true
	} else {
		return false
	}
}

type All(thing interface{}) bool {
	return true
}

type None(thing interface{}) bool {
	return false
}
