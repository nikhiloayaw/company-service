package handler

import (
	"api-gateway/pkg/api/handler/interfaces"
	"api-gateway/pkg/api/handler/request"
	"api-gateway/pkg/api/handler/response"
	clientInterfaces "api-gateway/pkg/client/interfaces"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type companyHandler struct {
	client clientInterfaces.CompanyServiceClient
}

func NewCompanyHandler(client clientInterfaces.CompanyServiceClient) interfaces.CompanyHandler {

	return &companyHandler{
		client: client,
	}
}

func (c *companyHandler) Create(ctx *gin.Context) {

	// create request body which have binding validation tags
	var body request.CompanyRequest

	// use gin binding which will also check the validation and return error if any.
	if err := ctx.ShouldBindJSON(&body); err != nil {

		response := response.ErrorResponse(http.StatusBadRequest, BindErrorMessage, err, body)
		ctx.JSON(http.StatusBadRequest, response)

		return
	}

	companyData, err := c.client.Create(body)

	if err != nil {
		response := response.ErrorResponse(http.StatusInternalServerError, "internal server error", err, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	var company response.Company

	if err := json.Unmarshal(companyData, &company); err != nil {
		response := response.ErrorResponse(http.StatusInternalServerError, "internal server error", err, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusCreated, "successfully company details created", company)
	ctx.JSON(http.StatusCreated, response)
}
