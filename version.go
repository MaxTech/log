package max_log

import "fmt"

var Version string

func init() {
    Version = fmt.Sprintf(
        "|- %s module:\t\t\t%s",
        "logger",
        "0.0.1")
}
