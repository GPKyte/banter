package sort // import "github.com/GPKyte/banter/sort"

type Member int

func SelectCommonMembers(numLists [][]int) []int {
	var memberLists = convertNumsToMembers(numLists)
	// Once all memberships verified for a shared member, retain it.
	// Basic version is a count

	// Or run ReadFirstToN on the buffered input of the shared member series
	// FirstTo(Len(memberLists))

	var maxMemberships int = len(memberLists) // Once reached by MMMM, item is in every collection
	var trackMembership = map[Member]int{}    // Must be initialized for first reference inside of trackAssistant
	var trackAssistant = func(forThis Member) {
		var someCount int = trackMembership[forThis] // Will default to zero if not initialized
		trackMembership[forThis] = someCount + 1
	}

	var countUpMemberships = applyFuncToAll
	countUpMemberships(memberLists, trackAssistant)

	var commonMembers = make([]Member, 0) // Capacity...Max of the memberlists
	var saveAllCommonMembers = applyFuncToKeysWithMatchingValues
	var saveAnyQualifiedMember = func(m Member, count int) {
		if count == maxMemberships {
			commonMembers = append(commonMembers, m)
		}
	}

	saveAllCommonMembers(trackMembership, saveAnyQualifiedMember)
	return convertMembersToNums([][]Member{commonMembers})[0]
}

func convertNumsToMembers(src [][]int) [][]Member {
	var conversion = make([][]Member, len(src))

	for j := range src {
		conversion[j] = make([]Member, 0, len(src[j]))

		for k := range src[j] {
			conversion[j] = append(conversion[j], Member(src[j][k]))
		}
	}
	return conversion
}

func convertMembersToNums(src [][]Member) [][]int {
	var conversion = make([][]int, len(src))

	for j := range src {
		conversion[j] = make([]int, 0, len(src[j]))

		for k := range src[j] {
			conversion[j] = append(conversion[j], int(src[j][k]))
		}
	}
	return conversion
}

func applyFuncToKeysWithMatchingValues(namesAreHard map[Member]int, do func(Member, int)) {
	for k, v := range namesAreHard {
		do(k, v)
	}
}

func applyFuncToAll(sliceOfSlices [][]Member, do func(Member)) {
	for _, slice := range sliceOfSlices {
		for _, thing := range slice {
			do(thing)
		}
	}
}

func ReadFirstToN(collection []int, N int) chan int {
	var lookup = map[int]int{}
	var firstToN = make(chan int, len(collection))
	var countToN = func(countThis int) bool {
		count := lookup[countThis]
		lookup[countThis]++

		if count != N {
			return false
		} else {
			return true
		}
	}
	var countSort = func() {
		for _, c := range collection {
			if countToN(c) {
				firstToN <- c
			}
		}
		defer close(firstToN)
	}
	go countSort()
	return firstToN
}
