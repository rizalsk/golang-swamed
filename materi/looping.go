package main

import "fmt"

func main() {

	for i := 0; i < 5; i++ {
		for j := i; j < 5; j++ {
			fmt.Print(j, " ")
		}

		fmt.Println()
	}

	var ys = [5]int{10, 20, 30, 40, 50} // array
	for _, v := range ys {
		fmt.Println("Value=", v)
	}

hentikan:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i == 3 {
				break hentikan //hentikan di label teratas
			}
			fmt.Print("matriks i[", i, "]:j[", j, "]", "\n")
		}
	}

}
