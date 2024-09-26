package assert

import "reflect"

func Equal[variable comparable](got, expected variable, crashMessage string) {
	if !reflect.DeepEqual(got, expected) {
		panic(crashMessage)
	}
}
