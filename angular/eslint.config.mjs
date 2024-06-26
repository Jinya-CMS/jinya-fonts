import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import angularEslint from 'angular-eslint';

export default tseslint.config({
  files: ['**/*.ts', '**/*.html'],
  extends: [
    js.configs.recommended,
    ...tseslint.configs.recommended,
    ...tseslint.configs.stylistic,
    ...angularEslint.configs.tsRecommended
  ],
  rules: {
    '@angular-eslint/directive-selector': ['error', {
      type: 'attribute',
      prefix: 'app',
      style: 'camelCase'
    }],
    '@angular-eslint/component-selector': ['error', {
      type: 'element',
      prefix: 'app',
      style: 'kebab-case'
    }]
  }
});
//
// export default [{
//   ignores: ['projects/**/*']
// }, ...compat.extends(
//   'eslint:recommended',
//   'plugin:@typescript-eslint/recommended',
//   'plugin:@angular-eslint/recommended',
//   'plugin:@angular-eslint/template/process-inline-templates'
// ).map(config => ({
//   ...config,
//   files: ['**/*.ts']
// })), {
//   files: ['**/*.ts'],
//
//   rules: {
//     '@angular-eslint/directive-selector': ['error', {
//       type: 'attribute',
//       prefix: 'app',
//       style: 'camelCase'
//     }],
//
//     '@angular-eslint/component-selector': ['error', {
//       type: 'element',
//       prefix: 'app',
//       style: 'kebab-case'
//     }]
//   }
// }, ...compat.extends(
//   'plugin:@angular-eslint/template/recommended',
//   'plugin:@angular-eslint/template/accessibility'
// ).map(config => ({
//   ...config,
//   files: ['**/*.html']
// })), {
//   files: ['**/*.html'],
//   rules: {}
// }];
