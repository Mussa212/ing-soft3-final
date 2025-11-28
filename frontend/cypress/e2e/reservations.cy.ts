describe('Reservations Flow', () => {
    const API_URL = Cypress.env('API_BASE_URL') || 'http://localhost:8080';

    beforeEach(() => {
        // Reset or mock backend if possible. For now we assume backend is running or we mock requests.
        // Since we don't have a real backend running in this environment, we will mock the API calls.

        cy.intercept('POST', '/auth/register', {
            statusCode: 201,
            body: { id: 1, name: 'Test User', email: 'test@example.com', is_admin: false }
        }).as('register');

        cy.intercept('POST', '/auth/login', {
            statusCode: 200,
            body: { id: 1, name: 'Test User', email: 'test@example.com', is_admin: false }
        }).as('login');

        cy.intercept('POST', '/reservations', {
            statusCode: 201,
            body: { id: 1, user_id: 1, date: '2025-12-25', time: '20:00', people: 2, status: 'pending' }
        }).as('createReservation');

        cy.intercept('GET', '/my/reservations*', {
            statusCode: 200,
            body: [
                { id: 1, user_id: 1, date: '2025-12-25', time: '20:00', people: 2, status: 'pending' }
            ]
        }).as('getReservations');
    });

    it('allows a user to register, login, and create a reservation', () => {
        // Register
        cy.visit('/register');
        cy.get('input[name="name"]').type('Test User'); // Assuming input id/name is name, or we use label
        // Check Input component implementation for id/name
        // Input component uses id passed or name.

        // Let's verify selectors based on our Input component
        // Input component: <input id={id || props.name} ... />
        // In RegisterPage: <Input label="Full Name" value={name} ... /> -> It doesn't have name prop, so id might be undefined?
        // Wait, Input component: const inputId = id || props.name;
        // If neither is provided, inputId is undefined. Label htmlFor will be undefined.
        // I should fix RegisterPage to provide id or name to Inputs.

        // Actually, let's fix RegisterPage and LoginPage to have name/id attributes for better accessibility and testing.
        // But for now, let's assume I can select by label text or I'll fix the code first.

        // Let's rely on testing-library like commands in Cypress if configured, or just standard cypress.
        // Standard cypress: cy.contains('label', 'Full Name').next('input').type(...)

        cy.contains('label', 'Full Name').parent().find('input').type('Test User');
        cy.contains('label', 'Email').parent().find('input').type('test@example.com');
        cy.contains('label', 'Password').parent().find('input').type('password123');
        cy.contains('button', 'Register').click();

        cy.wait('@register');
        cy.url().should('include', '/login');

        // Login
        cy.contains('label', 'Email').parent().find('input').type('test@example.com');
        cy.contains('label', 'Password').parent().find('input').type('password123');
        cy.contains('button', 'Sign In').click();

        cy.wait('@login');
        cy.url().should('include', '/reservations');

        // Create Reservation
        cy.contains('a', 'New Reservation').click();
        cy.url().should('include', '/reservations/new');

        cy.contains('label', 'Date').parent().find('input').type('2025-12-25');
        cy.contains('label', 'Time').parent().find('input').type('20:00');
        cy.contains('label', 'Number of People').parent().find('input').clear().type('2');
        cy.contains('button', 'Confirm Reservation').click();

        cy.wait('@createReservation');
        cy.url().should('include', '/reservations');

        // Verify list
        cy.wait('@getReservations');
        cy.contains('20:00').should('be.visible');
        cy.contains('PENDING').should('be.visible');

        // Cancel Reservation
        cy.intercept('PATCH', '/reservations/1/cancel', {
            statusCode: 200,
            body: { id: 1, status: 'cancelled' }
        }).as('cancelReservation');

        cy.intercept('GET', '/my/reservations*', {
            statusCode: 200,
            body: [
                { id: 1, user_id: 1, date: '2025-12-25', time: '20:00', people: 2, status: 'cancelled' }
            ]
        }).as('getReservationsCancelled');

        cy.contains('button', 'Cancel').click();
        cy.wait('@cancelReservation');
        cy.wait('@getReservationsCancelled');
        cy.contains('CANCELLED').should('be.visible');
    });

    it('admin flow', () => {
        cy.intercept('POST', '/auth/login', {
            statusCode: 200,
            body: { id: 2, name: 'Admin', email: 'admin@example.com', is_admin: true }
        }).as('adminLogin');

        cy.intercept('GET', '/admin/reservations*', {
            statusCode: 200,
            body: [
                {
                    id: 1,
                    user_id: 1,
                    date: '2025-12-25',
                    time: '20:00',
                    people: 2,
                    status: 'pending',
                    user: { id: 1, name: 'Test User', email: 'test@example.com', is_admin: false }
                }
            ]
        }).as('getAdminReservations');

        cy.intercept('PATCH', '/admin/reservations/1/confirm', {
            statusCode: 200,
            body: { id: 1, status: 'confirmed' }
        }).as('confirmReservation');

        // Login as admin
        cy.visit('/login');
        cy.contains('label', 'Email').parent().find('input').type('admin@example.com');
        cy.contains('label', 'Password').parent().find('input').type('admin123');
        cy.contains('button', 'Sign In').click();

        cy.wait('@adminLogin');
        cy.url().should('include', '/admin/reservations');

        // Check reservations
        cy.wait('@getAdminReservations');
        cy.contains('Test User').should('be.visible');

        // Confirm
        // We need to mock the refresh call which happens after confirm
        cy.intercept('GET', '/admin/reservations*', {
            statusCode: 200,
            body: [
                {
                    id: 1,
                    user_id: 1,
                    date: '2025-12-25',
                    time: '20:00',
                    people: 2,
                    status: 'confirmed',
                    user: { id: 1, name: 'Test User', email: 'test@example.com', is_admin: false }
                }
            ]
        }).as('getAdminReservationsConfirmed');

        cy.contains('button', 'Confirm').click();
        cy.wait('@confirmReservation');
        cy.wait('@getAdminReservationsConfirmed');

        cy.contains('CONFIRMED').should('be.visible');

        // Admin Cancel
        // Reset to pending for cancellation test or use a new one. 
        // For simplicity, let's assume we have another pending reservation or just cancel the confirmed one if allowed (logic might allow cancelling confirmed).
        // Let's assume we can cancel confirmed ones too as per logic.

        cy.intercept('PATCH', '/admin/reservations/1/cancel', {
            statusCode: 200,
            body: { id: 1, status: 'cancelled' }
        }).as('adminCancelReservation');

        cy.intercept('GET', '/admin/reservations*', {
            statusCode: 200,
            body: [
                {
                    id: 1,
                    user_id: 1,
                    date: '2025-12-25',
                    time: '20:00',
                    people: 2,
                    status: 'cancelled',
                    user: { id: 1, name: 'Test User', email: 'test@example.com', is_admin: false }
                }
            ]
        }).as('getAdminReservationsCancelled');

        // We need to handle the window.confirm
        cy.on('window:confirm', () => true);

        cy.contains('button', 'Cancel').click();
        cy.wait('@adminCancelReservation');
        cy.wait('@getAdminReservationsCancelled');
        cy.contains('CANCELLED').should('be.visible');
    });
});
