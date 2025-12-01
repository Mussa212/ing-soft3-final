import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { LoginPage } from './LoginPage';
import { AuthContext } from '../context/AuthContext';
import { BrowserRouter } from 'react-router-dom';
import { authApi } from '../api/auth';

// Mock authApi
vi.mock('../api/auth', () => ({
    authApi: {
        login: vi.fn(),
    },
}));

// Mock useNavigate
const navigateMock = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => navigateMock,
    };
});

describe('LoginPage', () => {
    const loginMock = vi.fn();

    const renderComponent = () => {
        render(
            <AuthContext.Provider value={{ currentUser: null, login: loginMock, logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <LoginPage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders the login form', () => {
        renderComponent();
        expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();

        // Check for images
        const artieImage = screen.getByAltText(/artie bucco/i);
        expect(artieImage).toBeInTheDocument();
        expect(artieImage).toHaveAttribute('src', expect.stringContaining('spoton.com'));

        const restaurantImage = screen.getByAltText(/vesuvio restaurant/i);
        expect(restaurantImage).toBeInTheDocument();
        expect(restaurantImage).toHaveAttribute('src', expect.stringContaining('spoton.com'));
    });

    it('calls login API and redirects on success (user)', async () => {
        const user = { id: 1, name: 'Test User', email: 'test@example.com', is_admin: false };
        (authApi.login as any).mockResolvedValue(user);

        renderComponent();

        fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password' } });
        fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

        await waitFor(() => {
            expect(authApi.login).toHaveBeenCalledWith({ email: 'test@example.com', password: 'password' });
            expect(loginMock).toHaveBeenCalledWith(user);
            expect(navigateMock).toHaveBeenCalledWith('/reservations');
        });
    });

    it('calls login API and redirects on success (admin)', async () => {
        const admin = { id: 2, name: 'Admin', email: 'admin@example.com', is_admin: true };
        (authApi.login as any).mockResolvedValue(admin);

        renderComponent();

        fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'admin@example.com' } });
        fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password' } });
        fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

        await waitFor(() => {
            expect(navigateMock).toHaveBeenCalledWith('/admin/reservations');
        });
    });

    it('shows error message on failure', async () => {
        (authApi.login as any).mockRejectedValue(new Error('Invalid credentials'));

        renderComponent();

        fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'wrong@example.com' } });
        fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'wrong' } });
        fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

        await waitFor(() => {
            expect(screen.getByText(/invalid credentials/i)).toBeInTheDocument();
        });
    });
});
