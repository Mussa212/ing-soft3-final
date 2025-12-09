describe('Reservations Flow', () => {
    const API_URL = Cypress.env('API_BASE_URL') || 'http://localhost:8080';

    const uniqueEmail = `test${Date.now()}@example.com`;

    beforeEach(() => {
        // Intercept requests to alias them, but do NOT stub the response.
        // This allows the request to go through to the real backend.
        cy.intercept('POST', '/auth/register').as('register');
        cy.intercept('POST', '/auth/login').as('login');
        cy.intercept('POST', '/reservations').as('createReservation');
        cy.intercept('GET', '/my/reservations*').as('getReservations');

        // Set viewport to desktop to ensure images are visible
        cy.viewport(1280, 720);
    });

    it('allows a user to register, login, and create a reservation', () => {
        // Home Page Verification & Navigation
        cy.visit('/');
        cy.contains('h1', 'Vesuvio Ristorante').should('be.visible');
        cy.get('img[alt="Artie Bucco"]').should('be.visible');
        cy.get('img[alt="Vesuvio Restaurant"]').should('be.visible');

        cy.contains('button', 'Register').click();
        cy.url().should('include', '/register');

        // Register
        cy.contains('label', 'Full Name').parent().find('input').type('Test User');
        cy.contains('label', 'Email').parent().find('input').type(uniqueEmail);
        cy.contains('label', 'Password').parent().find('input').type('password123');
        cy.contains('button', 'Register').click();

        cy.wait('@register').its('response.statusCode').should('eq', 201);
        cy.url().should('include', '/login');

        // Login
        cy.contains('label', 'Email').parent().find('input').type(uniqueEmail);
        cy.contains('label', 'Password').parent().find('input').type('password123');
        cy.contains('button', 'Sign In').click();

        cy.wait('@login').its('response.statusCode').should('eq', 200);
        cy.url().should('include', '/reservations');

        // Create Reservation
        cy.contains('a', 'New Reservation').click();
        cy.url().should('include', '/reservations/new');

        cy.contains('label', 'Date').parent().find('input').type('2025-12-25');
        cy.contains('label', 'Time').parent().find('input').type('20:00');
        cy.contains('label', 'Number of People').parent().find('input').clear().type('2');
        cy.contains('button', 'Confirm Reservation').click();

        cy.wait('@createReservation').its('response.statusCode').should('eq', 201);
        cy.url().should('include', '/reservations');

        // Verify list
        cy.wait('@getReservations');
        cy.contains('20:00').should('be.visible');
        cy.contains('PENDING').should('be.visible');

        // Cancel Reservation
        cy.intercept('PATCH', '/reservations/*/cancel').as('cancelReservation');

        // We need to reload or wait for the list to ensure we have the latest data if needed, 
        // but typically we just act on what's there.
        // Note: The previous test mocked the cancelled response. Now we rely on the backend.

        cy.contains('button', 'Cancel').click();
        cy.wait('@cancelReservation').its('response.statusCode').should('eq', 200);

        // Verify status change - might need a reload or the UI updates automatically
        // Assuming UI updates automatically on success
        cy.contains('CANCELLED').should('be.visible');
    });

    it('admin flow', () => {
        cy.intercept('POST', '/auth/login').as('adminLogin');
        cy.intercept('GET', '/admin/reservations*').as('getAdminReservations');
        cy.intercept('PATCH', '/admin/reservations/*/confirm').as('confirmReservation');
        cy.intercept('PATCH', '/admin/reservations/*/cancel').as('adminCancelReservation');

        // Login as admin (via Home)
        cy.visit('/');
        cy.contains('button', 'Login').click();
        cy.url().should('include', '/login');

        // Assuming these admin credentials exist in the real backend
        cy.contains('label', 'Email').parent().find('input').type('admin@vesuvio.test');
        cy.contains('label', 'Password').parent().find('input').type('ChangeMe123!');
        cy.contains('button', 'Sign In').click();

        cy.wait('@adminLogin').its('response.statusCode').should('eq', 200);
        cy.url().should('include', '/admin/reservations');

        // Wait for the admin reservations API call
        cy.wait('@getAdminReservations');

        // Verify we're on the admin panel - the title is "Admin Panel"
        cy.contains('h1', 'Admin Panel').should('be.visible');

        // The page uses ReservationsList component which renders cards, not a table.
        // If there are reservations for today, we'll see them. Otherwise we'll see "No reservations found."
        // Either outcome is valid for this test - we just verify the page loaded.
        cy.get('body').should('not.contain', 'Loading reservations...');
    });
});
