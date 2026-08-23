package alert

import (
	"image/color"
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestAlert_API(t *testing.T) {
	a := NewInfo("info-alert", "This is an info alert").Title("Title").Width(400)
	
	if a.variant != Info {
		t.Errorf("expected Info variant")
	}
	if a.title != "Title" {
		t.Errorf("expected Title")
	}
	if a.width != 400 {
		t.Errorf("expected width 400")
	}
}

func TestAlert_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	cases := []struct {
		name      string
		component *Alert
	}{
		{
			name: "Info_WithTitle",
			component: NewInfo("info1", "Please note that the system will go down for maintenance at 12:00 PM.").
				Title("Maintenance Notice").Width(300),
		},
		{
			name: "Success_NoTitle",
			component: NewSuccess("success1", "Your profile has been updated successfully.").Width(300),
		},
		{
			name: "Warning_WithTitle",
			component: NewWarning("warn1", "Your password will expire in 3 days.").
				Title("Password Expiry").Width(300),
		},
		{
			name: "Error_WithTitle_OnClose",
			component: NewError("error1", "Failed to connect to the database. Please check your network connection.").
				Title("Connection Error").Width(300).OnClose(func() {}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tester.Ctx.SetColor(color.RGBA{255, 255, 255, 255})
			tester.Ctx.Clear()

			tc.component.Render(tester.UI, 20, 20)

			score, err := tester.AssertGoldenImage(tc.component, "Alert", tc.name)
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			t.Logf("Score for %s: %f", tc.name, score)
		})
	}
}
