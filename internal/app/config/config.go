package config

import (
	"flag"
)

func Flags() bool {

	flag.String("d", "", "choosing a storage location")

	flag.Parse()

	return isFlagPassed("d")

}
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
