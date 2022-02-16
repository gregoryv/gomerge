// my command
package main

import "fmt"

func main() {
	db := &Memory{}
	db.Insert(Car{Brand: "Audi", Model: "A6", Year: 2018})
	db.Insert(Car{Brand: "Volvo", Model: "S90", Year: 2016})

	result := make([]Car, 2)
	db.Select(result)

	for i, car := range result {
		fmt.Printf("%v. %s\n", i+1, car.String())
	}
}


