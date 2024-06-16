package messages

type ChannelMessage struct {
	SongFinished bool
	SongSkipped  bool
	SongStopped  bool
	SongPaused   bool
	SongResumed  bool
}
