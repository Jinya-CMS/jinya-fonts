const js = require("@eslint/js");
const prettier = require("eslint-config-prettier");
const globals = require("globals");

module.exports = [
  {
    ...js.configs.recommended,
    ...prettier,
    files: ["/static/admin/js/**/*.js", "/static/frontend/js/**/*.js"],
    languageOptions: {
      ecmaVersion: "latest",
      sourceType: "module",
      globals: {
        ...globals.browser,
        Jodit: "readonly",
      },
    },
  },
];
