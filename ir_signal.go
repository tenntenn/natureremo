package natureremo

// IRSignal describes infrared ray signal.
type IRSignal struct {
	Freq   int64   `json:"freq"`
	Data   []int64 `json:"data"`
	Format string  `json:"format"`
}
