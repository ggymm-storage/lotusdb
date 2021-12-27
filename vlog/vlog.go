package vlog

import (
	"errors"
	"fmt"
	"github.com/flowercorp/lotusdb/logfile"
	"sync"
	"time"
)

var (
	ErrActiveLogFileNil = errors.New("active log file not exists")
	ErrLogFileNil       = errors.New("log file %d not exists")
)

type (
	// ValueLog value log.
	ValueLog struct {
		sync.RWMutex
		opt           options
		activeLogFile *logfile.LogFile            // current active log file for writing.
		logFiles      map[uint32]*logfile.LogFile // all log files. Must hold the mutex before modify it.
	}

	options struct {
		path      string
		blockSize int64
		IoType    logfile.IOType
	}

	// ValuePos value position.
	ValuePos struct {
		fid    uint32
		offset int64
		size   uint32
	}
)

// OpenValueLog create a new value log file.
func OpenValueLog() (*ValueLog, error) {
	return nil, nil
}

func (vlog *ValueLog) ReadValue(pos *ValuePos) ([]byte, error) {
	if pos == nil {
		return nil, nil
	}

	var logFile *logfile.LogFile
	if pos.fid == vlog.activeLogFile.Fid {
		logFile = vlog.activeLogFile
	} else {
		vlog.RLock()
		logFile = vlog.logFiles[pos.fid]
		vlog.RUnlock()
	}
	if logFile == nil {
		return nil, fmt.Errorf(ErrLogFileNil.Error(), pos.fid)
	}

	logEntry, err := logFile.Read(pos.offset)
	if err != nil {
		return nil, err
	}

	// check whether value is expired.
	if logEntry.ExpiredAt <= uint64(time.Now().Unix()) {
		// delete expired value.todo
		return nil, nil
	}
	return logEntry.Value, nil
}

func (vlog *ValueLog) Write(e *logfile.LogEntry) (*ValuePos, error) {
	err := vlog.activeLogFile.Write(e)
	if err != nil {
		return nil, nil
	}
	return nil, nil
}

func (vlog *ValueLog) Sync() error {
	if vlog.activeLogFile == nil {
		return ErrActiveLogFileNil
	}

	vlog.activeLogFile.Lock()
	defer vlog.activeLogFile.Unlock()
	return vlog.activeLogFile.Sync()
}

func (vlog *ValueLog) Close() error {
	if vlog.activeLogFile == nil {
		return ErrActiveLogFileNil
	}

	vlog.activeLogFile.Lock()
	defer vlog.activeLogFile.Unlock()
	return vlog.activeLogFile.Close()
}

// do it later.
func (vlog *ValueLog) compact() {
	// todo
}