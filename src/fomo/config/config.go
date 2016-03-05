package config

/**
[pis]
pi@pz.kellego.us
pi@turtle.kellego.us

[servers]
flint.kellego.us
**/

type Config struct {
	HostSets map[string][]string
}
