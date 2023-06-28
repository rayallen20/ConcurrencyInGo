package intermediate

import (
	"code/chapter5/02-propagationError/customError"
	"code/chapter5/02-propagationError/lowLevel"
	"os/exec"
)

type IntermediateErr struct {
	error
}

func RunJob(id string) error {
	const jobBinPath = "/bad/job/path"
	isExecutable, err := lowLevel.IsGloballyExec(jobBinPath)

	if err != nil {
		return IntermediateErr{customError.WrapError(err, "cannot run job %q: requisite binaries not available", id)}
	} else if isExecutable == false {
		return customError.WrapError(nil, "job binary is not executable") // 本包的错误
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}
