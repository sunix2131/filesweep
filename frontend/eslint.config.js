import js from '@eslint/js';
import tseslint from '@typescript-eslint/eslint-plugin';
import parser from '@typescript-eslint/parser';
import hooks from 'eslint-plugin-react-hooks';

export default [
  {
    ignores: ['dist/**', 'wailsjs/**', 'node_modules/**']
  },
  js.configs.recommended,
  {
    files: ['src/**/*.{ts,tsx}', 'tests/**/*.ts'],
    languageOptions: {
      parser,
      parserOptions: { ecmaFeatures: { jsx: true }, sourceType: 'module' },
      globals: {
        window: 'readonly',
        document: 'readonly',
        console: 'readonly'
      }
    },
    plugins: { '@typescript-eslint': tseslint, 'react-hooks': hooks },
    rules: { ...tseslint.configs.recommended.rules, ...hooks.configs.recommended.rules }
  }
];
