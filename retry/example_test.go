package retry

import (
	"context"
)

func ExampleRedundant() {
	count := 2
	p90 := 200
	param := "param"

	r1 := NewRedundant[string](count, p90, doHttpV1)

	r2 := NewRedundant[string](count, p90, func() (string, error) {
		return doHttpV2(param)
	})

	respInt := 0
	respString := ""
	r3 := NewRedundant[struct{}](count, p90, func() (struct{}, error) {
		if innerRespInt, innerRespString, err := doHttpV3(param); err == nil {
			respInt = innerRespInt
			respString = innerRespString
			return struct{}{}, nil
		} else {
			return struct{}{}, err
		}
	})

	_, _ = r1.Do(context.TODO())
	_, _ = r2.Do(context.TODO())
	_, _ = r3.Do(context.TODO())
	_ = respInt
	_ = respString
}

func doHttpV1() (string, error) {
	return "param", nil
}

func doHttpV2(param string) (string, error) {
	return param, nil
}

func doHttpV3(param string) (int, string, error) {
	return 0, param, nil
}
