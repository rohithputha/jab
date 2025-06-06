package memtbl

import (
	"fmt"
	"os"
)

type wal struct {
	logChan <-chan []byte
	logFile *os.File
}

func InitWal(logChan <-chan []byte, logFileName string) *wal {
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	twal := wal{
		logChan: logChan,
		logFile: file,
	}
	return &twal
}

func (wal *wal) ListenToLogs() {
	fmt.Println("listening to logs")
	for {
		select {
		case msg, ok := <-wal.logChan:
			if !ok {
				fmt.Println("logChan closed")
				return
			}
			wal.logFile.Write(append(msg, '\n'))
		}
	}
	fmt.Println("stopped listening to logs")
}

func (wal *wal) Close() {
	wal.logChan = nil
	wal.logFile.Close()
}
