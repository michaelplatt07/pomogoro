package channels

import (
	// Internal imports
	"pomogoro/internal/messages"
)

type SongChannel struct {
	channel chan messages.ChannelMessage
}
