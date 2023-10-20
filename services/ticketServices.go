package services

import (
	"github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

func NewTicketService(db *gorm.DB) TicketService {
	return TicketService{ticketRepository: core.NewTicketRepository(db)}
}

type TicketService struct {
	ticketRepository core.TicketRepository
}
