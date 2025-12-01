import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { RegisterPage } from './RegisterPage';
import { AuthContext } from '../context/AuthContext';
import { BrowserRouter } from 'react-router-dom';
import { authApi } from '../api/auth';

// Mock authApi
vi.mock('../api/auth', () => ({
    authApi: {
        register: vi.fn(),
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

describe('RegisterPage', () => {
    const loginMock = vi.fn();

    const renderComponent = () => {
        render(
            <AuthContext.Provider value={{ currentUser: null, login: loginMock, logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <RegisterPage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders registration form', () => {
        renderComponent();
        expect(screen.getByLabelText(/full name/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /register/i })).toBeInTheDocument();
    });

    it('handles successful registration', async () => {
        const user = userEvent.setup();
        const mockUser = { id: 1, name: 'Test User', email: 'test@test.com', is_admin: false };
        (authApi.register as any).mockResolvedValue(mockUser);

        renderComponent();

        await user.type(screen.getByLabelText(/full name/i), 'Test User');
        await user.type(screen.getByLabelText(/email/i), 'test@test.com');
        await user.type(screen.getByLabelText(/password/i), 'password123');
        await user.click(screen.getByRole('button', { name: /register/i }));

        await waitFor(() => {
            expect(authApi.register).toHaveBeenCalledWith({
                name: 'Test User',
                email: 'test@test.com',
                password: 'password123'
            });
            expect(navigateMock).toHaveBeenCalledWith('/login');
        });
    });

    it('displays error message on registration failure', async () => {
        const user = userEvent.setup();
        (authApi.register as any).mockRejectedValue({ response: { data: { message: 'Email already exists' } } });

        renderComponent();

        await user.type(screen.getByLabelText(/full name/i), 'Test User');
        await user.type(screen.getByLabelText(/email/i), 'existing@test.com');
        await user.type(screen.getByLabelText(/password/i), 'password123');
        await user.click(screen.getByRole('button', { name: /register/i }));

        expect(await screen.findByText(/email already exists/i)).toBeInTheDocument();
    });
});
