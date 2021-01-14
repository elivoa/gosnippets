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

func main() {
	do(0, make([]int, 9))
}

func do(currentLocation int, matrix []int) {
	// fill & judge

	// fmt.Println("-- ", currentLocation, matrix)

	for i := 1; i <= 9; i++ {
		// find if can use
		canuse := true
		for j := 0; j < min(currentLocation, 9); j++ {
			if i == matrix[j] {
				canuse = false
				break
			}
		}
		if canuse {
			matrix[currentLocation] = i

			if currentLocation == 8 {
				// todo judge and return
				// fmt.Println("reach bottom: ", matrix)

				if matrix[0]+matrix[1]+matrix[3]+matrix[4] == 20 &&
					matrix[1]+matrix[2]+matrix[4]+matrix[5] == 20 &&
					matrix[3]+matrix[4]+matrix[6]+matrix[7] == 20 &&
					matrix[4]+matrix[5]+matrix[7]+matrix[8] == 20 {

					fmt.Println("*****: ", matrix)
				}
				return
			} else {
				do(currentLocation+1, matrix)
			}

		}

	}

}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
