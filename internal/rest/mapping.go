package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/oligzeev/host-manager/internal/domain"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	ParamId = "id"
)

type MappingRestHandler struct {
	mappingService domain.MappingService
}

func NewMappingRestHandler(mappingService domain.MappingService) *MappingRestHandler {
	return &MappingRestHandler{mappingService: mappingService}
}

func (h MappingRestHandler) Register(router *gin.Engine) {
	group := router.Group("/mapping")
	group.GET("/", h.getMappings)
	group.GET("/:"+ParamId, h.getMappingById)
}

// GetMappings godoc
// @Summary Get Mappings
// @Description Method to get all mappings
// @Tags Mapping
// @Accept json
// @Produce json
// @Success 200 {array} domain.Mapping
// @Failure 500 {object} rest.Error
// @Router /mapping [get]
func (h MappingRestHandler) getMappings(c *gin.Context) {
	var results []domain.Mapping
	if err := h.mappingService.GetAll(c.Request.Context(), &results); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, results)
}

// GetMappingById godoc
// @Summary Get Mapping by Id
// @Description Method to get mapping by id
// @Tags Mapping
// @Accept json
// @Produce json
// @Param id path string true "Mapping Id"
// @Success 200 {object} domain.Mapping
// @Failure 500 {object} rest.Error
// @Router /mapping/{id} [get]
func (h MappingRestHandler) getMappingById(c *gin.Context) {
	id := c.Param(ParamId)
	var result domain.Mapping
	if err := h.mappingService.GetById(c.Request.Context(), id, &result); err != nil {
		log.Error(err)
		if domain.ECode(err) == domain.ErrNotFound {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, E(err))
		return
	}
	c.JSON(http.StatusOK, result)
}
