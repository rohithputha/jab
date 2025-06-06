package memtbl

import (
	"github.com/rohithputha/jab/datatypes"
	"sync"
	"time"
)

// Each transaction should be a start as a new go routine and during validate and commit phase should be streamlined into a critical section

type TrsnChange struct {
	finishTime time.Time
	changeKeys []*datatypes.Key // should this be pointers?
	getKeys    []*datatypes.Key // should this be pointers?
}
type TrsnCommand struct {
	command  string
	context  *trsnContext
	copyData bool
}

func (cmd *TrsnCommand) getCommand() string {
	return cmd.command
}

func (cmd *TrsnCommand) execute() error {
	// parse and then execute
	change := &TrsnChange{
		finishTime: time.Now(),
		changeKeys: make([]*datatypes.Key, 0),
		getKeys:    make([]*datatypes.Key, 0),
	}
	cmd.context.changeChannel <- change
	return nil
}

type trsnContext struct {
	memTablesMap  map[string]*Tbl
	contextTable  *Tbl
	changeChannel chan *TrsnChange
}
type Trsn struct {
	commands []*TrsnCommand
	context  *trsnContext

	startTime time.Time
	endTime   time.Time
}

func (tc *Trsn) StartTransaction() {
	for _, cmd := range tc.commands {
		err := cmd.execute()
		if err != nil {
			return // handle transaction failures
		}
	}
	validator := InitTrasactionValidator()
	validator.changeChannels <- tc.context.changeChannel

}
func CreateTrsn(commands []string) *Trsn {
	return nil
}

// #############################################################
// Single Transaction validator and writer

// this should be in a another file
var trsnVOnce sync.Once
var transactionValidator *trsnValidatorAndWriter

type trsnValidatorAndWriter struct {
	changeChannels chan chan *TrsnChange
}

func (tc *trsnValidatorAndWriter) log() {}

func (tc *trsnValidatorAndWriter) validateChangeChannel(changeChannel chan *TrsnChange) error {
	return nil
}
func (tvw *trsnValidatorAndWriter) validate() {
	select {
	case change := <-tvw.changeChannels:
		tvw.validateChangeChannel(change)
	}
}
func (tvw *trsnValidatorAndWriter) write() {

}

func InitTrasactionValidator() *trsnValidatorAndWriter {
	trsnVOnce.Do(func() {
		transactionValidator = &trsnValidatorAndWriter{
			changeChannels: make(chan chan *TrsnChange),
		}
		go transactionValidator.validate()
	})
	return transactionValidator
}
