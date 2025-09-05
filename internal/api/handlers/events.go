package handlers

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"sirteefyapps.com.ng/goroutines/internal/models"
// 	"sirteefyapps.com.ng/goroutines/internal/service"
// 	"sirteefyapps.com.ng/goroutines/pkg/errors"
// )

// type EventHandler struct {
// 	eventService service.EventService
// }

// func NewEventHandler(eventService service.EventService) *EventHandler {
// 	return &EventHandler{
// 		eventService: eventService,
// 	}
// }

// func (h *EventHandler) Create(c *gin.Context) {
// 	var newEvent models.Event
// 	if err := c.ShouldBindJSON(&newEvent); err != nil {
// 		result := errors.HandleError(
// 			func() (interface{}, error) {
// 				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
// 			},
// 			"binding JSON for new event",
// 		)
// 		result.RespondWithJSON(c)
// 		return
// 	}

// 	result := errors.HandleErrorWithStatusCode(
// 		func() (interface{}, error) {
// 			err := h.eventService.Create(&newEvent)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return newEvent, nil
// 		},
// 		"creating new event",
// 		http.StatusCreated,
// 	)
// 	result.RespondWithJSON(c)
// }

// func (h *EventHandler) GetAll(c *gin.Context) {
// 	result := errors.HandleError(
// 		func() (interface{}, error) {
// 			return h.eventService.GetAll()
// 		},
// 		"fetching all events",
// 	)
// 	result.RespondWithJSON(c)
// }

// func (h *EventHandler) GetByID(c *gin.Context) {
// 	eventID := c.Param("eventId")
	
// 	result := errors.HandleError(
// 		func() (interface{}, error) {
// 			return h.eventService.GetByID(eventID)
// 		},
// 		"fetching event by ID",
// 	)
// 	result.RespondWithJSON(c)
// }

// func (h *EventHandler) Update(c *gin.Context) {
// 	eventID := c.Param("eventId")
// 	var updateEvent models.Event
	
// 	if err := c.ShouldBindJSON(&updateEvent); err != nil {
// 		result := errors.HandleError(
// 			func() (interface{}, error) {
// 				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
// 			},
// 			"binding JSON for event update",
// 		)
// 		result.RespondWithJSON(c)
// 		return
// 	}

// 	result := errors.HandleError(
// 		func() (interface{}, error) {
// 			return h.eventService.Update(eventID, &updateEvent)
// 		},
// 		"updating event",
// 	)
// 	result.RespondWithJSON(c)
// }

// func (h *EventHandler) Delete(c *gin.Context) {
// 	eventID := c.Param("eventId")

// 	result := errors.HandleErrorWithStatusCode(
// 		func() (interface{}, error) {
// 			err := h.eventService.Delete(eventID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return gin.H{"message": "Event deleted successfully"}, nil
// 		},
// 		"deleting event",
// 		http.StatusOK,
// 	)
// 	result.RespondWithJSON(c)
// }