import {claimResource} from '../graphql-queries.js';
import {createIpv4RootPool, createRandomIntRootPool, createRdRootPool} from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;

test('create rd/ipv4/random root pool', async (t) => {
    t.ok(await createRdRootPool());
    t.ok(await createIpv4RootPool('192.168.1.0', 24));
    t.ok(await createRandomIntRootPool());
    t.end();
});

test('create AS and RD', async (t) => {
    const AS = 4545;
    const randomPoolId = await createRandomIntRootPool();
    const rdPoolId = await createRdRootPool();

    let randomNumber = (await claimResource(randomPoolId, {})).Properties.int;
    let rd = await claimResource(rdPoolId, {asNumber: 4545, assignedNumber: randomNumber});
    t.equal(rd.Properties.rd, `${AS}:${randomNumber}`);
    t.end();
});

test('create ipv4 and RD', async (t) => {
    const randomPoolId = await createRandomIntRootPool();
    const rdPoolId = await createRdRootPool();
    const ipv4PoolId = await createIpv4RootPool('192.168.1.0', 24);

    let ipv4 = (await claimResource(ipv4PoolId, {subnet: true})).Properties.address;
    let randomNumber = (await claimResource(randomPoolId, {})).Properties.int;
    let rd2 = await claimResource(rdPoolId, {ipv4: ipv4, assignedNumber: randomNumber});
    t.equal(rd2.Properties.rd, `${ipv4}:${randomNumber}`);
    t.end();
});
