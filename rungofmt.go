package rungofmt

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// ErrNotRegularFile indicate given source file is not regular file.
var ErrNotRegularFile = errors.New("given file is not regular file")

// ErrStopWithNonZeroExitCode indicate the exit code of stopped program is not zero.
type ErrStopWithNonZeroExitCode int

func (e ErrStopWithNonZeroExitCode) Error() string {
	return "[ErrStopWithNonZeroExitCode: " + strconv.FormatInt(int64(e), 10) + "]"
}

// RunGoFmt run gofmt on given sourceFilePath with `-w` and optional `-s` (simplifyCode) options.
func RunGoFmt(sourceFilePath string, simplifyCode bool) (err error) {
	if fileinfo, err := os.Stat(sourceFilePath); nil != err {
		return err
	} else if (fileinfo.Mode() & os.ModeType) != 0 {
		return ErrNotRegularFile
	}
	args := []string{"-w"}
	if simplifyCode {
		args = append(args, "-s")
	}
	args = append(args, sourceFilePath)
	cmd := exec.Command("gofmt", args...)
	if err = cmd.Start(); nil != err {
		log.Printf("ERROR: start gofmt (arguments: %v) failed %v.", args, err)
		return
	}
	log.Printf("INFO: run gofmt (PID=%v) with: %v", cmd.Process.Pid, args)
	if err = cmd.Wait(); nil != err {
		log.Printf("ERROR: gofmt stopped with error: %v", err)
	}
	if exitCode := cmd.ProcessState.ExitCode(); exitCode != 0 {
		log.Printf("WARN: gofmt stopped with non-zero exit code: %d.", exitCode)
		return ErrStopWithNonZeroExitCode(exitCode)
	}
	log.Printf("INFO: gofmt stopped.")
	return nil
}
