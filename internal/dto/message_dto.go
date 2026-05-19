package dto

type MessageResponse struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversationId"`
	SenderID       string `json:"senderId"`
	ReceiverID     string `json:"receiverId"`
	Text           string `json:"text"`
	Read           bool   `json:"read"`
	Delivered      bool   `json:"delivered"`
	CreatedAt      string `json:"createdAt"`
}

type SendMessageRequest struct {
	ReceiverID string `json:"receiverId" validate:"required"`
	Text       string `json:"text" validate:"required,min=1,max=2000"`
}

type ConversationResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Avatar          *string `json:"avatar"`
	LastMessage     *string `json:"lastMessage"`
	Online          bool    `json:"online"`
	LastSeen        string  `json:"lastSeen"`
	LastMessageTime *string `json:"lastMessageTime"`
}

type PaginationInfo struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int64 `json:"total"`
}

type MessagesWithPagination struct {
	Success    bool              `json:"success"`
	Data       []MessageResponse `json:"data"`
	Pagination PaginationInfo    `json:"pagination"`
}
