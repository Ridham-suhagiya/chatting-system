package service

import (
	"chatting-system-backend/database"
	"chatting-system-backend/model"
	"chatting-system-backend/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Messages interface {
	GetMessages(linkId string) (model.Messages, error)
	UpdateLinkMessages(linkId string, messages []byte) error
}

type message struct {
	DB *database.DB
}

type MessageContent struct {
	Content  string `json:"content"`
	Time     string `json:"time"`
	From     string `json:"from"`
	Username string `json:"username"`
}

func NewMessageService(db *database.DB) Messages {
	return &message{
		DB: db,
	}
}

func (s *message) GetMessages(linkId string) (model.Messages, error) {
	var messages = model.Messages{MessageContent: []byte{}}
	if err := s.DB.First(&messages, `link_id=?`, linkId).Error; err != nil {
		return messages, err
	}
	return messages, nil
}

func (s *message) InsetIntoMessages(linkId string, messages []byte) error {
	linkUUID, err := uuid.Parse(linkId)
	fmt.Println("some fonbfdionb", linkUUID, messages)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	return s.DB.Create(&model.Messages{LinkId: linkUUID, MessageContent: messages}).Error

}

func (s *message) UpdateLinkMessages(linkId string, messages []byte) error {
	linkUUID, err := uuid.Parse(linkId)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}

	// Attempt to update the existing record
	result := s.DB.Model(&model.Messages{}).Where("link_id = ?", linkUUID).Updates(map[string]interface{}{"message_content": messages})
	if result.Error != nil {
		fmt.Println("Update error:", result.Error)
		return result.Error
	}

	// Check if any rows were updated
	if result.RowsAffected == 0 {
		// No rows were updated, so insert a new record
		fmt.Println("No rows updated, inserting new record")
		err = s.InsetIntoMessages(linkId, messages)
		if err != nil {
			fmt.Println("Insert error:", err)
			return err
		}
	}

	return nil
}

type MessageInFile struct {
	MessageContent []string
	LinkId         string
}

var filePath string = "static/messages.jsonb"

func GetMessagesFromFile(linkId string) (map[string][]MessageContent, error) {
	utils.Mutex.Lock()
	defer utils.Mutex.Unlock()
	var messages map[string][]MessageContent
	data, err := os.ReadFile(filePath)
	if err != nil {
		return messages, fmt.Errorf("error reading file: %w", err)
	}
	if len(data) == 0 {
		return messages, nil
	}
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return messages, fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return messages, nil
}

func InsertMessageContentInFile(message []byte, linkId string) error {
	if linkId == "" || len(message) == 0 {
		return nil
	}

	var messageContent MessageContent

	err := json.Unmarshal(message, &messageContent)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	data, err := GetMessagesFromFile(linkId)

	if err != nil {
		return fmt.Errorf("something went wrong while fetching messags %w", err)
	}

	if len(data) == 0 {
		data = map[string][]MessageContent{linkId: {messageContent}}
	} else if len(data[linkId]) > 5 {

		db, err := database.GetDataBaseConnection()
		if err != nil {
			return fmt.Errorf("something went wrong while connecting db %w", err)
		}
		messageService := NewMessageService(db)
		messagesData, _ := messageService.GetMessages(linkId)
		completeMessage := append(data[linkId], messageContent)
		var unMarshalContent []MessageContent
		json.Unmarshal(messagesData.MessageContent, &unMarshalContent)
		messageArr := append(unMarshalContent, completeMessage...)
		marshalMessages, _ := json.Marshal(messageArr)
		messageService.UpdateLinkMessages(linkId, marshalMessages)
		data[linkId] = []MessageContent{}
	} else {
		data[linkId] = append(data[linkId], messageContent)
	}

	marshalData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	err = os.WriteFile(filePath, marshalData, 0644)
	if err != nil {
		fmt.Println("Something went wrong file writing into the file", err)
		return err
	}
	return nil
}

func EmptyMessageFile() error {
	err := os.WriteFile(filePath, []byte("{}"), 0644)
	if err != nil {
		return fmt.Errorf("error emptying the message file: %w", err)
	}
	return nil
}
