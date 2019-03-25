package main

import (
	"flag"
	"sync"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Counter struct {
	Str   string
	Count int
}

var action = flag.String("a", "single", "action")
var filename = flag.String("f", "list.txt", "list filename")

func main() {
	flag.Parse()
	var counters []Counter

	switch *action {
	default:
		fmt.Println("single mode(default)")
		counters = single()
	case "multi":
		fmt.Println("multi mode(unrestricted)")
		counters = multi()
	case "restrict":
		fmt.Println("multi mode(restricted)")
		counters = restrict()
	case "confluence":
		fmt.Println("synchronize mode")
		counters = confluence()
	}

	fmt.Println("done.")
	for _, counter := range counters {
		fmt.Print(counter.Str + ": " + strconv.Itoa(counter.Count) + " ")
	}
	fmt.Println("")
}

func search(str string, filename string) (count int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	count = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), str) {
			count++
		}
	}
	return
}

func targets() (targets []string) {
	for i := 0; i < 16; i++ {
		targets = append(targets, fmt.Sprintf("0%x", i))
	}
	for i := 16; i < 256; i++ {
		targets = append(targets, fmt.Sprintf("%x", i))
	}
	return
}

func sum(counters []Counter) (num int) {
	for _, counter := range counters {
		num += counter.Count
	}
	return
}

func single() (counters []Counter) {
	for _, str := range targets() {
		counters = append(counters, Counter{Str: str, Count: search(str, *filename)})
	}
	return
}

func multi() (counters []Counter) {
	for _, str := range targets() {
		go func(s string) {
			counters = append(counters, Counter{Str: s, Count: search(s, *filename)})
		}(str)
	}
	return
}

func confluence() (counters []Counter) {
	wg := &sync.WaitGroup{}
	for _, str := range targets() {
		go func(s string) {
			wg.Add(1)
			counters = append(counters, Counter{Str: s, Count: search(s, *filename)})
			wg.Done()
		}(str)
	}
	wg.Wait()
	return
}

func restrict() (counters []Counter) {
	c := make(chan bool, 10)
	wg := &sync.WaitGroup{}
	for _, str := range targets() {
		c <- true
		go func(s string) {
			defer func() { <-c }()
			wg.Add(1)
			counters = append(counters, Counter{Str: s, Count: search(s, *filename)})
			wg.Done()
		}(str)
	}
	wg.Wait()
	return
}
