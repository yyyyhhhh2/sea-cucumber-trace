package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sea-cucumber-trace/backend/internal/model"
	"sea-cucumber-trace/backend/internal/service"
)

func (h *Handler) ListBatches(c *gin.Context) {
	cl := claimsFromCtx(c)
	list, err := h.svc.ListBatches(cl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

type createBatchBody struct {
	BatchNo        string     `json:"batchNo" binding:"required"`
	OrgID          uint       `json:"orgId" binding:"required"`
	ProductName    string     `json:"productName"`
	FarmBase       string     `json:"farmBase"`
	Quality        string     `json:"quality"`
	CatchDate      *time.Time `json:"catchDate"`
	BreedArea      string     `json:"breedArea"`
	BreedStartDate *time.Time `json:"breedStartDate"`
	Spec           string     `json:"spec"`
	Quantity       string     `json:"quantity"`
	ExtraJSON      string     `json:"extraJson"`
}

func (h *Handler) CreateBatch(c *gin.Context) {
	cl := claimsFromCtx(c)
	var body createBatchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	b, err := h.svc.CreateBatch(cl, service.CreateBatchInput{
		BatchNo:        body.BatchNo,
		ProductName:    body.ProductName,
		FarmBase:       body.FarmBase,
		Quality:        body.Quality,
		CatchDate:      body.CatchDate,
		BreedArea:      body.BreedArea,
		BreedStartDate: body.BreedStartDate,
		Spec:           body.Spec,
		Quantity:       body.Quantity,
		ExtraJSON:      body.ExtraJSON,
		OrgID:          body.OrgID,
	})
	if err != nil {
		status := http.StatusForbidden
		if errors.Is(err, service.ErrBatchNoRequired) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *Handler) GetBatch(c *gin.Context) {
	cl := claimsFromCtx(c)
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
		return
	}
	b, err := h.svc.GetBatch(cl, uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

type addEventBody struct {
	Stage        model.TraceStage `json:"stage" binding:"required"`
	Title        string           `json:"title" binding:"required"`
	DetailJSON   string           `json:"detailJson"`
	Location     string           `json:"location"`
	OperatorName string           `json:"operatorName"`
	EvidenceURLs string           `json:"evidenceUrls"`
	OccurredAt   time.Time        `json:"occurredAt"`
}

func (h *Handler) AddEvent(c *gin.Context) {
	cl := claimsFromCtx(c)
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
		return
	}
	var body addEventBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if body.OccurredAt.IsZero() {
		body.OccurredAt = time.Now()
	}
	ev, cr, err := h.svc.AddEvent(cl, uint(id64), service.AddEventInput{
		Stage:        body.Stage,
		Title:        body.Title,
		DetailJSON:   body.DetailJSON,
		Location:     body.Location,
		OperatorName: body.OperatorName,
		EvidenceURLs: body.EvidenceURLs,
		OccurredAt:   body.OccurredAt,
	})
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"event": ev, "chain": cr})
}

func (h *Handler) Timeline(c *gin.Context) {
	cl := claimsFromCtx(c)
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
		return
	}
	items, err := h.svc.TimelineByBatchID(cl, uint(id64))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) PublicTimeline(c *gin.Context) {
	no := c.Param("batchNo")
	items, batch, err := h.svc.PublicTimeline(no)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"batch": batch, "items": items})
}

func (h *Handler) Verify(c *gin.Context) {
	no := c.Param("batchNo")
	res, err := h.svc.Verify(no)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": res})
}
