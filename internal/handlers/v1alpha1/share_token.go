package v1alpha1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apiServer "github.com/kubev2v/migration-planner/internal/api/server"
	"github.com/kubev2v/migration-planner/internal/auth"
	"github.com/kubev2v/migration-planner/internal/handlers/v1alpha1/mappers"
	"github.com/kubev2v/migration-planner/internal/service"
)

// (POST /api/v1/assessments/{id}/share-token)
func (s *ServiceHandler) CreateShareToken(ctx context.Context, request apiServer.CreateShareTokenRequestObject) (apiServer.CreateShareTokenResponseObject, error) {
	if errType, errMessage := s.mustAccessAssessment(ctx, request.Id); errType != errorNone {
		switch errType {
		case error403:
			return apiServer.CreateShareToken403JSONResponse{Message: errMessage}, nil
		case error404:
			return apiServer.CreateShareToken404JSONResponse{Message: errMessage}, nil
		case error500:
			return apiServer.CreateShareToken500JSONResponse{Message: errMessage}, nil
		default:
		}
	}

	// Create or get existing share token for this assessment
	shareToken, err := s.shareTokenSrv.CreateShareToken(ctx, request.Id)
	if err != nil {
		switch err.(type) {
		case *service.ErrResourceNotFound:
			return apiServer.CreateShareToken404JSONResponse{Message: err.Error()}, nil
		case *service.ErrSourceNoInventory:
			// Check if it's an inventory validation error
			return apiServer.CreateShareToken400JSONResponse{Message: err.Error()}, nil
		default:
			return apiServer.CreateShareToken500JSONResponse{Message: err.Error()}, nil
		}
	}

	return apiServer.CreateShareToken200JSONResponse(mappers.ShareTokenToApi(*shareToken)), nil
}

// (DELETE /api/v1/assessments/{id}/share-token)
func (s *ServiceHandler) DeleteShareToken(ctx context.Context, request apiServer.DeleteShareTokenRequestObject) (apiServer.DeleteShareTokenResponseObject, error) {
	if errType, errMessage := s.mustAccessAssessment(ctx, request.Id); errType != errorNone {
		switch errType {
		case error403:
			return apiServer.DeleteShareToken403JSONResponse{Message: errMessage}, nil
		case error404:
			return apiServer.DeleteShareToken404JSONResponse{Message: errMessage}, nil
		case error500:
			return apiServer.DeleteShareToken500JSONResponse{Message: errMessage}, nil
		default:
		}
	}

	// Delete share token for this assessment
	err := s.shareTokenSrv.DeleteShareToken(ctx, request.Id)
	if err != nil {
		return apiServer.DeleteShareToken500JSONResponse{Message: err.Error()}, nil
	}

	return apiServer.DeleteShareToken200JSONResponse{
		Message: func() *string { s := "Share token deleted successfully"; return &s }(),
		Status:  func() *string { s := "Success"; return &s }(),
	}, nil
}

// (GET /api/v1/assessments/{id}/share-token)
func (s *ServiceHandler) GetShareToken(ctx context.Context, request apiServer.GetShareTokenRequestObject) (apiServer.GetShareTokenResponseObject, error) {
	if errType, errMessage := s.mustAccessAssessment(ctx, request.Id); errType != errorNone {
		switch errType {
		case error403:
			return apiServer.GetShareToken403JSONResponse{Message: errMessage}, nil
		case error404:
			return apiServer.GetShareToken404JSONResponse{Message: errMessage}, nil
		case error500:
			return apiServer.GetShareToken500JSONResponse{Message: errMessage}, nil
		default:
		}
	}

	// Get share token for this assessment
	shareToken, err := s.shareTokenSrv.GetShareTokenBySourceID(ctx, request.Id)
	if err != nil {
		switch err.(type) {
		case *service.ErrResourceNotFound:
			return apiServer.GetShareToken404JSONResponse{Message: "share token not found"}, nil
		default:
			return apiServer.GetShareToken500JSONResponse{Message: err.Error()}, nil
		}
	}

	return apiServer.GetShareToken200JSONResponse(mappers.ShareTokenToApi(*shareToken)), nil
}

// (GET /api/v1/assessments/share-tokens)
func (s *ServiceHandler) ListShareTokens(ctx context.Context, request apiServer.ListShareTokensRequestObject) (apiServer.ListShareTokensResponseObject, error) {
	user := auth.MustHaveUser(ctx)

	// List all share tokens for the user's organization
	shareTokens, err := s.shareTokenSrv.ListShareTokens(ctx, user.Organization)
	if err != nil {
		return apiServer.ListShareTokens500JSONResponse{Message: err.Error()}, nil
	}

	return apiServer.ListShareTokens200JSONResponse(mappers.ShareTokenListToApi(shareTokens)), nil
}

// (GET /api/v1/assessments/shared/{token})
func (s *ServiceHandler) GetSharedAssessment(ctx context.Context, request apiServer.GetSharedAssessmentRequestObject) (apiServer.GetSharedAssessmentResponseObject, error) {
	// Get source by token (no authentication required)
	// TODO nargaman change to GetAssessmentByToken
	source, err := s.shareTokenSrv.GetSourceByToken(ctx, request.Token)
	if err != nil {
		switch err.(type) {
		case *service.ErrResourceNotFound:
			return apiServer.GetSharedAssessment404JSONResponse{Message: err.Error()}, nil
		default:
			return apiServer.GetSharedAssessment500JSONResponse{Message: err.Error()}, nil
		}
	}

	return apiServer.GetSharedAssessment200JSONResponse(mappers.SourceToApi(*source)), nil
}

type ErrorType int

const (
	error404 ErrorType = iota
	error403
	error500
	errorNone
)

// mustAccessAssessment checks if source exists and user has access
func (s *ServiceHandler) mustAccessAssessment(ctx context.Context, assessmentID uuid.UUID) (ErrorType, string) {
	// Check if source exists and user has access
	// TODO nargaman change to GetAssessment
	assessment, err := s.sourceSrv.GetSource(ctx, assessmentID)
	if err != nil {
		switch err.(type) {
		case *service.ErrResourceNotFound:
			return error404, err.Error()
		default:
			return error500, err.Error()
		}
	}

	user := auth.MustHaveUser(ctx)
	if user.Organization != assessment.OrgID {
		return error403, fmt.Sprintf("forbidden access to assessment %q", assessmentID)
	}
	return errorNone, ""
}
