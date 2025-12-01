# Vesuvio Ristorante Frontend

A React + TypeScript frontend for the Vesuvio Restaurant reservation system.

## Prerequisites

- Node.js (v18 or later)
- npm

## Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

## Running the Application

Start the development server:
```bash
npm run dev
```
The application will be available at `http://localhost:5173`.

## Testing

### Unit Tests
Run all unit tests with Vitest:
```bash
npm test
```

### End-to-End Tests
Run Cypress E2E tests:
```bash
npx cypress run
```
Or open the Cypress UI:
```bash
npx cypress open
```

## Features

- **Authentication**: Login and Registration.
- **Client**: Create and view reservations.
- **Admin**: Manage all reservations (confirm/cancel).
- **Responsive Design**: Mobile-friendly UI.
