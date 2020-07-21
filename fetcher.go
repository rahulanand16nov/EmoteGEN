package main

import (
	"os"
	"net/http"
	"io"
	"strconv"
	"fmt"
	"sync"
)

const MAX_EMOTES int = 50000
const MAX_WORKERS int = 100

func executeTask(wg *sync.WaitGroup, tasks chan int) error {
	for task := range tasks {
		fmt.Println("Executing %d",task)
		// Create the directory if not exisit
		_ = os.Mkdir("emotes", os.ModeDir)

		// Get the data
		resp, err := http.Get("https://static-cdn.jtvnw.net/emoticons/v1/"+strconv.Itoa(task)+"/1.0")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Status code not ok")
			continue
		}
		// Crease the file
		out, err := os.Create("emotes/"+strconv.Itoa(task)+".png")
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}
	wg.Done()
	return nil
}


func main () {
	tasks := make(chan int)
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func () {
		for i := 1; i<=MAX_EMOTES; i++ {
			tasks <- i
		}
		close(tasks)
		wg.Done()
	}()

	for i := 1; i <= MAX_WORKERS; i++ {
		wg.Add(1)
		go executeTask(&wg, tasks)
	}
	wg.Wait()
}