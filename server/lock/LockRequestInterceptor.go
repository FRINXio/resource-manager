package lock

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	log "github.com/net-auto/resourceManager/logging"
	"github.com/pkg/errors"
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

	fmt.Println("response interceptor")
	fmt.Println(oc.Operation.Directives)

	if !isMutation(oc) {
		if hasProperty(oc, "poolId") && (oc.Operation.Name == "ClaimResource" || oc.Operation.Name == "ClaimResourceWithAltId") {
			poolId := oc.Variables["poolId"].(string)
			log.Error(ctx, errors.Errorf("Locking pool"), "poolId", poolId)
			_, err := l.lockingService.Lock(poolId)
			if err != nil {
				return nil
			}
			l.lockingService.Unlock(poolId)
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

func hasProperty(oc *graphql.OperationContext, propertyName string) bool {
	for k, _ := range oc.Variables {
		if k == propertyName {
			return true
		}
	}
	return false
}
