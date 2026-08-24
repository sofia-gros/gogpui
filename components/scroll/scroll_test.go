package scroll

import (
	"image/color"
	"testing"

	"github.com/sofiagros/gogpui/core/layout"
	"github.com/sofiagros/gogpui/core/context"
	testutil "github.com/sofiagros/gogpui/testing"
)

func TestScroll_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	// Create a tall custom widget to scroll
	tallWidget := &layout.CustomWidget{
		RenderFunc: func(uictx *context.UIContext, x, y float64) (float64, float64) {
			if uictx.MeasureOnly {
				return 400, 1000
			}
			ctx := uictx.GG
			ctx.SetColor(color.RGBA{255, 0, 0, 255})
			ctx.DrawRectangle(x, y, 400, 200) // Top red box
			uictx.Fill()
			
			ctx.SetColor(color.RGBA{0, 255, 0, 255})
			ctx.DrawRectangle(x, y+400, 400, 200) // Middle green box
			uictx.Fill()
			
			ctx.SetColor(color.RGBA{0, 0, 255, 255})
			ctx.DrawRectangle(x, y+800, 400, 200) // Bottom blue box
			uictx.Fill()
			
			return 400, 1000
		},
	}

	sc := New("test-scroll").Size(400, 300).Child(tallWidget)

	// Inject scroll state manually to test offset
	t.Run("Scroll_Top", func(t *testing.T) {
		tester.Ctx.SetColor(color.RGBA{255, 255, 255, 255})
		tester.Ctx.Clear()
		
		sc.Render(tester.UI, 0, 0)
		_, err := tester.AssertGoldenImage(sc, "Scroll", "Top")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Scroll_Middle", func(t *testing.T) {
		tester.Ctx.SetColor(color.RGBA{255, 255, 255, 255})
		tester.Ctx.Clear()
		
		state := tester.UI.GetState("test-scroll")
		state.OffsetY = 450 // Scroll down to middle green box
		
		sc.Render(tester.UI, 0, 0)
		_, err := tester.AssertGoldenImage(sc, "Scroll", "Middle")
		if err != nil {
			t.Fatal(err)
		}
	})
}
