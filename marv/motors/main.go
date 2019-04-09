package motors

import "context"

var leftSpeed = 0.0
var rightSpeed = 0.0

// Command NEEDSCOMMENT
type Command struct {
	Left  float64
	Right float64
}

// Start NEEDSCOMMENT
func Start(ctx context.Context) chan<- *Command {
	commands := make(chan *Command)

	go func() {
		defer close(commands)
		for {
			select {
			case command := <-commands:
				handleCommand(command)
			case <-ctx.Done():
				return
			}
		}
	}()

	return commands
}

func handleCommand(command *Command) {
	leftSpeed = command.Left
	rightSpeed = command.Right
}
