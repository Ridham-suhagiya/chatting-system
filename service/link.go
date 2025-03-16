package service

import (
	"chatting-system-backend/database"
	"chatting-system-backend/model"
	"chatting-system-backend/utils"
)

type LinkService interface {
	CreateLink(link *model.ChatLinks) error
	GetLinkDetailsUsingLinkCode(linkCode string) (model.ChatLinks, error)
	GetLinkDetailUsingLinkId(id string) (model.ChatLinks, error)
}

type linkService struct {
	DB *database.DB
}

func NewlinkService(db *database.DB) LinkService {
	return &linkService{DB: db}
}

func (s *linkService) CreateLink(link *model.ChatLinks) error {
	link.LinkCode = utils.RandomCodeGenrator()
	return s.DB.Create(&link).Error
}

func (s *linkService) GetLinkDetailsUsingLinkCode(linkCode string) (model.ChatLinks, error) {
	var link model.ChatLinks
	if err := s.DB.First(&link, `link_code=?`, linkCode).Error; err != nil {
		return link, err
	}
	return link, nil
}

func (s *linkService) GetLinkDetailUsingLinkId(id string) (model.ChatLinks, error) {
	var link model.ChatLinks
	if err := s.DB.First(&link, `id=?`, id).Error; err != nil {
		return link, err
	}
	return link, nil
}
