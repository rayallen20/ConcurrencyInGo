package lowLevel

import (
	"code/chapter5/02-propagationError/customError"
	"os"
)

type LowLevelErr struct {
	error
}

func IsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{error: customError.WrapError(err, err.Error())}
	}

	return info.Mode().Perm()&0100 == 0100, nil
}
