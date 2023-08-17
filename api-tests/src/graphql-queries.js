// Contains graphQL queries and mutations

import {client} from "./client.js";
import gql from 'graphql-tag';
import _ from 'underscore';


export async function getAllTags(){
    return client.query({
        query: gql`
            query { QueryTags{
                id
                Tag
                Pools {
                    id
                    Name
                }
            }
            }
        `,
    })
    .then(result => result.data.QueryTags)
    .catch(error => console.log(error));
}

export async function findResourceTypeId(name){
    return client.query({
        query: gql`
            query { QueryResourceTypes {
                id
                Name
            }
        }
        `,
    })
    .then(result => {
        let rt = _.find(result.data.QueryResourceTypes, d => d.Name === name);
        if (rt) {
            return rt.id;
        }

        return rt;
    })
    .catch(error => console.log(error));
}

export async function findAllocationStrategyId(name){
    return client.query({
        query: gql`
            query { QueryAllocationStrategies {
                id
                Name
            }
            }
        `,
    })
    .then(result => {
        const strategy = _.find(result.data.QueryAllocationStrategies, d => d.Name === name);
        if (!strategy) {
            return null;
        }
        return strategy.id;
    })
    .catch(error => console.log(error));
}

export async function getCapacityForPool(poolId){
    return client.query({
        query: gql`
            query getCapacity($poolId: ID!) {
                QueryPoolCapacity(poolId: $poolId) {
                   utilizedCapacity
                   freeCapacity
                }
            }
        `,
        variables: {
            poolId: poolId,
        }
    })
    .then(result => result.data.QueryPoolCapacity)
    .catch(error => console.log(error));
}

export async function getResourcePool(poolId, before, after, first, last) {
    if (!before && !after && !first && !last) {
        first = 10;
    }
    return client.query({
        query: gql`
            query getResourcePool($poolId: ID!, $before: String, $after: String, $first: Int, $last: Int) {
                QueryResourcePool(poolId: $poolId) {
                   AllocationStrategy {
                       id
                       Name
                   }
                   Name
                   PoolType
                   allocatedResources(first: $first, last: $last, before: $before, after: $after) {
                       edges {
                           cursor {
                               ID
                           }
                           node {
                               Properties
                               id
                           }
                       }
                       pageInfo {
                           endCursor {
                               ID
                           }
                           startCursor {
                               ID
                           }
                           hasNextPage
                           hasPreviousPage
                       }
                   }
                   ParentResource {
                       id
                       ParentPool {
                           id
                       }
                   }
                    Resources{
                        id
                        Properties
                        NestedPool {
                            id
                            Name
                            PoolType
                        }
                    }
                }
            }
        `,
        variables: {
            poolId: poolId,
            before: before,
            after: after,
            first: first,
            last: last,
        }
    })
    .then(result => result.data.QueryResourcePool)
    .catch(error => console.log(error));
}

export async function getResourcesByDatetimeRange(fromDatetime, toDatetime) {
    return client.query({
        query: gql`
            query QueryRecentlyActiveResources($fromDatetime: String!, $toDatetime: String) {
                QueryRecentlyActiveResources(fromDatetime:$fromDatetime, toDatetime:$toDatetime) {
                    edges {
                        node {
                            id
                            Properties
                            NestedPool {
                                id
                                Name
                                PoolType
                            }
                        }
                    }
                }
            }
        `,
        variables: {
            fromDatetime: fromDatetime,
            toDatetime: toDatetime,
        }
    })
    .then(result => result.data.QueryRecentlyActiveResources)
    .catch(error => console.log(error));
}

export async function getResourcesForPool(poolId){
    return client.query({
        query: gql`
            query getResources($poolId: ID!) {
                 QueryResources(poolId: $poolId) {
                     edges {
                         node {
                             id
                             Properties
                             NestedPool {
                                 id
                                 Name
                                 PoolType
                             }
                             AlternativeId
                         }
                     }
                 }
            }
        `,
        variables: {
            poolId: poolId,
        }
    })
    .then(result => result.data.QueryResources)
    .catch(error => console.log(error));
}

export async function getPaginatedResourcesForPool(poolId, first, last, before, after){
    return client.query({
        query: gql`
            query getResources($poolId: ID!, $first: Int, $last: Int, $before: String, $after: String) {
                QueryResources(poolId: $poolId,  first: $first, last: $last, before: $before, after: $after) {
                    edges {
                        node {
                            id
                            Properties
                            NestedPool {
                                id
                                Name
                                PoolType
                            }
                        }
                        cursor{
                            ID
                        }
                    }
                }
            }
        `,
        variables: {
            poolId: poolId,
            first: first,
            last: last,
            before: before,
            after: after
        }
    })
    .then(result => result.data.QueryResources)
    .catch(error => console.log(error));
}

