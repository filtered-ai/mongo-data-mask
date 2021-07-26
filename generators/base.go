package generators

import "github.com/MichaelTJones/pcg"

type Base struct {
	Seed  uint64
	Pcg32 *pcg.PCG32
}
