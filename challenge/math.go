package challenge

func factorial(f int) int {
	return factorialByAddition(f)
}

func factorialByAddition(f int) int {
	var acc int = 1
	for f > 2 {
		for i, add := f, acc; i > 0; i-- {
			acc += add
		}
		f--
	}
	return acc
}

func ExampleFactorial() int {
	return factorial(5)
}
