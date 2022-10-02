package actions

type TelegramAction struct {
	ActionData
	IsFailureAction bool
	Token           string
	UserKeys        []string
	CustomSound     string
}

type telegramApiData struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
	Sound   string `json:"sound"`
}

func (a *TelegramAction) Run() error {
	return nil
}
