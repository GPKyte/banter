package sort // import "github.com/GPKyte/banter/sort"

import (
	"fmt"
	"time"
)

// NumbCollection is my feelings and also these numb-ers
type NumbCollection []int

// Wait right there. Okay keep going.
func Wait(nanos int) {
	time.Sleep(time.Duration(nanos) * time.Nanosecond)
}

// WaitTogether Everybody ready?? Now wait!
func WaitTogether(nanos int, sync <-chan bool) {
	<-sync
	Wait(nanos)
}

// WaitSort this first by leveraging race conditions to sort for you. Can your memory handle the threading?
func WaitSort(this *NumbCollection) *NumbCollection {
	var syncronizeSemaphorSignal = make(chan bool, 1)
	var racerPosition int = 0

	for _, each := range *this {
		go func(Everyone int) {
			// Might be kinda important for everyone to start racing at the same time *Looks directly at ~that~ person* :)
			WaitTogether(Everyone, syncronizeSemaphorSignal)

			(*this)[racerPosition] = Everyone
			racerPosition++
			fmt.Println(Everyone)
		}(each)
	}

	var start = true
	syncronizeSemaphorSignal <- start
	time.Sleep(1 * time.Second)

	return this
}
