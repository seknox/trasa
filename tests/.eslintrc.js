module.exports = {
    extends: [
      "plugin:jest/recommended",
      "plugin:jest/style",
    ],
    env: {
      browser: true,
      commonjs: true,
      es6: true,
      jest: true,
      node: true
    },
    globals: {
        page: true,
        browser: true,
        context: true,
        jestPuppeteer: true,
      },
  
      'max-len': [
        'warn',
        {
          code: 100,
          tabWidth: 2,
          comments: 100,
          ignoreComments: false,
          ignoreTrailingComments: true,
          ignoreUrls: true,
          ignoreStrings: true,
          ignoreTemplateLiterals: true,
          ignoreRegExpLiterals: true
        }
      ]
    }
  };
  