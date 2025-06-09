package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/geico-private/gts/clients/go/tlm/lg"

	"serverless_platform/cmd/serverless_controller/configs"
	"serverless_platform/internal/serverless_controller/managers"
	"serverless_platform/internal/serverless_controller/models"
	"serverless_platform/internal/serverless_controller/util"
	"serverless_platform/internal/shared/enums"
	sharedModels "serverless_platform/internal/shared/models"
	"serverless_platform/internal/shared/utils"
)

// Update the orchestrateHandler struct to include the validation error handler
type orchestrateHandler struct {
	config        configs.Config
	pm            managers.ProvisionManager
	alm           managers.AlmActivityManager
	pubm          managers.PublishManager
	validationHandler *ValidationErrorHandler // Add this new field
}

// Update the constructor to initialize the validation error handler
func NewOrchestrateHandler(config configs.Config, pm managers.ProvisionManager,
	alm managers.AlmActivityManager, pubm managers.PublishManager) (OrchestrateHandler, error) {

	// Initialize the validation error handler
	validationHandler, err := NewValidationErrorHandler()
	if err != nil {
		return nil, util.WrapErr(err, "failed to initialize validation handler")
	}

	return &orchestrateHandler{
		config:            config,
		pm:                pm,
		alm:               alm,
		pubm:              pubm,
		validationHandler: validationHandler, // Store the validation handler
	}, nil
}

// Updated Orchestrate method with enhanced validation error handling
func (o *orchestrateHandler) Orchestrate(c *gin.Context) {
	var almRequest *sharedModels.ALMRequest
	var err error
	ctx := c.Request.Context()
	logger := lg.FromContext(ctx).WithGroup("orchestrate")
	logger.InfoContext(ctx, "Orchestrate endpoint hit")

	traceId := util.GetTraceIdFromContext(ctx)

	// Bind the request to the OrchestrationRequest struct
	var request models.OrchestrationRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		logger.ErrorContext(ctx, BIND_FAILED, lg.Error, err)
		err = util.WrapErr(err, BIND_FAILED)
		c.JSON(http.StatusBadRequest, toOrchestrateResponse(request, traceId, almRequest, err))
		return
	}

	// Enhanced validation with user-friendly error messages
	if validationErrors := o.validationHandler.ValidateStruct(request); len(validationErrors) > 0 {
		logger.ErrorContext(ctx, "Validation failed", "errors", validationErrors)

		// Create a user-friendly error response
		errorResponse := o.validationHandler.FormatValidationErrors(validationErrors)

		// Create a custom error for internal logging
		validationErr := errors.New("validation failed: " + strings.Join(validationErrors, "; "))

		// Return the formatted error response to the client
		response := toOrchestrateResponse(request, traceId, almRequest, validationErr)

		// Enhance the response with the detailed validation errors
		if responseMap, ok := response.(map[string]interface{}); ok {
			responseMap["validationErrors"] = errorResponse
		}

		c.JSON(http.StatusBadRequest, response)
		return
	}

	logger.With("CorrelationId", request.Cicd.CorrelationId)
	// Inject the updated logger back into the context
	ctx = lg.NewContext(ctx, logger)

	if request.Action == enums.ActionEmpty {
		logger.WarnContext(ctx, "Action is empty, defaulting to deploy")
		request.Action = enums.ActionDeploy
	}

	ownerId := util.GetOwnerIdFromContext(c)
	if ownerId == "" {
		logger.ErrorContext(ctx, AUTHZ_FAILED)
		c.JSON(http.StatusForbidden, toOrchestrateResponse(request, traceId,
			almRequest, errors.New(AUTHZ_FAILED)))
		return
	}

	// Set the ownerId in the Request context
	ctx = util.SetOwnerIdToRequestContext(ctx, ownerId)

	almRequest, err = o.pm.Provision(ctx,
		managers.ProvisionManagerParameters{IncomingRequest: request, OwnerId: ownerId})
	if err != nil {
		logger.ErrorContext(ctx, PROVISION_FAILED, lg.Error, err)
		err = util.WrapErr(err, PROVISION_FAILED)
		c.JSON(http.StatusBadRequest, toOrchestrateResponse(request, traceId, almRequest, err))
		return
	}

	err = o.alm.TriggerActivity(ctx, request, *almRequest)

	if err != nil {
		logger.ErrorContext(ctx, ALM_PUBLISH_FAILED, lg.Error, err)
		err = util.WrapErr(err, ALM_PUBLISH_FAILED)
		c.JSON(http.StatusBadRequest, toOrchestrateResponse(request, traceId, almRequest, err))
		return
	}

	c.JSON(http.StatusOK, toOrchestrateResponse(request, traceId, almRequest, err))
}

// Enhanced error response structure
func toOrchestrateResponse(request models.OrchestrationRequest,
	traceId string, almRequest *sharedModels.ALMRequest, err error) interface{} {

	baseResponse := map[string]interface{}{
		"OrchestrationRequest": request,
		"DeploymentActivityId": traceId,
		"Error":                getErrorResponse(err),
	}

	if utils.IsLocalEnv() {
		baseResponse["Dynamic_Internal_Local"] = almRequest
	}

	return baseResponse
}

// Enhanced error response that works well with validation errors
func getErrorResponse(err error) map[string]interface{} {
	if err == nil {
		return nil
	}

	return map[string]interface{}{
		"Description": err.Error(),
	}
}
