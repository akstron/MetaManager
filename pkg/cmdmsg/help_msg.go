package cmdmsg

func GetForceCommandMessage() string {
	return `Use -f or --force to force initialize.`
}

func ErrorOccurredMessage() string {
	return `Error occurred while running the command: `
}
