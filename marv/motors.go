package marv

import "context"

// Command NEEDSCOMMENT
type Command struct {
	Left  float64
	Right float64
}

// StartMotors NEEDSCOMMENT
func StartMotors(ctx context.Context) chan<- *Command {
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
	// leftSpeed = command.Left
	// rightSpeed = command.Right
}
