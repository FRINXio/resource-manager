import {
    deleteAllocationStrategy,
    testStrategy,
    createAllocationStrategy,
    findAllocationStrategyId,
    createResourceType, findResourceTypeId, deleteResourceType, getRequiredPoolProperties
} from '../graphql-queries.js';
import {cleanup, getUniqueName} from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;

test('create and call JS strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await createAllocationStrategy(
        poolName,
        'function invoke() {return {vlan: userInput.desiredVlan};}',
        'js', null);
    let strategyOutput = await testStrategy(strategyId, { ResourcePoolName: 'testpool'},
        poolName, poolId, [], {desiredVlan: 85} );

    t.equal(strategyOutput.stdout.vlan, 85);

    await cleanup()
    t.end();
});

test('create and call Py strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await createAllocationStrategy(
        poolName,
        'log(json.dumps({ \'respool\': resourcePool[\'ResourcePoolName\'], \'currentRes\': currentResources }))\nreturn {\'vlan\': userInput[\'desiredVlan\']}',
        'py', null);
    let strategyOutput = await testStrategy(strategyId, { ResourcePoolName: poolName},
        poolName, poolId,[], {desiredVlan: 11} );

    t.equal(strategyOutput.stdout.vlan, 11);

    await cleanup()
    t.end();
});

test('delete strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let strategyId = await createAllocationStrategy(
        poolName,
        'function invoke() {return {vlan: userInput.desiredVlan};}', 'js', null);
    let foundStrategyId = await findAllocationStrategyId(poolName);
    t.equal(foundStrategyId, strategyId);
    await deleteAllocationStrategy(strategyId);
    foundStrategyId = await findAllocationStrategyId(poolName);
    t.notOk(foundStrategyId);
    t.end();
});

test('simple ipv4 prefix strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let ipv4PrefixStrategyId = await findAllocationStrategyId('ipv4_prefix');
    let x = await testStrategy(ipv4PrefixStrategyId, {prefix: 8, address: '10.0.0.0', subnet: false},
        poolName, poolId,
        [], {desiredSize: 8388608});

    t.equal(x.stdout.address, '10.0.0.0');
    t.equal(x.stdout.prefix, 9);

    await cleanup()
    t.end();
});

test('ipv4 prefix strategy one resource already claimed', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let ipv4PrefixStrategyId = await findAllocationStrategyId('ipv4_prefix');
    let allocated = await testStrategy(ipv4PrefixStrategyId, {prefix: 8, address: '10.0.0.0', subnet: false},
        poolName, poolId,
        [{Properties: { prefix: 9, address: '10.0.0.0'},
            Status: 'claimed',
            UpdatedAt: '2020-08-18 11:38:48.0 +0200 CEST'
        }], {desiredSize: 8388608});

    t.equal(allocated.stdout.address, '10.128.0.0');
    t.equal(allocated.stdout.prefix, 9);

    await cleanup()
    t.end();
});

test('ipv4 prefix strategy pool has no capacity left', async (t) => {
    const poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    const allocatedResources = [
        {Properties: { prefix: 9, address: '10.0.0.0'},
            Status: 'claimed',
            UpdatedAt: '2020-08-18 11:38:48.0 +0200 CEST'
        },
        {Properties: { prefix: 9, address: '10.128.0.0'},
            Status: 'claimed',
            UpdatedAt: '2020-08-18 11:38:48.0 +0200 CEST'
        }];
    let ipv4PrefixStrategyId = await findAllocationStrategyId('ipv4_prefix');
    let allocated = await testStrategy(ipv4PrefixStrategyId, {prefix: 8, address: '10.0.0.0', subnet: false},
        poolName, poolId, allocatedResources, {desiredSize: 8388608}, true);

    t.notOk(allocated);

    await cleanup()
    t.end();
});

test('ipv4 strategy just get an IP', async (t) => {
    const poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let ipv4StrategyId = await findAllocationStrategyId('ipv4');
    let allocated = await testStrategy(ipv4StrategyId, {prefix: 8, address: '10.0.0.0', subnet: false},
        poolName, poolId, [], {});

    t.equal(allocated.stdout.address, '10.0.0.0');

    await cleanup()
    t.end();
});


