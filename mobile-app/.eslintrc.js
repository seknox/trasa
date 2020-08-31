module.exports = {
  extends: [
    'airbnb',
    // 'airbnb-typescript/base',
    'plugin:prettier/recommended',
    'prettier/react',
    'airbnb/hooks',
    'eslint:recommended',
    'plugin:import/errors',
    'plugin:import/warnings',
    //  'plugin:import/typescript',
  ],
  env: {
    browser: true,
    commonjs: true,
    es6: true,
    jest: true,
    node: true,
  },
  // parser: "babel-eslint",
  parserOptions: {
    project: './tsconfig.json',
  },
  // overrides:[
  //   {rules: {
  //     'react/destructuring-assignment': ['off']
  //   }
  // ]
  rules: {
    'react/destructuring-assignment': ['off'],
    'react/prop-types': ['off'],
    'react-hooks/rules-of-hooks': 'error', // Checks rules of Hooks
    'react-hooks/exhaustive-deps': 'warn', // Checks effect dependencies
    'jsx-a11y/href-no-hash': ['off'],
    'react/jsx-filename-extension': ['warn', { extensions: ['.js', '.jsx', '.ts', '.tsx'] }],
    // 'import/parsers': [
    //   'warn',
    //   {
    //     '@typescript-eslint/parser': ['.ts', '.tsx', '.js', '.jsx']
    //   }
    // ],

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
        ignoreRegExpLiterals: true,
      },
    ],
  },
};
