package lock

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	log "github.com/net-auto/resourceManager/logging"
	"github.com/vektah/gqlparser/v2/ast"
)

// LockRequestInterceptor is a graphql interceptor that locks a resource pool when a mutation is performed on it.
type LockRequestInterceptor struct {
	lockingService LockingService
}

// NewLockRequestInterceptor creates a new LockRequestInterceptor.
func NewLockRequestInterceptor(lockingService LockingService) *LockRequestInterceptor {
	return &LockRequestInterceptor{lockingService: lockingService}
}

var _ interface {
	graphql.ResponseInterceptor
	graphql.OperationInterceptor
	graphql.HandlerExtension
} = &LockRequestInterceptor{}

// InterceptResponse intercepts the graphql response.
// If the response is a mutation, have poolId as argument and is trying to claim resource, it locks the resource pool.
// We are wrapping the next(ctx) in a lock/unlock pair to ensure that the resource pool is unlocked even if the response is nil.
// In this response interceptor we are providing also DB read/write safety access when there are multiple concurrent requests/commits.
func (l *LockRequestInterceptor) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	oc := graphql.GetOperationContext(ctx)

	if isMutation(oc) {
		if hasProperty(oc, "poolId", "ClaimResource") || hasProperty(oc, "poolId", "ClaimResourceWithAltId") {
			poolId := oc.Variables["poolId"].(string)
			l.lockingService.Acquire(poolId).Lock()

			select {
			case <-ctx.Done():
				l.lockingService.Unlock(poolId)

				// TODO finish this "DONE" to prevent allocation after request timeout
				log.Warn(ctx, "HTTP request finished earlier then allocation of resource")
				return graphql.ErrorResponse(ctx, "Request has been canceled")
			default:
				response := next(ctx)
				l.lockingService.Unlock(poolId)
				return response
			}

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

// isMutation checks if the operation is a mutation based on OperationContext.
func isMutation(oc *graphql.OperationContext) bool {
	return oc.Operation != nil && oc.Operation.Operation == ast.Mutation
}

// hasProperty checks if the operation has a property with the given  and operation testID.
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
