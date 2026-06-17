---
name: nuxt-developer
description: Develop web applications using Nuxt.js, a powerful Vue.js framework for server-side rendering and static site generation and test with playwright.
---

## Technology Stack
| Concern        | Choice                                      |
|----------------|---------------------------------------------|
| Language       | JavaScript / TypeScript                     |
| Framework      | [Nuxt.js](https://nuxtjs.org/)             |
| UI Library     | Nuxt UI              |
| State Management | [Pinia](https://pinia.vuejs.org/)        |
| HTTP Client    | [Axios](https://axios-http.com/)           |
| Testing        | [Playwright](https://playwright.dev/)      |

## Project Structure
```
├── components/          # Reusable Vue components
├── layouts/             # Application layouts
├── pages/               # Application pages (routes)
├── plugins/             # Nuxt plugins
├── store/               # Pinia stores for state management
├── assets/              # Static assets (images, styles, etc.)
├── middleware/          # Middleware for route handling
├── tests/               # Playwright tests
├── nuxt.config.js       # Nuxt configuration file
```

## Workflow of Nuxt developer
1. Understand the requirements and design the application structure
2. Read HTML template from input and convert it into Nuxt components and compare with the existing codebase to identify missing components and features
3. Implement the necessary components, pages, and state management logic using Nuxt.js and Pinia
4. Add test ids to the components for testing with playwright
5. Integrate with backend APIs using Axios for data fetching and manipulation
6. Write Playwright tests to ensure the application works as expected and has good test coverage
   - Must use network interception to mock API responses in tests
   - Must passed all tests in `tests/` directory, if failed then fix the code and make sure all tests passed

## Rules
- Always write clean and maintainable code following best practices
- Ensure that all implemented features are covered by tests
- Regularly commit code changes with meaningful commit messages
- Collaborate with team members and seek feedback to improve the code quality
- Continuously learn and stay updated with the latest Nuxt.js features and best practices