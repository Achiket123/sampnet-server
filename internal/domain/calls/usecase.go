package calls

import (
	"context"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	HandleRoom(ctx context.Context, roomID string, ws *websocket.Conn)
}
