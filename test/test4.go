package main

import "fmt"

// func main() {

// 	// matrix := [][]int{
// 	// 	{1, 2, 3},
// 	// 	{4, 5, 6},
// 	// 	{7, 8, 9},
// 	// }

// 	for a := 1; a <= 9; a++ {
// 		for b := 1; b <= 9; b++ {
// 			for c := 1; c <= 9; c++ {
// 				for d := 1; d <= 9; d++ {
// 					mmap := map[int]bool{}
// 					if a+b+c+d == 20 {
// 						mmap[a] = true
// 						mmap[b] = true
// 						mmap[c] = true
// 						mmap[d] = true
// 						if len(mmap) == 4 {
// 							fmt.Println("::: ", a, b, c, d)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// }

// 012
// 345
// 678
// 4 - 1 group = 5 group
// 5>2, 5>4,
func main() {
	// do(0, make([]int, 9))
	numbers := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12}
	sum := map[int]int{}
	allresult := [][]int{}

	var add = func(a int) {
		if _, ok := sum[a]; !ok {
			sum[a] = 0
		}
		sum[a] = sum[a] + 1
	}

	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			for k := j + 1; k < len(numbers); k++ {
				for l := k + 1; l < len(numbers); l++ {
					a := numbers[i]
					b := numbers[j]
					c := numbers[k]
					d := numbers[l]
					if a+b+c+d == 24 {
						add(a)
						add(b)
						add(c)
						add(d)
						// fmt.Println(">>>>>>>", a, b, c, d)
						allresult = append(allresult, []int{a, b, c, d})
					}
				}
			}
		}
	}

	// step2
	for _, aaa := range allresult {
		fmt.Println(aaa)
	}

	for i := 0; i < len(allresult); i++ {
		for j := i + 1; j < len(allresult); j++ {
			for k := j + 1; k < len(allresult); k++ {
				for l := k + 1; l < len(allresult); l++ {
					for m := l + 1; m < len(allresult); m++ {
						// check(allresult[i], allresult[j], allresult[k], allresult[l], allresult[m])
						// check := map[int]int{}
						// allresult[i]

					}
				}
			}
		}
	}

	fmt.Println(sum)
}

func check() bool {

}
