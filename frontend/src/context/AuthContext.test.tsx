import { render, screen, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthProvider, AuthContext } from './AuthContext';
import React, { useContext } from 'react';

// Test component to consume context
const TestComponent = () => {
    const { currentUser, login, logout, isLoading } = useContext(AuthContext)!;
    return (
        <div>
            <div data-testid="user">{currentUser ? currentUser.name : 'No User'}</div>
            <div data-testid="loading">{isLoading.toString()}</div>
            <button onClick={() => login({ id: 1, name: 'Test User', email: 'test@test.com', is_admin: false })}>Login</button>
            <button onClick={logout}>Logout</button>
        </div>
    );
};

describe('AuthContext', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        localStorage.clear();
    });

    it('provides default state', () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );
        expect(screen.getByTestId('user')).toHaveTextContent('No User');
        expect(screen.getByTestId('loading')).toHaveTextContent('false');
    });

    it('updates state on login', async () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await act(async () => {
            screen.getByText('Login').click();
        });

        expect(screen.getByTestId('user')).toHaveTextContent('Test User');
        expect(localStorage.getItem('currentUser')).toBeTruthy();
    });

    it('updates state on logout', async () => {
        // Setup initial state with user
        const mockUser = { id: 1, name: 'Test User', email: 'test@test.com', is_admin: false };
        localStorage.setItem('currentUser', JSON.stringify(mockUser));

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        // Should initialize with user from localStorage
        expect(screen.getByTestId('user')).toHaveTextContent('Test User');

        await act(async () => {
            screen.getByText('Logout').click();
        });

        expect(screen.getByTestId('user')).toHaveTextContent('No User');
        expect(localStorage.getItem('currentUser')).toBeNull();
    });
});
