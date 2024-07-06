//go:generate easyjson ./message.go
package notify

import "time"

//easyjson:json
type Message struct {
	Type     string    `json:"type"`
	SlotID   int       `json:"slotId"`
	BannerID int       `json:"bannerId"`
	GroupID  int       `json:"groupId"`
	Time     time.Time `json:"time"`
}
