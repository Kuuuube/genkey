package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"runtime"
	"sort"
)

func Score(l string) float64 {
	var score float64
	speeds := CalcFingerSpeed(l)

	weightedSpeed := 0.00

	leftIndex, rightIndex := CalcIndexUsage(l)
	score += float64(25*(13-leftIndex))

	score += float64(25*(13-rightIndex))

	lowest := 100000000000.0
	highest := 0.0
	for i, speed := range speeds {
		weightedSpeed += speed*speed / (KPS[i]*KPS[i])
		if weightedSpeed < lowest {
			lowest = weightedSpeed
		}
		if weightedSpeed > highest {
			highest = weightedSpeed
		}
	}

	deviation := highest - lowest

	score += weightedSpeed
	score += 0.3*deviation
	return score/1000000
}

func randomLayout() string {
	chars := "abcdefghijklmnopqrstuvwxyz,./'"
	length := len(chars)
	l := ""
	for i := 0; i < length; i++ {
		char := string([]rune(chars)[rand.Intn(len(chars))])
		l += char
		chars = strings.ReplaceAll(chars, char, "")
	}
	return l
}

type Layout struct {
	Keys  string
	Score float64
}

func Populate(n int) []string {
	rand.Seed(time.Now().Unix())
	layouts := []string{}
	for i := 0; i < n; i++ {
		layouts = append(layouts, randomLayout())
		fmt.Printf("%d random created...\r", i+1)

	}
	fmt.Println()
	for i, _ := range layouts {
		go greedyImprove(&layouts[i])
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d improving...\r", runtime.NumGoroutine()-1)
	}
	fmt.Println()

	fmt.Println("Sorting...")
	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i]) < Score(layouts[j]) 
	})
	PrintLayout(layouts[0])
	fmt.Println(Score(layouts[0]))
	PrintLayout(layouts[0])
	fmt.Println(Score(layouts[1]))
	PrintLayout(layouts[0])
	fmt.Println(Score(layouts[2]))

	layouts = layouts[0:49]

	for i, _ := range layouts {
		go fullImprove(&layouts[i])
	}
	for runtime.NumGoroutine() > 1 {
		fmt.Printf("%d fully improving...\r", runtime.NumGoroutine()-1)
	}

	sort.Slice(layouts, func(i, j int) bool {
		return Score(layouts[i]) < Score(layouts[j]) 
	})
	
	fmt.Println()
	PrintLayout(layouts[0])
	fmt.Println(Score(layouts[0]))
	fmt.Println(CalcIndexUsage(layouts[0]))
	PrintLayout(layouts[1])
	fmt.Println(Score(layouts[1]))
	PrintLayout(layouts[2])
	fmt.Println(Score(layouts[2]))
	return layouts
}

func greedyImprove(layout *string)  {
	stuck := 0
	for {
		stuck++
		prop := cycleRandKeys(*layout, 1)

		first := Score(*layout)
		second := Score(prop)

		if second < first {
			// accept
			*layout = prop
			stuck = 0
		} else {
			stuck++
		}

		if stuck > 100 {
			return
		}
	}

}

func fullImprove(layout *string) {
	i := 0
	tier := 1
	i = 0
	changed := false
	for {
		i += 1
		prop := cycleRandKeys(*layout, tier)
		first := Score(*layout)
		second := Score(prop)

		if second < first {
			*layout = prop
			i = 0
			changed = true
			continue
		} else if second == first {
			*layout = prop
		} else {
			if i > 5000*tier {
				if changed {
					tier = 1
				} else {
					tier++
				}

				changed = false

				if tier > 6 {
					break
				}

				i = 0
			}
		}
		continue
	}

}

func Anneal(l *string) {
	for temp:=100;temp>-5;temp-- {
		for i:=0;i<300;i++ {
			prop := cycleRandKeys(*l, 1)
			first := Score(*l)
			second := Score(prop)
			if second < first || rand.Intn(100) < temp {
				*l = prop
			} 
			
		}
	}
}

func cycleRandKeys(l string, count int) string {
	first := rand.Intn(30)
	a := first
	b := rand.Intn(30)
	for i := 0; i < count-1; i++ {
		l = swap(l, a, b)
		a = b
		b = rand.Intn(30)
	}
	a = first
	l = swap(l, a, b)
	return l
}

func swap(l string, a int, b int) string {
	r := []rune(l)
	r[a], r[b] = r[b], r[a]
	return string(r)
}
