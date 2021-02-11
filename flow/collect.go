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

// Collection can be used with filter operations defined here or elsewhere to collect subsets neatly.
type Collection interface {
	Len() int
	Get(this int) interface{}
}

// newGenericCollection is simply a helper method to init a GenericCollection
func newGenericCollection(ofThese ...[]interface{}) GenericCollection {
	var cap int = 0 // no cap, this is a cap, no cap
	for _, each := range ofThese {
		cap += len(each)
	}

	var collection = make([]interface{}, 0, cap)
	for _, each := range ofThese {
		for eachCollectable := range each {
			collection = append(collection, eachCollectable)
		}
	}

	return GenericCollection(collection)
}

// Select of the matching collectable will be returned in one Collection
func Select(c Collection) func(FilterFunc) Collection {
	var this = func(mrfun FilterFunc) Collection {
		return ImmutableMapReduce(mrfun, c)
	}

	return this
}

// ImmutableMapReduce is a variant of a wellknown method
// in effect, given a selective filter, return that which matches 'true'
func ImmutableMapReduce(mrfun FilterFunc, c Collection) Collection {
	var buffer = make([]interface{}, 0, c.Len())
	var countCollectables int = c.Len()

	for thisHere := 0; thisHere < countCollectables; thisHere++ {
		var this = c.Get(thisHere)

		if mrfun(this) {
			buffer = append(buffer, this)
		}
	}

	return newGenericCollection(buffer)
}

// GenericCollection is a simple working example
type GenericCollection []interface{}

// Len over a GenericCollection returns length of underlying slice
func (GC GenericCollection) Len() int {
	return len(GC)
}

// Get this element from the GeneralCollection. Note that default collectable is returned in OOB cases
func (GC GenericCollection) Get(this int) interface{} {
	if this < 0 || this >= GC.Len() {
		// Logging this transgression would be okay behavior.
		// Allow out-of-bound runtime error is not okay
		// Default value returned
		return nil
	}
	return GC[this]
}

// FilterFunc ust be provided with runtime safety implemented by user of library for given types
type FilterFunc func(this interface{}) bool

// Even will return all even integers when used to select elements in Collection
func Even(integer interface{}) bool {
	var whetherEven = func(isUnknown int) bool {
		// bitmasking is by far the fastest method with unsigned integers, but we want to account for more
		// Since I care more about the example's simplicity and less so about efficiency, we will use modulus operation instead
		return (isUnknown%2 == 0)
	}
	switch i := integer.(type) {
	case int, int8, int16, int32, int64:
		return whetherEven(int(i.(int)))
	default:
		return false
	}
}

// Odd this is still here
func Odd(integer interface{}) bool {
	return false // Implement later
}

// Short is a example of something to work on strings
func Short(stick interface{}) bool {
	var whetherShort bool
	const shortie int = 5

	switch stick.(type) {
	case string:
		whetherShort = len(stick.(string)) <= shortie
	default:
		whetherShort = false
	}

	return whetherShort
}

var nrrTracer int = 0

// NotReallyRandom sure looks random...
func NotReallyRandom(aldsjk interface{}) bool {
	nrrTracer += 1 + nrrTracer/23

	if nrrTracer*9%10 == 7 || nrrTracer*2%10 < 3 {
		return true
	}
	return false
}

// All or nothing, just kidding, it's all. it's EVERYTHING.
func All(thing interface{}) bool {
	return true
}

// None shall pass!
func None(thing interface{}) bool {
	return false
}
