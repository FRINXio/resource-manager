import {claimResource} from '../graphql-queries';
import {createIpv4RootPool, createRandomIntRootPool, createRdRootPool} from '../test-helpers';

test('create rd/ipv4/random root pool', async () => {
    expect(await createRdRootPool()).toBeTruthy();
    expect(await createIpv4RootPool('192.168.1.0', 24)).toBeTruthy();
    expect(await createRandomIntRootPool()).toBeTruthy();
});

test('create AS and RD', async () => {
    const AS = 4545;
    const randomPoolId = await createRandomIntRootPool();
    const rdPoolId = await createRdRootPool();

    let randomNumber = (await claimResource(randomPoolId, {})).Properties.int;
    let rd = await claimResource(rdPoolId, {asNumber: 4545, assignedNumber: randomNumber});
    expect(rd.Properties.rd).toBe(`${AS}:${randomNumber}`);
});

test('create ipv4 and RD', async () => {
    const randomPoolId = await createRandomIntRootPool();
    const rdPoolId = await createRdRootPool();
    const ipv4PoolId = await createIpv4RootPool('192.168.1.0', 24);

    let ipv4 = (await claimResource(ipv4PoolId, {subnet: true})).Properties.address;
    let randomNumber = (await claimResource(randomPoolId, {})).Properties.int;
    let rd2 = await claimResource(rdPoolId, {ipv4: ipv4, assignedNumber: randomNumber});
    expect(rd2.Properties.rd).toBe(`${ipv4}:${randomNumber}`);
});
