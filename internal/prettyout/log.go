package prettyout

import (
	"github.com/the-yex/gvm/internal/prettyout/color"
	"io"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/9 下午6:03
* @Package:
 */

func PrettyInfo(out io.Writer, format string, a ...interface{}) {
	color.New(color.FgGreen).Fprintf(out, format, a...)
}

func PrettyWarm(out io.Writer, format string, a ...interface{}) {
	color.New(color.FgYellow).Fprintf(out, format, a...)

}

func PrettyError(out io.Writer, format string, a ...interface{}) {
	color.New(color.FgRed).Fprintf(out, format, a...)

}
