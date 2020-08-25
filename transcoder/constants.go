package transcoder

import "fmt"

var (
	FailedToGenerateCommand    = fmt.Errorf("failed to generate command")
	FailedToStartCommand       = fmt.Errorf("failed to start command")
	FailedToWaitCommand        = fmt.Errorf("failed to wait command")
	FailedInitializeStdPipe    = fmt.Errorf("failed to initialize std pipe")
	FailedToGenerateMasterFile = fmt.Errorf("failed to generate master file")
)
