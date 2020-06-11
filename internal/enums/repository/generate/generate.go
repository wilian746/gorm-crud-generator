package generate

type GenerateCommand string

const (
	App     GenerateCommand = "app"
	Unknown GenerateCommand = "unknown"
)

func Values() []GenerateCommand {
	return []GenerateCommand{
		App,
	}
}

func ValueOf(value string) GenerateCommand {
	for _, db := range Values() {
		if string(db) == value {
			return db
		}
	}
	return Unknown
}

func Valid(value string) bool {
	return ValueOf(value) != Unknown
}
