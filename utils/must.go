package utils

func Must0(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[R any](res R, err error) R {
	Must0(err)
	return res
}
