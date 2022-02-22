package urlprocessor

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"sync"
)

const (
	httpPrefix  = "http://"
	httpsPrefix = "https://"
)

type Args struct {
	ParallelLimit int
	URLs          []string
	HTTPClient    *http.Client
}

func New(args Args) (*Processor, error) {
	if len(args.URLs) == 0 {
		return nil, errors.New("URLs cannot be empty")
	}
	if args.ParallelLimit == 0 {
		return nil, errors.New("ParallelLimit should be a value greater than zero")
	}
	if args.HTTPClient == nil {
		return nil, errors.New("HTTPClient is required")
	}

	ch := make(chan struct{}, args.ParallelLimit)
	for i := 0; i < args.ParallelLimit; i++ {
		ch <- struct{}{}
	}

	return &Processor{
		urls:   args.URLs,
		ch:     ch,
		client: args.HTTPClient,
	}, nil
}

type Processor struct {
	urls   []string
	ch     chan struct{}
	client *http.Client
}

type HashResult struct {
	URL     string
	MD5Hash string
}

func (p Processor) Process(resultCh chan<- HashResult, errCh chan<- error, doneCh chan<- struct{}) {
	var wg sync.WaitGroup
	wg.Add(len(p.urls))

	for _, url := range p.urls {
		url := url
		// will block if more than the limit is running in parallel
		<-p.ch
		go func() {
			defer wg.Done()
			defer func() {
				// signal channel to allow other blocked goroutines to run
				p.ch <- struct{}{}
			}()

			parsedURL, err := parseUrl(url)
			if err != nil {
				errCh <- err
				return
			}

			response, err := p.client.Get(*parsedURL)
			if err != nil {
				errCh <- err
				return
			}

			defer response.Body.Close()

			content, err := io.ReadAll(response.Body)
			if err != nil {
				errCh <- err
				return
			}

			md5Hex := toMD5Hex(content)
			resultCh <- HashResult{MD5Hash: md5Hex, URL: url}
		}()
	}

	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
}

func parseUrl(urlStr string) (*string, error) {
	if !strings.HasPrefix(urlStr, httpPrefix) && !strings.HasPrefix(urlStr, httpsPrefix) {
		// let http.Get redirect to https if found
		urlStr = httpPrefix + urlStr
	}
	// try parse url
	url, err := neturl.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	parsed := url.String()
	return &parsed, nil
}

func toMD5Hex(bytes []byte) string {
	hasher := md5.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}
