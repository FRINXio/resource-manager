// Contains helper functions used in tests (so that test functions are more readable)

import {
    getAllTags, findResourceTypeId, createNestedSingletonPool,
    findAllocationStrategyId, createAllocationPool, createNestedAllocationPool, getResourcesForPool
} from "./graphql-queries";
import _ from "underscore";

export async function getTag(tagName) {
    let allTags = await getAllTags();
    return _.find(allTags, d => d.Tag === tagName);
}

export async function createRandomIntRootPool() {
    let resourceTypeId = await findResourceTypeId('random_signed_int32');
    let strategyId = await findAllocationStrategyId('random_signed_int32');
    let poolName = getUniqueName('root-random-int');
    return await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { from: "int", to: "int"},
        {from: '1', to: 999},)
}

export async function createIpv4RootPool(address, prefix) {
    let resourceTypeId = await findResourceTypeId('ipv4');
    let strategyId = await findAllocationStrategyId('ipv4');
    let poolName = getUniqueName('root-ipv4');
    return await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { address: "string", prefix: "int"},
        {address: address, prefix: prefix},)
}

export async function createRdRootPool(){
    let resourceTypeId = await findResourceTypeId('route_distinguisher');
    let strategyId = await findAllocationStrategyId('route_distinguisher');
    let poolName = getUniqueName('root-rd');
    return await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { rd: "int"},
        {rd: 0},)
}


export async function createIpv6RootPool(){
    let resourceTypeId = await findResourceTypeId('ipv6');
    let strategyId = await findAllocationStrategyId('ipv6');
    let poolName = getUniqueName('root-ipv6-range');
    return await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { address: "string", prefix: "int"},
        {address: "dead::", prefix: 16},)
}

export async function createIpv6PrefixRootPool(){
    let resourceTypeId = await findResourceTypeId('ipv6_prefix');
    let strategyId = await findAllocationStrategyId('ipv6_prefix');
    let poolName = getUniqueName('root-ipv6-range');
    return await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { address: "string", prefix: "int"},
        {address: "dead::", prefix: 120},)
}


export async function createIpv6NestedPool(parentResourceId){
    let resourceTypeId = await findResourceTypeId('ipv6');
    let strategyId = await findAllocationStrategyId('ipv6');
    return await createNestedAllocationPool(
        getUniqueName('nested-ipv6'),
        resourceTypeId,
        strategyId,
        parentResourceId);
}

export async function createIpv4PrefixRootPool(){
    let resourceTypeId = await findResourceTypeId('ipv4_prefix');
    let strategyId = await findAllocationStrategyId('ipv4_prefix');
    return await createAllocationPool(
        getUniqueName('ipv4-root'),
        resourceTypeId,
        strategyId,
        {prefix: "int", address: "string"},
        {prefix: 8, address: "10.0.0.0"},);
}

export async function createIpv4NestedPool(parentResourceId) {
    let resourceTypeId = await findResourceTypeId('ipv4');
    let strategyId = await findAllocationStrategyId('ipv4');
    return await createNestedAllocationPool(
        getUniqueName('ipv4-nested'),
        resourceTypeId,
        strategyId,
        parentResourceId);
}

export async function createIpv4PrefixNestedPool(parentResourceId){
    let resourceTypeId = await findResourceTypeId('ipv4_prefix');
    let strategyId = await findAllocationStrategyId('ipv4_prefix');
    return await createNestedAllocationPool(
        getUniqueName('ipv4prefix-nested'),
        resourceTypeId,
        strategyId,
        parentResourceId);
}

export async function get2ChildrenIds(poolId) {
    let resourceForPool = await getResourcesForPool(poolId);
    return [resourceForPool[0].NestedPool.id, resourceForPool[1].NestedPool.id]
}

export async function createSingletonIpv4PrefixNestedPool(parentResourceId){
    let resourceTypeId = await findResourceTypeId('ipv4_prefix');
    return await createNestedSingletonPool(
        getUniqueName('singleton-ipv4prefix-nested'),
        resourceTypeId,
        [{address: "10.10.0.0", prefix: 11},],
        parentResourceId);
}

export async function createVlanRangeRootPool(){
    let resourceTypeId = await findResourceTypeId('vlan_range');
    let strategyId = await findAllocationStrategyId('vlan_range');
    return await createAllocationPool(
        getUniqueName('root-vlan-range'),
        resourceTypeId,
        strategyId,
        { from: "int", to: "int"},
        {from: 0, to: 4095},)
}

export async function createVlanRootPool(){
    let resourceTypeId = await findResourceTypeId('vlan');
    let strategyId = await findAllocationStrategyId('vlan');
    return await createAllocationPool(
        getUniqueName('root-vlan'),
        resourceTypeId,
        strategyId,
        { from: "int", to: "int"},
        {from: 0, to: 4095},)
}

export async function createVlanNestedPool(parentResourceId){
    let resourceTypeId = await findResourceTypeId('vlan');
    let strategyId = await findAllocationStrategyId('vlan');
    return await createNestedAllocationPool(
        getUniqueName('nested'),
        resourceTypeId,
        strategyId,
        parentResourceId)
}

export function getUniqueName(prefix){
    return prefix + Math.random().toString(36).substring(7);
}
