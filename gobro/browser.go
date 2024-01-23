package gobro

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type BrowserOpts struct {
	width, height int64
	quality int
	counter int
}

func defaultBrowserOpts() *BrowserOpts {
	return &BrowserOpts{
		width: 1920,
		height: 1080,
		quality: 100,
	}
}

type BrowserOptsFunc func(*BrowserOpts)

func WithQuality(quality int) BrowserOptsFunc { return func(opts *BrowserOpts) { opts.quality = quality } }

func WithWidth(width int64) BrowserOptsFunc { return func(opts *BrowserOpts) { opts.width = width } }

func WithHeight(height int64) BrowserOptsFunc { return func(opts *BrowserOpts) { opts.height = height } }

type Browser struct { 
	Opts *BrowserOpts
	Ctx context.Context
	Cancel context.CancelFunc
}

func NewBrowser(url string, opts ...BrowserOptsFunc) (*Browser, error) {
	defaultOpts := defaultBrowserOpts()

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Headless,
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), allocOpts...)

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	for _, opt := range opts { opt(defaultOpts) }

	browser := &Browser{
		Opts: defaultOpts,
		Ctx: ctx,
		Cancel: cancel,
	}

	tasks := chromedp.Tasks{
		chromedp.EmulateViewport(browser.Opts.width, browser.Opts.height),
		chromedp.Navigate(url),
		InitCursor(),
	}

	if err := chromedp.Run(browser.Ctx, tasks); err != nil { return nil, err }

	return browser, nil
}


func (b *Browser) MoveCursor(x, y float64, buffer *[]byte) error {
	click := MouseClickXY(x, y, chromedp.ButtonNone, chromedp.ButtonModifiers())
	screenshot := chromedp.FullScreenshot(buffer, b.Opts.quality)

	start := time.Now()
	if err := chromedp.Run(
		b.Ctx,
		chromedp.Tasks{
			click,
			screenshot,
		},
	); err != nil { return err }
	log.Printf("Execution time %s\n", time.Since(start))

	return nil
}
