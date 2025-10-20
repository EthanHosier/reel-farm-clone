# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

## API Client Generation

This project uses OpenAPI TypeScript code generation to create type-safe API clients from the backend OpenAPI specification.

### Generating the API Client

To generate the TypeScript API client from the OpenAPI specification:

```bash
npm run generate-api
```

This command will:
- Read the OpenAPI spec from `../openapi/api.yaml`
- Generate TypeScript types and service classes in `./src/api/`
- Create type-safe API clients for all backend endpoints

### API Client Usage

The generated API client provides:
- **Type-safe models**: `UserAccount`, `HealthResponse`, `ErrorResponse`
- **Service classes**: `UsersService`, `HealthService`
- **Core utilities**: `OpenAPI`, `ApiError`, `CancelablePromise`

Example usage:
```typescript
import { UsersService, OpenAPI } from './api';

// Configure the base URL and authentication
OpenAPI.BASE = 'http://localhost:3000';
OpenAPI.TOKEN = 'Bearer your-jwt-token';

// Make API calls
const userAccount = await UsersService.getUserAccount();
```

### Regenerating After Backend Changes

Whenever the backend API changes (new endpoints, modified schemas, etc.), regenerate the client:

```bash
npm run generate-api
```

This ensures your frontend stays in sync with the backend API.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) (or [oxc](https://oxc.rs) when used in [rolldown-vite](https://vite.dev/guide/rolldown)) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## React Compiler

The React Compiler is enabled on this template. See [this documentation](https://react.dev/learn/react-compiler) for more information.

Note: This will impact Vite dev & build performances.

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...

      // Remove tseslint.configs.recommended and replace with this
      tseslint.configs.recommendedTypeChecked,
      // Alternatively, use this for stricter rules
      tseslint.configs.strictTypeChecked,
      // Optionally, add this for stylistic rules
      tseslint.configs.stylisticTypeChecked,

      // Other configs...
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
