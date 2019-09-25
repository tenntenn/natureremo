package natureremo

type TV struct {
	State   *TVState
	Buttons []DefaultButton
}

// TVState is TV state
type TVState struct {
	Input TVInputType `json:"input"`
}

type TVInputType string

const (
	TVInputTypeT  TVInputType = "t"
	TVInputTypeBS TVInputType = "bs"
	TVInputTypeCS TVInputType = "cs"
)
