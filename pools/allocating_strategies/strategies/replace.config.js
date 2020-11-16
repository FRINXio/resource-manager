module.exports = {
    encoding: 'utf8',
    from: [/var currentResources = \[\]/g, /var resourcePoolProperties = {}/g, /var userInput = {}/g],
    to: '//',
    files: [
        'generated/**/*.js',
    ],
};