test('simple ipv6 prefix strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let ipv6PrefixStrategyId = await findAllocationStrategyId('ipv6_prefix');
    let allocated = await testStrategy(ipv6PrefixStrategyId, {prefix: 120, address: 'dead::', subnet: false},
        poolName, poolId, [], {desiredSize: 101});
    t.equal(allocated.stdout.address, 'dead::');
    t.equal(allocated.stdout.prefix, 121);

    await cleanup()
    t.end();
});

test('simple ipv6 strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let ipv6StrategyId = await findAllocationStrategyId('ipv6');
    let allocated = await testStrategy(ipv6StrategyId, {prefix: 120, address: 'dead::', subnet: true},
        poolName, poolId, [], {});
    t.equal(allocated.stdout.address, 'dead::1');

    await cleanup()
    t.end();
});

test('ipv4-rd strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await findAllocationStrategyId('route_distinguisher');
    let allocated = await testStrategy(strategyId, {},
        poolName, poolId, [], {ipv4: '1.2.3.4', assignedNumber: 2});

    t.equal(allocated.stdout.rd, '1.2.3.4:2');

    await cleanup()
    t.end();
});

test('check required resource types for different strategies', async (t) => {
    let requiredPoolProperties = await getRequiredPoolProperties('vlan');

    t.equal(requiredPoolProperties[0].Name, 'from')
    t.equal(requiredPoolProperties[0].Type, 'int')
    t.equal(requiredPoolProperties[1].Name, 'to')
    t.equal(requiredPoolProperties[1].Type, 'int')

    requiredPoolProperties = await getRequiredPoolProperties("ipv4");

    t.equal(requiredPoolProperties[0].Name, 'address')
    t.equal(requiredPoolProperties[0].Type, 'string')
    t.equal(requiredPoolProperties[1].Name, 'prefix')
    t.equal(requiredPoolProperties[1].Type, 'int')
    t.equal(requiredPoolProperties[2].Name, 'subnet')
    t.equal(requiredPoolProperties[2].Type, 'bool')
    t.end();
});

test('ipv4-rd strategy duplicate already claimed', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    const claimed = [{Properties: {rd: '1.2.3.4:2'},
        Status: 'claimed',
        UpdatedAt: '2020-08-30 11:38:48.0 +0200 CEST'
    }];

    let strategyId = await findAllocationStrategyId('route_distinguisher');
    let allocated = await testStrategy(strategyId, {},
        poolName, poolId, claimed, {ipv4: '1.2.3.4', assignedNumber: 2}, true);

    t.notOk(allocated);

    await cleanup()
    t.end();
});


test('as-rd strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await findAllocationStrategyId('route_distinguisher');
    let allocated = await testStrategy(strategyId, {},
        poolName, poolId, [], {asNumber: 22, assignedNumber: 288});

    t.equal(allocated.stdout.rd, '22:288');

    await cleanup()
    t.end();
});

test('vlan range strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await findAllocationStrategyId('vlan_range');
    let allocated = await testStrategy(strategyId,
        {from: 0, to: 4095}, poolName, poolId, [], {desiredSize: 101});

    t.deepEqual(allocated.stdout, {from: 0, to:100});

    await cleanup()
    t.end();
});

test('vlan range strategy range partly claimed', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await findAllocationStrategyId('vlan_range');
    const claimed = [
        {
            Properties: {from: 0, to: 1000},
            Status: 'claimed',
            UpdatedAt: '2020-08-30 11:38:48.0 +0200 CEST'
        },];

        let allocated = await testStrategy(strategyId, {from: 0, to: 4095},
        poolName, poolId, claimed, {desiredSize: 101});

    t.deepEqual(allocated.stdout, {from: 1001, to:1101});

    await cleanup()
    t.end();
});

test('vlan strategy', async (t) => {
    let poolName = getUniqueName('testJSstrategy');
    let poolId = 0;
    let strategyId = await findAllocationStrategyId('vlan');
    let allocated = await testStrategy(strategyId, {from: 0, to: 4095},
        poolName, poolId, [], {});

    t.equal(allocated.stdout.vlan, 0);

    await cleanup()
    t.end();
});