export async function searchPoolsByTags(searchExpression){
    return client.query({
        query: gql`
            query getPoolsByTags($searchExpression: TagOr) {
                SearchPoolsByTags(tags: $searchExpression) {
                    id
                }
            }
        `,
        variables: {
            searchExpression: searchExpression,
        }
    })
    .then(result => result.data.SearchPoolsByTags)
    .catch(error => console.log(error));
}

export async function getPoolHierarchyPath(poolId){
    return client.query({
        query: gql`
            query QueryPoolHierarchyPath($poolId: ID!) {
                QueryResourcePoolHierarchyPath(poolId: $poolId) {
                    id
                    Name
                }
            }
        `,
        variables: {
            poolId: poolId,
        }
    })
    .then(result => result.data.QueryResourcePoolHierarchyPath)
    .catch(error => console.log(error));
}


export async function createNestedAllocationPool(poolName, resourceTypeId, strategyId, parentResourceId, tags = null){
    return client.mutate({
        mutation: gql`
            mutation createNestedAllocPool($poolName: String!, $resourceTypeId: ID!, $strategyId: ID!, $parentResourceId: ID!, $tags: [String!]) {
                CreateNestedAllocatingPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    description: "pool for testing"
                    allocationStrategyId: $strategyId
                    poolDealocationSafetyPeriod: 0
                    parentResourceId: $parentResourceId
                    tags: $tags
                }){
                    pool
                    {
                        id
                    }
                }
            }
        `,
        variables: {
            poolName: poolName,
            resourceTypeId: resourceTypeId,
            strategyId: strategyId,
            parentResourceId: parentResourceId,
            tags: tags,

        }
    })
    .then(result => result.data.CreateNestedAllocatingPool.pool.id)
    .catch(error => console.log(error));
}

export async function testStrategy(allocationStrategyId, poolProperties, poolName, poolId, currentResources, userInput, suppressErrors = false) {
    return client.mutate({
        mutation: gql`
            mutation testStrategy(
                $allocationStrategyId: ID!,
                $poolProperties: Map!,
                $poolName: String!,
                $poolId: ID!,
                $currentResources: [ResourceInput!]!,
                $userInput: Map!) {
                TestAllocationStrategy(
                    allocationStrategyId: $allocationStrategyId
                    resourcePool: {
                        ResourcePoolID: $poolId
                        poolProperties: $poolProperties
                        ResourcePoolName: $poolName
                    }
                    currentResources: $currentResources
                    userInput: $userInput
                )
            }
        `,
        variables: {
            allocationStrategyId: allocationStrategyId,
            poolProperties: poolProperties,
            poolName: poolName,
            currentResources: currentResources,
            userInput: userInput,
            poolId: poolId
        }
    })
    .then(result => result.data.TestAllocationStrategy)
    .catch(error => suppressErrors?null:console.log(error));
}

export async function createTag(tagText) {
    return client.mutate({
        mutation: gql`
            mutation createTag($tagText: String!) {
                CreateTag( input:  {
                    tagText: $tagText
                }){
                    tag {
                        id
                    }
                }
            }
        `,
        variables: {
            tagText: tagText,
        }
    })
    .then(result => result.data.CreateTag.tag.id)
    .catch(error => console.log(error));
}

export async function tagPool(tagId, poolId) {
    return client.mutate({
        mutation: gql`
            mutation tagPool($tagId: ID!, $poolId: ID!) {
                TagPool( input:  {
                    tagId: $tagId
                    poolId: $poolId
                }){
                    tag {
                        id
                    }
                }
            }
        `,
        variables: {
            tagId: tagId,
            poolId: poolId,
        }
    })
    .then(result => result.data.TagPool.tag.id)
    .catch(error => console.log(error));
}

export async function untagPool(tagId, poolId) {
    return client.mutate({
        mutation: gql`
            mutation untagPool($tagId: ID!, $poolId: ID!) {
                UntagPool( input:  {
                    tagId: $tagId
                    poolId: $poolId
                }){
                    tag {
                        id
                    }
                }
            }
        `,
        variables: {
            tagId: tagId,
            poolId: poolId,
        }
    })
    .then(result => result.data.UntagPool.tag.id)
    .catch(error => console.log(error));
}

