package main

import "fmt"

var ch chan int

// https://stackoverflow.com/questions/39107003/select-with-single-case-blocks-adding-default-unblocks
func test() {
	// single

	// multiple
	for {
		select {
		case <-pause:
			fmt.Println("pause")
			select {
			case <-play:
				fmt.Println("play")
			case <-quit:
				wg.Done()
				return
			}
		case <-quit:
			wg.Done()
			return
		default:
			work()
		}
	}

}
