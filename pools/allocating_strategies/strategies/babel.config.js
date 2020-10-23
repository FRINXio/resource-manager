module.exports = {
    presets: [
        [
            "@babel/preset-env",
            {
                targets: {
                    node: "current",
                    chrome: "58",
                },
                corejs: 3,
                useBuiltIns: "entry",
            },
        ],
        '@babel/preset-flow',
    ],
    plugins: [
        "@babel/plugin-proposal-class-properties",
        "@babel/plugin-proposal-nullish-coalescing-operator",
        "@babel/plugin-proposal-optional-chaining",
        "@babel/plugin-transform-runtime",
    ],
    env: {
        test: {
            sourceMaps: "both",
        },
    },
};