export async function deleteTag(tagId) {
    return client.mutate({
        mutation: gql`
            mutation deleteTag($tagId: ID!) {
                DeleteTag( input:  {
                    tagId: $tagId
                }){
                    tagId
                }
            }
        `,
        variables: {
            tagId: tagId,
        }
    })
    .then(result => result.data.DeleteTag.tagId)
    .catch(error => console.log(error));
}

export async function deleteAllocationStrategy(strategyId) {
    return client.mutate({
        mutation: gql`
            mutation deleteStrategy($allocationStrategyId: ID!) {
                DeleteAllocationStrategy( input:  {
                    allocationStrategyId: $allocationStrategyId
                }){
                    strategy {
                        id
                    }
                }
            }
        `,
        variables: {
            allocationStrategyId: strategyId,
        }
    })
    .then(result => result.data.DeleteAllocationStrategy.strategy.id)
    .catch(error => console.log(error));
}

export async function createAllocationStrategy(strategyName, scriptBody, strategyType, expectedPropertyTypes) {
    return client.mutate({
        mutation: gql`
            mutation createStrategy($strategyName: String!, $scriptBody: String!, $strategyType: AllocationStrategyLang!,
                $expectedPropertyTypes: Map) {
                CreateAllocationStrategy( input:  {
                    name: $strategyName
                    description: "testing strategy"
                    script: $scriptBody
                    lang: $strategyType
                    expectedPoolPropertyTypes: $expectedPropertyTypes
                }){
                    strategy {
                        id
                    }
                }
            }
        `,
        variables: {
            strategyName: strategyName,
            scriptBody: scriptBody,
            strategyType: strategyType,
            expectedPropertyTypes: expectedPropertyTypes
        }
    })
    .then(result => result.data.CreateAllocationStrategy.strategy.id)
    .catch(error => console.log(error));
}

export async function deleteResourcePool(poolId) {
    return client.mutate({
        mutation: gql`
            mutation deletePool($resourcePoolId: ID!) {
                DeleteResourcePool( input:  {
                    resourcePoolId: $resourcePoolId
                }){
                   resourcePoolId
                }
            }
        `,
        variables: {
            resourcePoolId: poolId,
        }
    })
    .then(result => result.data.DeleteResourcePool)
    .catch(error => console.log(error));
}

export async function deleteResourceType(resourceTypeId) {
    return client.mutate({
        mutation: gql`
            mutation deleteResType($resourceTypeId: ID!) {
                DeleteResourceType( input:  {
                    resourceTypeId: $resourceTypeId
                }){
                    resourceTypeId
                }
            }
        `,
        variables: {
            resourceTypeId: resourceTypeId,
        }
    })
    .then(result => result.data.DeleteResourceType)
    .catch(error => console.log(error));
}

export async function createResourceType(resourceName, resourceProperties) {
    return client.mutate({
        mutation: gql`
            mutation createResType($resourceName: String!, $resourceProperties: Map!) {
                CreateResourceType( input:  {
                    resourceName: $resourceName
                    resourceProperties: $resourceProperties
                }){
                    resourceType {
                        id
                    }
                }
            }
        `,
        variables: {
            resourceName: resourceName,
            resourceProperties: resourceProperties,
        }
    })
    .then(result => result.data.CreateResourceType.resourceType.id)
    .catch(error => console.log(error));
}

export async function createNestedSingletonPool(poolName, resourceTypeId, poolValues, parentResourceId){
    return client.mutate({
        mutation: gql`
            mutation createSingPool($poolName: String!, $resourceTypeId: ID!, $poolValues: [Map!]!, $parentResourceId: ID!) {
                CreateNestedSingletonPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    description: "singleton nested pool for testing"
                    poolValues: $poolValues
                    parentResourceId: $parentResourceId
                }){
                    pool
                    {
                        id
                    }
                }
            }
        `,
        variables: {
            poolName: poolName,
            resourceTypeId: resourceTypeId,
            poolValues: poolValues,
            parentResourceId: parentResourceId,
        }
    })
    .then(result => result.data.CreateNestedSingletonPool.pool.id)
    .catch(error => console.log(error));
}

export async function createSingletonPool(poolName, resourceTypeId, poolValues){
    return client.mutate({
        mutation: gql`
            mutation createSingPool($poolName: String!, $resourceTypeId: ID!, $poolValues: [Map!]!) {
                CreateSingletonPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    description: "singleton pool for testing"
                    poolValues: $poolValues
                }){
                    pool
                    {
                        id
                    }
                }
            }
        `,
        variables: {
            poolName: poolName,
            resourceTypeId: resourceTypeId,
            poolValues: poolValues,
        }
    })
    .then(result => result.data.CreateSingletonPool.pool.id)
    .catch(error => console.log(error));
}

