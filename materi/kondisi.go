package main

import "fmt"

func main() {
	var hari string = "jumat"

	if hari == "sabtu" || hari == "minggu" {
		fmt.Println("Hari libur kita")
	} else if hari == "jumat" {
		fmt.Println("hari ini kita jumatan")
	} else {
		fmt.Println("hari ini kita kerja yaaa")
	}

	switch hari {
	case "sabtu", "minggu":
		fmt.Println("Hari libur kita")
	case "jumat":
		fmt.Println("hari ini kita jumatan")
	default:
		fmt.Println("hari ini kita kerja yaaa")
	}
}
