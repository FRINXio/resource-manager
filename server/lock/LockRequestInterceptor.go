package lock

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

type LockRequestInterceptor struct {
	lockingService *LockingService
}

func NewLockRequestInterceptor(lockingService *LockingService) *LockRequestInterceptor {
	return &LockRequestInterceptor{lockingService: lockingService}
}

var _ interface {
	graphql.ResponseInterceptor
	graphql.OperationInterceptor
	graphql.HandlerExtension
} = &LockRequestInterceptor{}

func (l *LockRequestInterceptor) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	oc := graphql.GetOperationContext(ctx)

	if isMutation(oc) {
		if hasProperty(oc, "poolId", "ClaimResource") || hasProperty(oc, "poolId", "ClaimResourceWithAltId") {
			poolId := oc.Variables["poolId"].(string)
			_, err := l.lockingService.Lock(poolId)
			if err != nil {
				return nil
			}
			response := next(ctx)
			l.lockingService.Unlock(poolId)
			return response
		}
	}

	return next(ctx)
}

func (l *LockRequestInterceptor) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	return next(ctx)
}

func (l *LockRequestInterceptor) ExtensionName() string {
	return "LockRequestInterceptor"
}

func (l *LockRequestInterceptor) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func isMutation(oc *graphql.OperationContext) bool {
	return oc.Operation != nil && oc.Operation.Operation == ast.Mutation
}

func hasProperty(oc *graphql.OperationContext, propertyName string, operationName string) bool {
	for _, selection := range oc.Operation.SelectionSet {
		if field, ok := selection.(*ast.Field); ok {
			if operationName == field.Name && field.Arguments.ForName(propertyName) != nil {
				return true
			}
		}
	}

	return false
}
