package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const baseUrl = "https://rickandmortyapi.com/api/character?page=%v"

type Response struct {
	Info       Info   `json:"info"`
	Characters []Char `json:"results"`
}

type Info struct {
	Total int `json:"count"`
	Page  int `json:"page"`
}

type Char struct {
	Name string `json:"name"`
}

func find(name string, ch []Char) bool {
	for _, n := range ch {
		if strings.Contains(strings.ToLower(n.Name), strings.ToLower(name)) {
			fmt.Printf("found %s charachter %s \n", name, n.Name)

			return true
		}
	}

	return false
}

func fetchResults(wg *sync.WaitGroup, page int, resultChan chan<- Response, errChan chan<- error) {
	defer wg.Done()

	url := fmt.Sprintf(baseUrl, page)
	res, err := http.Get(url)
	if err != nil {
		errChan <- fmt.Errorf("failed to get page %d: %w", page, err)
		return
	}
	defer res.Body.Close()

	var currResponse Response
	if err = json.NewDecoder(res.Body).Decode(&currResponse); err != nil {
		errChan <- fmt.Errorf("failed to decode page %d: %w", page, err)
		return
	}

	// If response is valid
	if len(currResponse.Characters) == 0 {
		errChan <- fmt.Errorf("page %d not found", page)
		return
	}

	resultChan <- currResponse
}

func TestGoRoutines(t *testing.T) {

	name := os.Args[1]
	start := time.Now()

	totalPages := 50
	//var currResponse Response
	var wg sync.WaitGroup
	resultChan := make(chan Response, totalPages)
	errChan := make(chan error, totalPages)
	wgFind := sync.WaitGroup{}
	wgFind.Add(5)

	for i := 0; i < totalPages; i++ {
		wg.Add(1)
		go fetchResults(&wg, i, resultChan, errChan)
	}

	wg.Wait()
	close(resultChan)
	close(errChan)

	for i := range resultChan {
		find(name, i.Characters)
	}

	/*	var page int
		for {
			var currResponse Response

			url := fmt.Sprintf(baseUrl, page)
			res, err := http.Get(url)

			if err != nil {
				fmt.Printf("failed to get page")
				break
			}

			if err = json.NewDecoder(res.Body).Decode(&currResponse); err != nil {
				fmt.Printf("failed to get page %v", page)
				continue
			}

			if len(currResponse.Characters) == 0 {
				fmt.Printf("finished \n")
				break
			}

			find(name, currResponse.Characters)

			page++
		}*/

	end := time.Now()

	fmt.Printf("took %v\n", end.Sub(start))

}