export async function createSetPool(poolName, resourceTypeId, poolValues){
    return client.mutate({
        mutation: gql`
            mutation createSetPool($poolName: String!, $resourceTypeId: ID!, $poolValues: [Map!]!) {
                CreateSetPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    description: "set pool for testing"
                    poolDealocationSafetyPeriod: 0
                    poolValues: $poolValues
                }){
                    pool
                    {
                        id
                    }
                }
            }
        `,
        variables: {
            poolName: poolName,
            resourceTypeId: resourceTypeId,
            poolValues: poolValues,
        }
    })
    .then(result => result.data.CreateSetPool.pool.id)
    .catch(error => console.log(error));
}

export async function createAllocationPool(poolName, resourceTypeId, strategyId, poolPropertyTypes, poolProperties, tags = null, suppressErrors = false){
    return client.mutate({
        mutation: gql`
            mutation createAllocPool($poolName: String!, $resourceTypeId: ID!, $strategyId: ID!, $poolProperties: Map!, $poolPropertyTypes: Map!, $tags: [String!]) {
                CreateAllocatingPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    allocationStrategyId: $strategyId
                    poolDealocationSafetyPeriod: 0
                    poolPropertyTypes: $poolPropertyTypes
                    poolProperties: $poolProperties
                    tags: $tags
                    }){
                pool
                    {
                        id
                        PoolProperties
                    }
            }
            }
        `,
        variables: {
            poolName: poolName,
            resourceTypeId: resourceTypeId,
            strategyId: strategyId,
            poolProperties: poolProperties,
            poolPropertyTypes: poolPropertyTypes,
            tags: tags,
        }
    })
    .then(result => result.data.CreateAllocatingPool.pool)
    .catch(error => suppressErrors?null:console.log(error));
}

export async function claimResource(poolId, params, description = null, suppressErrors = false){
    return client.mutate({
        mutation: gql`
            mutation claimRes($poolId: ID!, $description: String, $userInput: Map!) {
               ClaimResource( poolId: $poolId, description: $description, userInput: $userInput){
                   id
                   Properties
                }
            }
        `,
        variables: {
            poolId: poolId,
            userInput: params,
            description: description,
        }
    })
    .then(result => result.data.ClaimResource)
    .catch(error => suppressErrors?null:console.log(error));
}

export async function claimResourceWithAltId(poolId, params, altId, description = null, suppressErrors = false){
    return client.mutate({
        mutation: gql`
            mutation claimResAltId($poolId: ID!, $description: String, $userInput: Map!, $alternativeId: Map!) {
                ClaimResourceWithAltId( poolId: $poolId, description: $description, userInput: $userInput, alternativeId: $alternativeId){
                    id
                    Properties
                }
            }
        `,
        variables: {
            poolId: poolId,
            userInput: params,
            alternativeId: altId,
            description: description,
        }
    })
    .then(result => {
        if (!result || !result.data) {
            return null;
        }

        return result.data.ClaimResourceWithAltId;
    })
    .catch(error => {
        if (!suppressErrors) {
            console.log(error);
        }
    });
}

export async function queryResource(poolId, params){
    return client.query({
        query: gql`
            query QueryResource($poolId: ID!, $input: Map!) {
                QueryResource(input: $input, poolId: $poolId) {
                   id
                   Properties
                   Description
                }
            }
        `,
        variables: {
            poolId: poolId,
            input: params,
        }
    })
    .then(result => result.data.QueryResource)
    .catch(error => console.log(error));
}

export async function queryResourcesByAltIdAndPoolId(poolId, params){
    return client.query({
        query: gql`
            query QueryResourcesByAltId($poolId: ID, $input: Map!) {
                QueryResourcesByAltId(input: $input, poolId: $poolId) {
                    edges {
                        node {
                            id
                            Properties
                            NestedPool {
                                id
                                Name
                                PoolType
                            }
                            AlternativeId
                        }
                    }
                }
            }
        `,
        variables: {
            poolId: poolId,
            input: params,
        }
    })
    .then(result => {
        if (result.data) {
            return result.data.QueryResourcesByAltId;
        }

        return null;
    })
    .catch(error => {
        console.log(error);
        return null;
    });
}

