package watermill

import "github.com/ThreeDotsLabs/watermill"

func NewLogger(debug, trace bool) watermill.LoggerAdapter {
	return watermill.NewStdLogger(debug, trace)
}
