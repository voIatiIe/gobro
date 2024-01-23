package gobro

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func InitCursor(opts ...chromedp.MouseOption) chromedp.Action {
	cursor := 
		`
			var img = new Image();
			img.src = 'https://cdn130.picsart.com/275339252006211.png';
			img.style.position = 'absolute';
			img.style.left = '50%%';
			img.style.top = '50%%';
			img.style.height = '25px';
			img.style.zIndex = '9999';
			img.id = 'cursor'

			document.body.appendChild(img);
		`

	return chromedp.Evaluate(cursor, nil)
}

func drawCursor(x, y float64, ctx context.Context) chromedp.MouseAction {
	cursor := fmt.Sprintf(
		`
			var cursorImage = document.getElementById('cursor');
			cursorImage.style.left = '%f%%';
			cursorImage.style.top = '%f%%';
		`,
		100.0 * x,
		100.0 * y,
	)

	return chromedp.Evaluate(cursor, nil)
}

func MouseClickXY(x, y float64, opts ...chromedp.MouseOption) chromedp.MouseAction {
	click := chromedp.MouseClickXY(x, y, opts...)

	return chromedp.ActionFunc(func(ctx context.Context) error {
		drawCursor(x, y, ctx).Do(ctx)
		return click.Do(ctx)
	})
}