export async function queryResourcesByAltId(params){
    return client.query({
        query: gql`
            query QueryResourcesByAltId($poolId: ID, $input: Map!) {
                QueryResourcesByAltId(input: $input, poolId: $poolId) {
                    edges {
                        node {
                            id
                            Properties
                            NestedPool {
                                id
                                Name
                                PoolType
                            }
                            AlternativeId
                        }
                    }
                }
            }
        `,
        variables: {
            input: params,
        }
    })
    .then(result => {
        if (result.data) {
            return result.data.QueryResourcesByAltId;
        }

        return null;
    })
    .catch(error => {
        console.log(error);
        return null;
    });
}

export async function freeResource(poolId, propInput){
    return client.mutate({
        mutation: gql`
            mutation freeRes($poolId: ID!, $propInput: Map!) {
               FreeResource( input: $propInput, poolId: $poolId )
            }
        `,
        variables: {
            poolId: poolId,
            propInput: propInput,
        }
    })
    .then(result => result.data.FreeResource)
    .catch(error => console.log(error));
}

export async function getLeafPools(resourceTypeId, tags) {
    return client.query({
        query: gql`
            query QueryLeafResourcePools($resourceTypeId: ID, $tags: TagOr) {
                QueryLeafResourcePools(resourceTypeId: $resourceTypeId, tags: $tags) {
                    id
                    Name
                    AllocationStrategy {
                        id
                        Name
                        Description
                    }
                    ResourceType {
                        id
                        Name
                    }
                    Resources{
                        id
                        Properties
                        NestedPool {
                            id
                            Name
                            PoolType
                        }
                    }
                }
            }
        `,
        variables: {
            resourceTypeId: resourceTypeId,
            tags: tags
        }
    })
    .then(result => result.data.QueryLeafResourcePools)
    .catch(error => console.log(error));
}

export async function getAllPoolsByTypeOrTag(resourceTypeId, tags) {
    return client.query({
        query: gql`
            query QueryResourcePools($resourceTypeId: ID, $tags: TagOr) {
                QueryResourcePools(resourceTypeId: $resourceTypeId, tags: $tags) {
                    id
                    Name
                    AllocationStrategy {
                        id
                        Name
                        Description
                    }
                    ResourceType {
                        id
                        Name
                    }
                    Capacity {
                        utilizedCapacity
                        freeCapacity
                    }
                    Resources{
                        id
                        Properties
                        NestedPool {
                            id
                            Name
                            PoolType
                        }
                    }
                }
            }
        `,
        variables: {
            resourceTypeId: resourceTypeId,
            tags: tags
        }
    })
    .then(result => result.data.QueryResourcePools)
    .catch(error => console.log(error));
}

export async function updateResourceAltId(poolId, input, alternativeId){
    return client.mutate({
        mutation: gql`
            mutation UpdateResourceAltId($input: Map!, $poolId: ID!, $alternativeId: Map!) {
                UpdateResourceAltId(input: $input, poolId: $poolId, alternativeId:$alternativeId) {
                    AlternativeId
                }
            }
        `,
        variables: {
            poolId: poolId,
            input: input,
            alternativeId: alternativeId
        }
    })
    .then(result => result.data.UpdateResourceAltId)
    .catch(error => console.log(error));
}

export async function getEmptyPools(resourceTypeId){
    return client.query({
        query: gql`
            query QueryEmptyResourcePools($resourceTypeId: ID) {
                QueryEmptyResourcePools(resourceTypeId: $resourceTypeId) {
                    id
                    Name
                    AllocationStrategy {
                        id
                        Name
                        Description
                    }
                    ResourceType {
                        id
                        Name
                    }
                    Resources{
                        id
                        Properties
                        NestedPool {
                            id
                            Name
                            PoolType
                        }
                    }
                }
            }
        `,
        variables: {
            resourceTypeId: resourceTypeId
        }
    })
    .then(result => result.data.QueryEmptyResourcePools)
    .catch(error => console.log(error));
}

export async function getRequiredPoolProperties(allocationStrategyName) {
    return client.query({
        query: gql`
            query QueryRequiredPoolProperties($allocationStrategyName: String!) {
                QueryRequiredPoolProperties(allocationStrategyName: $allocationStrategyName) {
                    Name
                    Type
                    FloatVal
                    IntVal
                    StringVal
                }
            }
        `,
        variables: {
            allocationStrategyName: allocationStrategyName
        }
    })
    .then(result => result.data.QueryRequiredPoolProperties)
    .catch(error => console.log(error));
}