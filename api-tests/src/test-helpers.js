// Contains helper functions used in tests (so that test functions are more readable)

import {
    getAllTags, findResourceTypeId, createNestedSingletonPool,
    findAllocationStrategyId, createAllocationPool,
    createNestedAllocationPool, getResourcesForPool,
    claimResource, getResourcePool
} from "./graphql-queries.js";
import _ from "underscore";

export async function getTag(tagName) {
    let allTags = await getAllTags();
    return _.find(allTags, d => d.Tag === tagName);
}

export async function createRandomIntRootPool() {
    let resourceTypeId = await findResourceTypeId('random_signed_int32');
    let strategyId = await findAllocationStrategyId('random_signed_int32');
    let poolName = getUniqueName('root-random-int');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { from: "int", to: "int"},
        {from: '1', to: 999},)

    return pool.id;
}

export async function createIpv4RootPool(address, prefix) {
    let resourceTypeId = await findResourceTypeId('ipv4');
    let strategyId = await findAllocationStrategyId('ipv4');
    let poolName = getUniqueName('root-ipv4');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { address: "string", prefix: "int"},
        {address: address, prefix: prefix},)

    return pool.id
}

export async function createRdRootPool(){
    let resourceTypeId = await findResourceTypeId('route_distinguisher');
    let strategyId = await findAllocationStrategyId('route_distinguisher');
    let poolName = getUniqueName('root-rd');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { rd: "int"},
        {rd: 0},)

    return pool.id;
}


export async function createIpv6RootPool(){
    let resourceTypeId = await findResourceTypeId('ipv6');
    let strategyId = await findAllocationStrategyId('ipv6');
    let poolName = getUniqueName('root-ipv6-range');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { address: "string", prefix: "int"},
        {address: "dead::", prefix: 16},)

    return pool.id;
}

export async function createIpv6PrefixRootPool(){
    let resourceTypeId = await findResourceTypeId('ipv6_prefix');
    let strategyId = await findAllocationStrategyId('ipv6_prefix');
    let poolName = getUniqueName('root-ipv6-range');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        { address: "string", prefix: "int"},
        {address: "dead::", prefix: 120},)

    return pool.id;
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

export async function createIpv4PrefixRootPool() {
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
    const pool = await createAllocationPool(
        getUniqueName('root-vlan-range'),
        resourceTypeId,
        strategyId,
        { from: "int", to: "int"},
        {from: 0, to: 4095},)

    return pool.id;
}

export async function createVlanRootPool(tags = null){
    let resourceTypeId = await findResourceTypeId('vlan');
    let strategyId = await findAllocationStrategyId('vlan');
    const pool = await createAllocationPool(
        getUniqueName('root-vlan'),
        resourceTypeId,
        strategyId,
        { from: "int", to: "int"},
        {from: 0, to: 4095},
        tags)

    return pool.id;
}

export async function createVlanNestedPool(parentResourceId, tags = null){
    let resourceTypeId = await findResourceTypeId('vlan');
    let strategyId = await findAllocationStrategyId('vlan');
    return await createNestedAllocationPool(
        getUniqueName('nested'),
        resourceTypeId,
        strategyId,
        parentResourceId,
        tags)
}

export function getUniqueName(prefix){
    return prefix + Math.random().toString(36).substring(5);
}

export async function prepareIpv4Pool() {
    let rootPoolId = (await createIpv4PrefixRootPool()).id;
    let resource12Id = (await claimResource(rootPoolId, { desiredSize: 4194304 })).id;
    return await createIpv4NestedPool(resource12Id);
}

export async function allocateFromIPv4PoolSerially(poolId, count, claimResourceParams) {
    const result = [];
    for (let i = 0; i < count; i++) {
        result.push(await claimResource(poolId, claimResourceParams));
    }
    const ips = result.filter(it => it != null).map(it => it.Properties.address);
    const uniqIPs = [...new Set(ips)];
    if (uniqIPs.length != count) {
        console.error({ result, ips, uniqIPs });
        throw new Error("Unexpected number of IPs:" + uniqIPs.length);
    }
    return uniqIPs;
}

export async function allocateFromIPv4PoolParallelly(poolId, count, retries, claimResourceParams) {
    const promises = [];
    for (let i = 0; i < count; i++) {
        promises.push(claimResource(poolId, claimResourceParams, null, true));
    }
    const result = await Promise.all(promises);
    const ips = result.filter(it => it != null).map(it => it.Properties.address);
    if (ips.length != count) {
        if (ips.length == 0) {
            throw new Error("Nothing alocated at " + JSON.stringify({count, retries}))
        } else if (retries-- > 0) {
            // call recursively
            const moreIPs = await allocateFromIPv4PoolParallelly(poolId, count - ips.length, retries, claimResourceParams);
            ips.push(...moreIPs);
        } else {
            console.error({ result, ips });
            throw new Error("Unexpected number of IPs:" + ips.length);
        }
    }
    return ips;
}

export async function queryIPs(poolId, count) {
    const queried = await getResourcePool(poolId, null, null, count);
    const allocated = queried.allocatedResources;
    if (!allocated || !allocated.edges) {
        console.error("Unexpected query result", {queried});
        throw new Error("Unexpected query result");
    }
    const foundIPs = allocated.edges.map(it => it.node.Properties.address);
    const uniqIPs = [...new Set(foundIPs)];
    if (foundIPs.length != uniqIPs.length) {
        console.error("Query returned duplicates", {foundIPs, uniqIPs});
        throw new Error("Query returned duplicates");
    }
    return uniqIPs;
}
