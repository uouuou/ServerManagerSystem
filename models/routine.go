package models

import (
	"log"
	"runtime/debug"
)

// NewRoutine 采用提前recover的方式终止因为goroutine错误导致的整体崩溃
func NewRoutine(f func()) {
	go func() {
		defer func() {
			// Recover from panic.
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Println(err)
				log.Println(stack)
			}
		}()

		f()
	}()
}
