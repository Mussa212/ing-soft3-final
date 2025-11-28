import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { HomePage } from './HomePage';
import { AuthContext } from '../context/AuthContext';
import { BrowserRouter } from 'react-router-dom';

// Mock Navigate
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        Navigate: vi.fn(({ to }) => <div data-testid="navigate" data-to={to} />),
    };
});

describe('HomePage', () => {
    it('redirects to reservations if user is logged in', () => {
        render(
            <AuthContext.Provider value={{ currentUser: { id: 1, name: 'User', email: 'u@u.com', is_admin: false }, login: vi.fn(), logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <HomePage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
        expect(screen.getByTestId('navigate')).toHaveAttribute('data-to', '/reservations');
    });

    it('redirects to admin reservations if admin is logged in', () => {
        render(
            <AuthContext.Provider value={{ currentUser: { id: 1, name: 'Admin', email: 'a@a.com', is_admin: true }, login: vi.fn(), logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <HomePage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
        expect(screen.getByTestId('navigate')).toHaveAttribute('data-to', '/admin/reservations');
    });

    it('renders landing page content if not logged in', () => {
        render(
            <AuthContext.Provider value={{ currentUser: null, login: vi.fn(), logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <HomePage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
        expect(screen.getByRole('heading', { name: /vesuvio ristorante/i })).toBeInTheDocument();
        expect(screen.getByText(/authentic neapolitan cuisine/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /register/i })).toBeInTheDocument();
    });
});
