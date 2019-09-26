package natureremo

type Light struct {
	State   *LightState
	Buttons []DefaultButton
}

// LightState is Light state
type LightState struct {
	Brightness string `json:"brightness"`
	Power      string `json:"power"`
	LastButton string `json:"last_button"`
}
