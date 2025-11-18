package prettyout

import (
	"github.com/the-yex/gvm/internal/prettyout/color/colorable"
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/9 下午6:09
* @Package:
 */
func Test_Pretty(t *testing.T) {
	//out := colorable.NewColorableStdout()
	t.Run("info", func(t *testing.T) {
		out := colorable.NewColorableStdout()
		PrettyInfo(out, "%s\n", "hello world")
	})
	t.Run("warning", func(t *testing.T) {
		out := colorable.NewColorableStdout()
		PrettyWarm(out, "%s\n", "hello world")
	})
	t.Run("error", func(t *testing.T) {
		out := colorable.NewColorableStdout()
		PrettyError(out, "%s\n", "hello world")
	})
}
