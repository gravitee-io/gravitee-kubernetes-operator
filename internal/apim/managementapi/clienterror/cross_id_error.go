package clienterror

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

type CrossIdNotFoundError struct {
	CrossId string
}

func (e *CrossIdNotFoundError) Error() string {
	return "No API found for CrossId " + e.CrossId
}

type CrossIdMultipleFoundError struct {
	CrossId string
	Apis    []model.ApiListItem
}

func (e *CrossIdMultipleFoundError) Error() string {
	return fmt.Sprintf("Multiple APIs found for CrossId %s. (%d APIs found)", e.CrossId, len(e.Apis))
}

type CrossIdUnauthorizedError struct {
	CrossId string
}

func (e CrossIdUnauthorizedError) Error() string {
	return "Unauthorized error for CrossId " + e.CrossId
}
