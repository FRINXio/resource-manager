package lock

import (
	"context"
	"fmt"
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
		if isLockable(oc) {
			poolId, err := getArgument(oc, "poolId")
			if err != nil {
				log.Warn(ctx, "Unable to find poolId for query %s. Query will not be locked", oc.OperationName)
				return next(ctx)
			}

			l.lockingService.Acquire(*poolId).Lock()

			select {
			case <-ctx.Done():
				l.lockingService.Unlock(*poolId)
				log.Warn(ctx, "HTTP request finished before successfully locking. Skipping further executions")
				return graphql.ErrorResponse(ctx, "Request has been canceled")
			default:
				response := next(ctx)
				l.lockingService.Unlock(*poolId)
				return response
			}
		}
	}

	return next(ctx)
}

func isLockable(oc *graphql.OperationContext) bool {
	return matchesNameAndArgument(oc, "poolId", "ClaimResource") ||
		matchesNameAndArgument(oc, "poolId", "ClaimResourceWithAltId")
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

// matchesNameAndArgument checks if the operation has a property with the given  and operation testID.
func matchesNameAndArgument(oc *graphql.OperationContext, propertyName string, operationName string) bool {
	for _, selection := range oc.Operation.SelectionSet {
		if field, ok := selection.(*ast.Field); ok {
			if operationName == field.Name && field.Arguments.ForName(propertyName) != nil {
				return true
			}
		}
	}

	return false
}
func getArgument(oc *graphql.OperationContext, propertyName string) (*string, error) {
	err := fmt.Errorf("cannot find %s argument", propertyName)
	for _, selection := range oc.Operation.SelectionSet {
		if field, ok := selection.(*ast.Field); ok {
			if field.Arguments.ForName(propertyName) != nil {
				argMap := selection.(*ast.Field).ArgumentMap(oc.Variables)
				propertyValueAsStr := fmt.Sprintf("%v", argMap[propertyName])
				return &propertyValueAsStr, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, err
}
