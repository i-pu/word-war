module.exports = {
  root: true,
  env: {
    node: true
  },
  // extends: [
  //   'plugin:vue/recommended',
  //   '@vue/prettier',
  //   '@vue/typescript'
  // ],
  rules: {
    // 'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-console': 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'prettier/prettier': [
      'error',
      {
        'semi': false,
        'singleQuote': true
      }
    ],
    // "vue/html-indent": [
    //   "error",
    //   {
    //     "attribute": 1,
    //     "baseIndent": 1,
    //     "closeBracket": 0,
    //     "alignAttributesVertically": true,
    //     "ignores": []
    //   },
    // ]
  },
  parserOptions: {
    parser: '@typescript-eslint/parser'
  },
  overrides: [
    {
      files: [
        '**/__tests__/*.{j,t}s?(x)',
        '**/tests/unit/**/*.spec.{j,t}s?(x)'
      ],
      env: {
        jest: true
      }
    }
  ]
}
