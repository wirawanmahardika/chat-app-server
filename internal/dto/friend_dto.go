package dto

type FriendResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Avatar      *string `json:"avatar"`
	LastMessage *string `json:"lastMessage"`
	Online      bool    `json:"online"`
	LastSeen    string  `json:"lastSeen"`
}

type FriendRequestResponse struct {
	ID         string  `json:"id"`
	FromID     string  `json:"fromId"`
	FromName   string  `json:"fromName"`
	FromAvatar *string `json:"fromAvatar"`
	ToID       string  `json:"toId"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"createdAt"`
}

type SendFriendRequest struct {
	FriendID string `json:"friendId" validate:"required"`
}

type RespondFriendRequest struct {
	RequestID string `json:"requestId" validate:"required"`
}
