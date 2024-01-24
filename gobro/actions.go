package gobro

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func InitCursor(opts ...chromedp.MouseOption) chromedp.Action {
	cursor := 
		`
			var cursor = document.createElement('div');
			cursor.style.width = '6px';
			cursor.style.height = '6px';
			cursor.style.backgroundColor = 'black';
			cursor.style.borderRadius = '50%%';
			cursor.style.position = 'absolute';
			// cursor.style.transform = 'translate(-50%%, -50%%)';
			cursor.style.zIndex = '9999';
			cursor.id = 'cursor';

			document.body.appendChild(cursor);
		`

	return chromedp.Evaluate(cursor, nil)
}

func drawCursor(x, y float64, ctx context.Context) chromedp.MouseAction {
	cursor := fmt.Sprintf(
		`
			var cursor = document.getElementById('cursor');

			try {
				if (cursor === null) {
					cursor = document.createElement('div');

					cursor.style.width = '6px';
					cursor.style.height = '6px';
					cursor.style.backgroundColor = 'black';
					cursor.style.borderRadius = '50%%';
					cursor.style.position = 'absolute';
					// cursor.style.transform = 'translate(-50%%, -50%%)';
					cursor.style.zIndex = '9999';
					cursor.id = 'cursor';

					document.body.appendChild(cursor);
				}
				cursor.style.left = '%f%%';
				cursor.style.top = '%f%%';
			}
			catch(err) {}
		`,
		100.0 * x,
		100.0 * y,
	)

	return chromedp.Evaluate(cursor, nil)
}

func MouseClickXY(x, y float64, opts ...chromedp.MouseOption) chromedp.MouseAction {
	click := chromedp.MouseClickXY(x, y, opts...)

	return chromedp.ActionFunc(func(ctx context.Context) error {
		if err := click.Do(ctx); err != nil {
			return err
		}

		return drawCursor(x, y, ctx).Do(ctx)
	})
}


func Screenshot(res *[]byte, quality int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		format := page.CaptureScreenshotFormatPng
		if quality != 100 {
			format = page.CaptureScreenshotFormatJpeg
		}

		var err error
		*res, err = page.CaptureScreenshot().
			WithCaptureBeyondViewport(false).
			WithFromSurface(true).
			WithFormat(format).
			WithQuality(int64(quality)).
			Do(ctx)

		return err
	})
}
