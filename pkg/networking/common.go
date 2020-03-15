package networking

import (
	"time"
	"math/rand"
	"fmt"
	"strconv"
)
func makeInterfaceName() string{
	rand.Seed(time.Now().UTC().UnixNano())
	return fmt.Sprintf("startc-%s", strconv.Itoa(rand.Int())[:8])
}
