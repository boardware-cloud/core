package controllers

import (
	"net/http"

	api "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"

	"github.com/boardware-cloud/common/constants"

	"github.com/gin-gonic/gin"
)

type ServiceApi struct{}

var serviceApi ServiceApi

func (ServiceApi) CreateService(c *gin.Context) {
	middleware.IsRoot(c, func(c *gin.Context, account api.Account) {
		var createServicesRequest api.CreateServiceRequest
		if err := c.ShouldBindJSON(&createServicesRequest); err != nil {
			// TODO: Error message
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		service, err := core.CreateService(
			createServicesRequest.Name,
			createServicesRequest.Title,
			createServicesRequest.Description,
			createServicesRequest.Url,
			constants.ServiceType(createServicesRequest.Type),
		)
		if err != nil {
			err.GinHandler(c)
			return
		}
		c.JSON(http.StatusCreated, ServiceBackward(service))
	})
}

func (ServiceApi) ListServices(c *gin.Context) {
	c.JSON(http.StatusOK, ServiceListBackward(core.ListServices()))
}
