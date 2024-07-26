package fontsync

import (
	"fmt"
	"log"
)

func logWithCpu(cpu int, message string, v ...any) {
	log.Printf(fmt.Sprintf("CPU %d: ", cpu)+message, v...)
}
