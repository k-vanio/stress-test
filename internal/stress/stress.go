package stress

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Stress struct {
	mu          *sync.Mutex
	url         string
	requests    int
	concurrency int
}

func New() *Stress {
	return &Stress{
		mu: &sync.Mutex{},
	}
}

func (s *Stress) Run(cmd *cobra.Command, args []string) {
	if err := s.validateArgs(args); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	statusOk := 0
	statusNotOk := make(map[int]int)
	startTime := time.Now()

	await := make(chan struct{}, s.concurrency)
	var wg sync.WaitGroup

	for i := 0; i < s.requests; i++ {
		await <- struct{}{}
		wg.Add(1)
		go func() {
			s.makeRequest(&statusNotOk, &statusOk)
			defer func() {
				<-await
				defer wg.Done()
			}()
		}()
	}

	wg.Wait()
	duration := time.Since(startTime).String()

	statusNotOkString := ""
	for k, v := range statusNotOk {
		statusNotOkString += fmt.Sprintf("statusCode %d = %d\n", k, v)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Total time spent during execution",
		"Total quantity of requests executed",
		"Number of requests with HTTP 200 status",
		"Distribution of other HTTP status codes",
	})
	table.Append([]string{duration, strconv.Itoa(s.requests), strconv.Itoa(statusOk), statusNotOkString})

	table.Render()
}

func (s *Stress) validateArgs(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("—url, —requests, —concurrency are required")
	}

	m := make(map[string]string, 3)
	for _, arg := range args {
		part := strings.SplitN(arg, "=", 2)
		if len(part) == 2 {
			m[part[0]] = part[1]
		}
	}

	if _, ok := m["—url"]; !ok {
		return fmt.Errorf("—url is required")
	}
	s.url = m["—url"]

	if _, ok := m["—requests"]; !ok {
		return fmt.Errorf("—requests is required")
	}

	requests, err := strconv.Atoi(m["—requests"])
	if err != nil {
		return fmt.Errorf("—requests must be a number")
	}
	s.requests = requests

	if _, ok := m["—concurrency"]; !ok {
		return fmt.Errorf("—concurrency is required")
	}

	concurrency, err := strconv.Atoi(m["—concurrency"])
	if err != nil {
		return fmt.Errorf("—concurrency must be a number")
	}
	s.concurrency = concurrency

	return nil
}

func (s *Stress) makeRequest(statusNotOk *map[int]int, statusOk *int) {
	resp, err := http.Get(s.url)
	if err != nil {
		s.mu.Lock()
		(*statusNotOk)[500]++
		s.mu.Unlock()
		return
	}

	s.mu.Lock()
	if resp.StatusCode == 200 {
		*statusOk++
	} else {
		(*statusNotOk)[resp.StatusCode]++
	}
	s.mu.Unlock()
}
