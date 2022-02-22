package urlprocessor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestProcessor_Process(t *testing.T) {
	tt := []struct {
		name                string
		statusCodes         []int
		contents            []string
		urls                []string
		expectedResult      []string
		expectedErrorsCount int
	}{
		{
			name: "happy path",
			statusCodes: []int{
				200,
				200,
				200,
			},
			contents: []string{
				"1",
				"2",
				"3",
			},
			urls: []string{
				"%v/news",
				"%v/home",
				"xx",
			},
			expectedResult: []string{
				"c4ca4238a0b923820dcc509a6f75849b",
				"c81e728d9d4c2f636f067f89cc14862c",
			},
			expectedErrorsCount: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var opResults []string
			var opErrors []string
			resultCh := make(chan HashResult)
			errorCh := make(chan error)
			doneCh := make(chan struct{})
			triesCount := 0

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				current := triesCount
				triesCount++
				w.WriteHeader(tc.statusCodes[current])
				t.Logf("url %v", r.URL)
				_, _ = w.Write([]byte(tc.contents[current]))
			}))
			defer server.Close()

			for i := 0; i < len(tc.urls); i++ {
				tc.urls[i] = fmt.Sprintf(tc.urls[i], server.URL)
			}

			processor, _ := New(Args{
				URLs:          tc.urls,
				ParallelLimit: 10,
				HTTPClient:    server.Client(),
			})

			processor.Process(resultCh, errorCh, doneCh)

		loop:
			for {
				select {
				case opResult := <-resultCh:
					opResults = append(opResults, opResult.MD5Hash)
				case opError := <-errorCh:
					opErrors = append(opErrors, opError.Error())
				case <-doneCh:
					break loop
				}
			}

			sort.Strings(opResults)
			sort.Strings(opErrors)

			if !reflect.DeepEqual(opResults, tc.expectedResult) {
				t.Logf("expected %v results, got %v", tc.expectedResult, opResults)
				t.Fail()
			}

			if tc.expectedErrorsCount != len(opErrors) {
				t.Logf("expected %v errors, got %v", tc.expectedErrorsCount, len(opErrors))
				t.Fail()
			}
		})
	}
}
