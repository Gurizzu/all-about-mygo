package service

func OverrideError(err error) (res string) {
	if err == nil {
		return ""
	}
	return err.Error()
}
