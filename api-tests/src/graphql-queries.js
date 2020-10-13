// Contains graphQL queries and mutations

import {client} from "./client";
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

export async function getResourcesForPool(poolId){
    return client.query({
        query: gql`
            query getResources($poolId: ID!) {
                 QueryResources(poolId: $poolId) {
                     id
                     Properties
                     NestedPool {
                         id
                         Name
                         PoolType
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


export async function createNestedAllocationPool(poolName, resourceTypeId, strategyId, parentResourceId){
    return client.mutate({
        mutation: gql`
            mutation createNestedAllocPool($poolName: String!, $resourceTypeId: ID!, $strategyId: ID!, $parentResourceId: ID!) {
                CreateNestedAllocatingPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    description: "pool for testing"
                    allocationStrategyId: $strategyId
                    poolDealocationSafetyPeriod: 0
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
            strategyId: strategyId,
            parentResourceId: parentResourceId,

        }
    })
    .then(result => result.data.CreateNestedAllocatingPool.pool.id)
    .catch(error => console.log(error));
}

export async function testStrategy(allocationStrategyId, poolProperties, poolName, currentResources, userInput) {
    return client.mutate({
        mutation: gql`
            mutation testStrategy(
                $allocationStrategyId: ID!,
                $poolProperties: Map!,
                $poolName: String!, 
                $currentResources: [ResourceInput!]!,
                $userInput: Map!) {
                TestAllocationStrategy(
                    allocationStrategyId: $allocationStrategyId
                    resourcePool: {
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
        }
    })
    .then(result => result.data.TestAllocationStrategy)
    .catch(error => console.log(error));
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

export async function createAllocationStrategy(strategyName, scriptBody, strategyType) {
    return client.mutate({
        mutation: gql`
            mutation createStrategy($strategyName: String!, $scriptBody: String!, $strategyType: AllocationStrategyLang!) {
                CreateAllocationStrategy( input:  {
                    name: $strategyName
                    description: "testing strategy"
                    script: $scriptBody
                    lang: $strategyType
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

export async function createAllocationPool(poolName, resourceTypeId, strategyId, poolPropertyTypes, poolProperties){
    return client.mutate({
        mutation: gql`
            mutation createAllocPool($poolName: String!, $resourceTypeId: ID!, $strategyId: ID!, $poolProperties: Map!, $poolPropertyTypes: Map!) { 
                CreateAllocatingPool( input:  {
                    resourceTypeId: $resourceTypeId
                    poolName: $poolName
                    allocationStrategyId: $strategyId
                    poolDealocationSafetyPeriod: 0
                    poolPropertyTypes: $poolPropertyTypes
                    poolProperties: $poolProperties
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
            poolProperties: poolProperties,
            poolPropertyTypes: poolPropertyTypes,

        }
    })
    .then(result => result.data.CreateAllocatingPool.pool.id)
    .catch(error => console.log(error));
}

export async function claimResource(poolId, params){
    return client.mutate({
        mutation: gql`
            mutation claimRes($poolId: ID!, $userInput: Map!) {
               ClaimResource( poolId: $poolId, userInput: $userInput){
                   id
                   Properties
                }
            }
        `,
        variables: {
            poolId: poolId,
            userInput: params,
        }
    })
    .then(result => result.data.ClaimResource)
    .catch(error => console.log(error));
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