package notifications

import "context"

type Broadcaster interface {
	BroadcastNotification(ctx context.Context, n *Notification) error
}
