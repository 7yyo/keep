package sys

import (
	"fmt"
	"os"
)

func Exit() {
	fmt.Println("\ngoodbye!")
	os.Exit(1)
}
