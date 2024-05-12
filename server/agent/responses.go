package agent

import (
	"context"
)

type ResponsesPipe chan string

func (rp ResponsesPipe) Send(ctx context.Context, msg string) {
	select {
	case <-ctx.Done():
	default:
		rp <- msg
	}
}

func (rp ResponsesPipe) SendError(ctx context.Context, err error) {
	select {
	case <-ctx.Done():
	default:
		// todo: recognise an error and respond accordingly
		rp <- err.Error()
	}
}
