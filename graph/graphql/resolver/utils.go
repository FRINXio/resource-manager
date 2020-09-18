package resolver

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"sort"
)

func createNestedPool(ctx context.Context,
	parentResourceID int,
	client *ent.Client,
	nestedPoolFactory func() (*ent.ResourcePool, error)) (*ent.ResourcePool, error) {

	parentResource, errFindRes := client.Resource.Get(ctx, parentResourceID)
	if errFindRes != nil {
		return nil, gqlerror.Errorf("Cannot create nested pool, Unable to find parent resource: %v", errFindRes)
	}

	pool, errCreatePool := nestedPoolFactory()
	if errCreatePool != nil {
		return nil, errCreatePool
	}

	//create pool properties from parent resource type
	if pool.PoolType == resourcepool.PoolTypeAllocating {

		properties, err := parentResource.QueryProperties().All(ctx)

		if err != nil {
			return nil, gqlerror.Errorf("Cannot retrieve properties, error: %v", err)
		}

		poolProperties, err := client.PoolProperties.Create().AddProperties(properties...).Save(ctx)

		if err != nil {
				return nil, gqlerror.Errorf("Cannot create pool properties, error: %v", err)
		}

		_, err = pool.Update().SetPoolProperties(poolProperties).Save(ctx)

		if err != nil {
			return nil, gqlerror.Errorf("Unable to set pool properties on the given pool, error: %v", err)
		}
	}

	pool, errSetParentResource := pool.Update().SetParentResource(parentResource).Save(ctx)
	if errSetParentResource != nil {
		return pool, nil
	}

	return pool, nil
}

type ById []*ent.Resource

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func toResourceEdges (res []*ent.Resource) []*model.ResourceEdge {
	var edges []*model.ResourceEdge

	for _, v := range res {
		edges = append(edges, &model.ResourceEdge{
			Node:   v,
			Cursor: v.ID,
		},
		)
	}

	return edges
}

func populateResourceConnection(res []*ent.Resource, hasNext bool, hasPrevious bool) *model.ResourceConnection {
	pi := &model.PageInfo{
		HasPreviousEdge: hasPrevious,
		HasNextEdge:     hasNext,
		StartCursor:     res[0].ID,
		EndCursor:       res[len(res)-1].ID,
	}

	con := model.ResourceConnection{
		PageInfo: pi,
		Edges:    toResourceEdges(res),
	}

	return &con
}

func sortResources(res []*ent.Resource) []*ent.Resource {
	sort.Sort(ById(res))
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getAfterResources(res []*ent.Resource, after int, first int) ([]*ent.Resource, bool, bool) {
	sorted := sortResources(res)

	for i, v := range sorted {
		if v.ID == after {
			upper := min(len(res), i+first)
			hasNext := len(res) < i+first
			return res[i : upper], hasNext, i > 0
		}
	}

	return nil, false, false //TODO throw an error?
}

func getBeforeResources(res []*ent.Resource, before int, last int) ([]*ent.Resource, bool, bool) {
	sorted := sortResources(res)

	for i, v := range sorted {
		if v.ID == before {
			lower := max(0, i - last)
			return res[lower : i], len(res) > i, i - last > 0
		}
	}

	return nil, false, false //TODO throw an error?
}