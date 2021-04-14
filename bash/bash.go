package bash

import "os"

func init() {
	os.Setenv("prod", "true")
}
