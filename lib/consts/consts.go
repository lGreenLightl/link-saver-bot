package consts

const (
	GetUpdatesMethod  = "getUpdates"
	SendMessageMethod = "sendMessage"
)

type Type int

const (
	Unknown Type = iota
	Message
)

const (
	DefaultPermission = 0774
)

const (
	HelpCommand  = "/help"
	RndCommand   = "/rnd"
	StartCommand = "/start"
)

const MessageHelp = `I can save and keep your pages. Also I can offer you them to read.

In order to save the page, just send me all link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After that, this page will be removed from your list!`

const (
	MessageHello          = "Hi! 👋\n\n" + MessageHelp
	MessageUnknownCommand = "Unknown command 🤨"
	MessageNoSavedPages   = "You have no saved pages 🤷🏼‍♂️"
	MessageSaved          = "Saved! 👌"
	MessageAlreadyExists  = "You already have this page in your list ✌️"
)

const (
	HostPath    = "api.telegram.org"
	StoragePath = "../file-storage"
)

const (
	BatchSize = 100
)
