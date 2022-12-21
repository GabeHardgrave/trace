package trace

import (
	"fmt"
	"io"
)

type ErrorLocation struct {
	// The name of the file the error occurred at.
	File string `json:"file"`
	// The line the error occurred at.
	Line int `json:"line"`
	// Any extra data attached to the error site
	Details []any `json:"details"`
}

func (loc *ErrorLocation) format(w io.Writer) {
	fmt.Fprintf(w, "%s:%d", loc.File, loc.Line)

	if len(loc.Details) > 0 {
		fmt.Fprintf(w, "{ ")
		fmtArg, rest := loc.Details[0], loc.Details[1:]
		if fmtStr, ok := fmtArg.(string); ok {
			fmt.Fprintf(w, fmtStr, rest...)
		} else {
			fmt.Fprintf(w, "%+v", loc.Details)
		}
		fmt.Fprintf(w, " }")
	}
}

func (loc *ErrorLocation) formatSizeHint() int {
	l := len(loc.File) + len(":") + 4 // 4 because most code files aren't more than 9999 lines
	if len(loc.Details) > 0 {
		l += len("{ ") + len(" }")
		l += 20 // completely arbitary, but may as well guess something
	}
	return l
}
