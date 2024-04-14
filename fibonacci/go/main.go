package main

import (
	"fmt"
	"math"
)

func main() {
	var formula, loop float64
	var isEqual bool
	for n := 1; n <= 100; n++ {
		formula = fibonacciFormula(n)
		loop = fibonacciLoop(n)
		isEqual = formula == loop
		fmt.Println("Nth =", n, ", Formula =", formula, ", Loop =", loop, ", Is equal =", isEqual)
		if !isEqual {
			break
		}
	}
}

// https://www.youtube.com/watch?v=Zo1wu6tO0_g
func fibonacciFormula(n int) float64 {
	var satu float64 = 1
	var akarLima float64 = math.Sqrt(float64(5))
	var dua float64 = 2

	a := (satu + akarLima) / dua
	b := (satu - akarLima) / dua
	c := float64(n)
	result := (math.Pow(a, c) - math.Pow(b, c)) / akarLima
	return math.Round(result)
}

func fibonacciLoop(n int) float64 {
	f := make([]float64, n+2)
	f[0], f[1] = 0, 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}
