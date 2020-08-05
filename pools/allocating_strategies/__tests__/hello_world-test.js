test('hello world strategy returns always the same resource', () => {
    let script = require('../hello_world_strategy');
    const val1 = script.invoke()
    const val2 = script.invoke()

    expect(val1).toStrictEqual(val2);
    expect(val1).toStrictEqual({"message": "Hello World!"});
});