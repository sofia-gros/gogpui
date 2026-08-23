//go:build !windows

package gogpui

// checkWindowMinimized は Windows 以外のプラットフォーム用のスタブ。
func checkWindowMinimized(title string) (bool, int, int) {
	return false, 0, 0
}
