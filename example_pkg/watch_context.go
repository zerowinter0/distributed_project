package example_pkg

type WatchContext struct {
	channel  chan string
	messages []string
}

var watchCtx *WatchContext = nil

func setWatchCtx(ctx *WatchContext) {
	watchCtx = ctx
}

func watchAppendMessage(m string) {
	if watchCtx != nil {
		watchCtx.channel <- m
	}
}
