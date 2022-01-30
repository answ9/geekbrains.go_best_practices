package main

import (
	"context"
	"fmt"
	"io"
	"lesson1/pkg/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
	neturl "net/url"
)

type Crawler interface {
	Scan(ctx context.Context, url string, curDepth int)
	GetResultChan() <-chan CrawlResult
	Wait()
	Start()
}

type crawler struct {
	maxDepth  int
	req       Requester
	res       chan CrawlResult
	visited   map[string]struct{}
	visitedMu sync.RWMutex
}

type CrawlResult struct {
	Title string
	Url   string
	Err   error
}

func (c *crawler) GetResultChan() <-chan CrawlResult {
	return c.res
}

func (c *crawler) Scan(ctx context.Context, url string, curDepth int) {
	c.visitedMu.RLock()
	if _, ok := c.visited[url]; ok {
		c.visitedMu.RUnlock()
		return
	}
	c.visitedMu.RUnlock()
	if curDepth >= c.maxDepth {
		return
	}
	select {
	case <-ctx.Done():
		return
	default:
		page, err := c.req.GetPage(ctx, url)
		c.visitedMu.Lock()
		c.visited[url] = struct{}{}
		c.visitedMu.Unlock()
		if err != nil {
			c.res <- CrawlResult{Url: url, Err: err}
			return
		}
		title := page.GetTitle()
		c.res <- CrawlResult{
			Title: title,
			Url:   url,
			Err:   nil,
		}
		links := page.GetLinks()
		for _, link := range links {
			go c.Scan(ctx, link, curDepth+1)
		}
	}
}

func NewCrawler(maxDepth int, req Requester) *crawler {
	return &crawler{
		maxDepth: maxDepth,
		req:      req,
		res:      make(chan CrawlResult, 100),
		visited:  make(map[string]struct{}),
	}
}

type Requester interface {
	GetPage(ctx context.Context, url string) (Page, error)
}

type reqWithDelay struct {
	delay time.Duration
	req   Requester
}

func NewRequestWithDelay(delay time.Duration, req Requester) *reqWithDelay {
	return &reqWithDelay{delay: delay, req: req}
}

func (r reqWithDelay) GetPage(ctx context.Context, url string) (Page, error) {
	time.Sleep(r.delay)
	return r.req.GetPage(ctx, url)
}

type requester struct {
	timeout time.Duration
}

func NewRequester(timeout time.Duration) *requester {
	return &requester{timeout: timeout}
}

func (r requester) GetPage(ctx context.Context, url string) (Page, error) {
	cl := &http.Client{
		Timeout: r.timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	rawPage, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawPage.Body.Close()
	return NewPage(rawPage.Body)
}

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return page{}, err
	}
	return page{doc}, nil
}

func (p page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

func (p page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			if !IsAbsoluteUrl(url) {
				return
			}
			urls = append(urls, url)
		}
	})
	return urls
}

func IsAbsoluteUrl(str string) bool {
	u, err := neturl.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func main() {
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}
	var r Requester
	r = NewRequester(time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)
	crawler := NewCrawler(cnfg.MaxDepth, r)
	crawler.Scan(ctx, cnfg.StartUrl, cnfg.CurDepth)
	go processResult(ctx, crawler.GetResultChan(), cnfg, cancel)
	for {
		select {
		case s := <-chSig:

			switch s {
			case syscall.SIGINT:
				fmt.Println("SIGINT was received")
				cancel()
			case syscall.SIGTERM:
				fmt.Println("SIGTERM was received")
				cancel()
			case syscall.SIGUSR1:
				fmt.Printf("SIGUSR1 was received. Max depth was increased by 2\n")
				crawler := NewCrawler(cnfg.MaxDepth+2, r)
				crawler.Scan(ctx, cnfg.StartUrl, cnfg.CurDepth)
				go processResult(ctx, crawler.GetResultChan(), cnfg, cancel)
			default:
				fmt.Println("Unknown signal")
			}

		case <-ctx.Done():
			fmt.Printf("context canceled")
			return
		}
	}
	cancel()
}

func processResult(ctx context.Context, in <-chan CrawlResult, cnfg *config.AppConfig, cancel context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(cnfg.TimeOut)*time.Second)
	defer func() {
		cancelFunc()
		os.Exit(1)
	}()

	var errCount int
	for {
		select {
		case res := <-in:
			if res.Err != nil {
				errCount++
				fmt.Printf("ERROR Link: %s, err: %v\n", res.Url, res.Err)
				if errCount >= cnfg.MaxErrors {
					fmt.Printf("Max errors count is reached - %d\n", cnfg.MaxErrors)
					cancel()
				}
			} else {
				fmt.Printf("Link: %s, Title: %s\n", res.Url, res.Title)
			}
		case <-ctx.Done():
			fmt.Printf("context canceled")
			return
		}
	}
}
