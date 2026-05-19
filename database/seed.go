package database

import (
	"chatapp/database/model"
	"log"
	"time"

	"gorm.io/gorm"
)

func strPtr(s string) *string {
	return &s
}

type SeedMessage struct {
	SenderEmail string
	Text        string
	TimeOffset  time.Duration
	Read        bool
}

type SeedConvo struct {
	User1Email string
	User2Email string
	Messages   []SeedMessage
}

func Seed(db *gorm.DB) error {
	log.Println("Clearing existing data...")

	// Delete in dependency order
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Message{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ConversationParticipant{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Conversation{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Friendship{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.User{}).Error; err != nil {
		return err
	}

	log.Println("Creating mock users...")
	users := []model.User{
		{
			Name:     "Alice Smith",
			Email:    "alice@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=Alice"),
			Online:   true,
		},
		{
			Name:     "Bob Jones",
			Email:    "bob@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=Bob"),
			Online:   true,
		},
		{
			Name:     "Charlie Brown",
			Email:    "charlie@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=Charlie"),
			Online:   false,
			LastSeen: time.Now().Add(-5 * time.Minute),
		},
		{
			Name:     "David Miller",
			Email:    "david@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=David"),
			Online:   false,
			LastSeen: time.Now().Add(-1 * time.Hour),
		},
		{
			Name:     "Emma Watson",
			Email:    "emma@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=Emma"),
			Online:   true,
		},
		{
			Name:     "Frank Castle",
			Email:    "frank@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=Frank"),
			Online:   false,
			LastSeen: time.Now().Add(-12 * time.Hour),
		},
		{
			Name:     "Grace Hopper",
			Email:    "grace@example.com",
			Password: "password123",
			Avatar:   strPtr("https://api.dicebear.com/7.x/adventurer/svg?seed=Grace"),
			Online:   false,
			LastSeen: time.Now().Add(-2 * time.Hour),
		},
	}

	emailToID := make(map[string]string)
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
		emailToID[users[i].Email] = users[i].ID
	}
	log.Printf("Successfully created %d users.", len(users))

	log.Println("Creating mock friendships...")
	friendships := []model.Friendship{
		{SenderID: emailToID["alice@example.com"], ReceiverID: emailToID["bob@example.com"], Status: "accepted"},
		{SenderID: emailToID["alice@example.com"], ReceiverID: emailToID["charlie@example.com"], Status: "accepted"},
		{SenderID: emailToID["alice@example.com"], ReceiverID: emailToID["emma@example.com"], Status: "accepted"},
		{SenderID: emailToID["bob@example.com"], ReceiverID: emailToID["charlie@example.com"], Status: "accepted"},
		{SenderID: emailToID["bob@example.com"], ReceiverID: emailToID["emma@example.com"], Status: "accepted"},
		{SenderID: emailToID["charlie@example.com"], ReceiverID: emailToID["david@example.com"], Status: "accepted"},
		{SenderID: emailToID["david@example.com"], ReceiverID: emailToID["emma@example.com"], Status: "pending"},
		{SenderID: emailToID["frank@example.com"], ReceiverID: emailToID["alice@example.com"], Status: "pending"},
		{SenderID: emailToID["grace@example.com"], ReceiverID: emailToID["bob@example.com"], Status: "declined"},
	}
	for i := range friendships {
		if err := db.Create(&friendships[i]).Error; err != nil {
			return err
		}
	}
	log.Printf("Successfully created %d friendships.", len(friendships))

	log.Println("Creating mock conversations and messages...")
	seedConvos := []SeedConvo{
		{
			User1Email: "alice@example.com",
			User2Email: "bob@example.com",
			Messages: []SeedMessage{
				{SenderEmail: "alice@example.com", Text: "Hi Bob! How are you?", TimeOffset: -2 * time.Hour, Read: true},
				{SenderEmail: "bob@example.com", Text: "Hey Alice! I'm doing great, how about you?", TimeOffset: -1*time.Hour - 50*time.Minute, Read: true},
				{SenderEmail: "alice@example.com", Text: "I'm doing well too! Are you free for lunch today?", TimeOffset: -1*time.Hour - 30*time.Minute, Read: true},
				{SenderEmail: "bob@example.com", Text: "Yes, I am! The usual place at 12:30?", TimeOffset: -1*time.Hour - 15*time.Minute, Read: true},
				{SenderEmail: "alice@example.com", Text: "Perfect, see you there!", TimeOffset: -1 * time.Hour, Read: true},
				{SenderEmail: "bob@example.com", Text: "Awesome, see ya!", TimeOffset: -45 * time.Minute, Read: true},
			},
		},
		{
			User1Email: "alice@example.com",
			User2Email: "charlie@example.com",
			Messages: []SeedMessage{
				{SenderEmail: "charlie@example.com", Text: "Hi Alice, did you finish the design draft?", TimeOffset: -5 * time.Hour, Read: true},
				{SenderEmail: "alice@example.com", Text: "Hey Charlie, yes I did. I've sent it to your email.", TimeOffset: -4 * time.Hour - 45*time.Minute, Read: true},
				{SenderEmail: "charlie@example.com", Text: "Great! Let me review it. Thanks!", TimeOffset: -4 * time.Hour - 30*time.Minute, Read: true},
				{SenderEmail: "alice@example.com", Text: "No problem. Let me know if you need any changes.", TimeOffset: -4 * time.Hour, Read: true},
				{SenderEmail: "charlie@example.com", Text: "Looks good at first glance, just need to tweak the colors slightly.", TimeOffset: -3 * time.Hour, Read: true},
			},
		},
		{
			User1Email: "alice@example.com",
			User2Email: "emma@example.com",
			Messages: []SeedMessage{
				{SenderEmail: "emma@example.com", Text: "Hey Alice, are we meeting up this weekend?", TimeOffset: -1 * time.Hour, Read: true},
				{SenderEmail: "alice@example.com", Text: "Hi Emma! Yes, let's hang out on Saturday afternoon.", TimeOffset: -40 * time.Minute, Read: true},
				{SenderEmail: "emma@example.com", Text: "Sounds good to me! Let's go to that new coffee shop.", TimeOffset: -30 * time.Minute, Read: true},
				{SenderEmail: "alice@example.com", Text: "Oh, the one downtown? Count me in!", TimeOffset: -10 * time.Minute, Read: false},
			},
		},
		{
			User1Email: "bob@example.com",
			User2Email: "charlie@example.com",
			Messages: []SeedMessage{
				{SenderEmail: "bob@example.com", Text: "Hey Charlie, did you watch the game last night?", TimeOffset: -10 * time.Hour, Read: true},
				{SenderEmail: "charlie@example.com", Text: "Yeah! What a finish! Unbelievable goal.", TimeOffset: -9 * time.Hour - 50*time.Minute, Read: true},
				{SenderEmail: "bob@example.com", Text: "I know right? I thought they were going to lose.", TimeOffset: -9 * time.Hour - 45*time.Minute, Read: true},
			},
		},
		{
			User1Email: "bob@example.com",
			User2Email: "emma@example.com",
			Messages: []SeedMessage{
				{SenderEmail: "emma@example.com", Text: "Bob, can you share the project repository link?", TimeOffset: -3 * time.Hour, Read: true},
				{SenderEmail: "bob@example.com", Text: "Sure, here you go: github.com/wirawanmahardika/chat-app", TimeOffset: -2 * time.Hour - 55*time.Minute, Read: true},
				{SenderEmail: "emma@example.com", Text: "Thanks a lot! Lifesaver.", TimeOffset: -2 * time.Hour - 50*time.Minute, Read: true},
			},
		},
	}

	for _, sc := range seedConvos {
		u1ID := emailToID[sc.User1Email]
		u2ID := emailToID[sc.User2Email]

		if u1ID == "" || u2ID == "" {
			continue
		}

		convo := model.Conversation{
			Type: "private",
		}
		if err := db.Create(&convo).Error; err != nil {
			return err
		}

		// Create participants
		p1 := model.ConversationParticipant{
			ConversationID: convo.ID,
			UserID:         u1ID,
		}
		p2 := model.ConversationParticipant{
			ConversationID: convo.ID,
			UserID:         u2ID,
		}
		if err := db.Create(&p1).Error; err != nil {
			return err
		}
		if err := db.Create(&p2).Error; err != nil {
			return err
		}

		var lastMsgID *string
		var lastMsgTime time.Time
		now := time.Now()

		for _, sm := range sc.Messages {
			senderID := emailToID[sm.SenderEmail]
			var receiverID string
			if senderID == u1ID {
				receiverID = u2ID
			} else {
				receiverID = u1ID
			}

			msgCreatedAt := now.Add(sm.TimeOffset)
			msg := model.Message{
				ConversationID: convo.ID,
				SenderID:       senderID,
				ReceiverID:     receiverID,
				Text:           sm.Text,
				Read:           sm.Read,
				Delivered:      true,
				CreatedAt:      msgCreatedAt,
			}
			if err := db.Create(&msg).Error; err != nil {
				return err
			}
			lastMsgID = &msg.ID
			lastMsgTime = msgCreatedAt
		}

		if lastMsgID != nil {
			updates := map[string]interface{}{
				"last_message_id": lastMsgID,
				"updated_at":      lastMsgTime,
			}
			if err := db.Model(&convo).Updates(updates).Error; err != nil {
				return err
			}
		}
	}
	log.Println("Conversations and messages seeded successfully.")

	return nil
}
