package shell

import (
	"fmt"
	"os"
)

type LogHandler interface {
	Log(level int, args ...any)
}

type ConsoleLogHandler struct {
	level int
}

func NewConsoleLogHandler() *ConsoleLogHandler {
	return &ConsoleLogHandler{
		level: 0,
	}
}

func (c *ConsoleLogHandler) Log(level int, args ...any) {
	if c.level < level {
		return
	}

	fmt.Print(args...)
}

func (c *ConsoleLogHandler) SetLogLevel(level int) {
	c.level = level
}

type FileLogHandler struct {
	FilePath string
	level    int
}

func NewFileLogHandler(filePath string) *FileLogHandler {
	return &FileLogHandler{
		FilePath: filePath,
		level:    4, // Default log level set to debug
	}
}

func (f *FileLogHandler) Log(level int, args ...any) {
	if f.level < level {
		return
	}

	file, err := os.OpenFile(f.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}

	_, err = fmt.Fprint(file, args...)

	if err != nil {
		fmt.Println("Error writing to log file:", err)
	}

	err = file.Close()

	if err != nil {
		fmt.Println("Error closing log file:", err)
	}

	f.rotateFile()
}

func (f *FileLogHandler) rotateFile() {
	const maxSize = 10 * 1024 * 1024 // 10 MB

	info, err := os.Stat(f.FilePath)

	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	if info.Size() > maxSize {
		err = os.Rename(f.FilePath, f.FilePath+".old")

		if err != nil {
			fmt.Println("Error rotating log file:", err)
		}
	}
}
