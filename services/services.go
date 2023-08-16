package services

import (
	"net/http"

	"github.com/boardware-cloud/model/core"

	"github.com/boardware-cloud/common/errors"

	"github.com/boardware-cloud/common/constants"
)

func ServiceBackward(service core.Service) Service {
	return Service{
		ID:          service.ID,
		Name:        service.Name,
		Title:       service.Title,
		Description: service.Description,
		Url:         service.Url,
		Type:        service.Type,
	}
}

type Service struct {
	ID          uint                  `json:"id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Title       string                `json:"title,omitempty"`
	Description string                `json:"description,omitempty"`
	Url         string                `json:"url,omitempty"`
	Type        constants.ServiceType `json:"type,omitempty"`
}

func CreateService(name, title, description, url string, serviceType constants.ServiceType) (Service, *errors.Error) {
	var services []core.Service
	DB.Find(&services, "name = ?", name)
	if len(services) > 0 {
		return Service{}, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.SERVICE_KEY_DUPLICATION_ERROR,
			Message:    "Service name duplication.",
		}
	}
	service := core.NewService(name, title, description, url, constants.ServiceType(serviceType))
	DB.Create(&service)
	return ServiceBackward(service), nil
}

func ListServices() List[Service] {
	var serviceList List[Service]
	var services []core.Service
	DB.Find(&services)
	for _, v := range services {
		serviceList.Data = append(serviceList.Data, ServiceBackward(v))
	}
	return serviceList
}
