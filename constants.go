package main

import "fmt"

var (
	FailedToGenerateCommand                = fmt.Errorf("failed to generate command")
	FailedToStartCommand                = fmt.Errorf("failed to start command")
	FailedToWaitCommand                = fmt.Errorf("failed to wait command")
	FailedInitializeStdPipe               = fmt.Errorf("failed to initialize std pipe")
	FailedToGenerateMasterFile               = fmt.Errorf("failed to generate master file")
	IllegalStatemachineEvent     = fmt.Errorf("illegalStatemachineEvent")
	NoRowFound                   = fmt.Errorf("no rows found")
	EmptyStringNotExpected       = fmt.Errorf("non-empty string expected")
	IllegalEventHandlerEntity    = fmt.Errorf(" illegal event in event handler")
	TransactionNotFoundInContext = fmt.Errorf("transaction not found in context")
)
