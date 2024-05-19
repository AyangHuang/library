package check

import "context"

type Checker interface {
	Check() bool
}

type DirectChecker interface {
	DirectCheck() bool
}

// DChecker DoubleCheck，DirectChecker || Checker
type DChecker interface {
	DirectChecker
	Checker
}

type CheckerFun func() bool

func (check CheckerFun) Check() bool {
	return check()
}

type DirectCheckerFun func() bool

func (check DirectCheckerFun) DirectCheck() bool {
	return check()
}

func Check(ctx context.Context, checkers ...interface{}) bool {
	for _, check := range checkers {
		switch check.(type) {
		case DChecker:
			if ok := check.(DirectChecker).DirectCheck(); ok {
				continue
			}
			if ok := check.(Checker).Check(); !ok {
				return false
			}
		case DirectChecker:
			if ok := check.(DirectChecker).DirectCheck(); ok {
				return true
			}
		case Checker:
			if ok := check.(Checker).Check(); !ok {
				return false
			}
		default:
			panic("do not implement Checker or DirectChecker interface")
		}

	}
	return true
}
