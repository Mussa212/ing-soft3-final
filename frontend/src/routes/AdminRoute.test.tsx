import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { AdminRoute } from './AdminRoute';
import { useAuth } from '../hooks/useAuth';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

vi.mock('../hooks/useAuth');

describe('AdminRoute', () => {
    it('renders loading state', () => {
        (useAuth as any).mockReturnValue({ isLoading: true });
        render(<AdminRoute />);
        expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('redirects to login if not authenticated', () => {
        (useAuth as any).mockReturnValue({ isLoading: false, currentUser: null });
        render(
            <MemoryRouter initialEntries={['/admin']}>
                <Routes>
                    <Route path="/admin" element={<AdminRoute />} />
                    <Route path="/login" element={<div>Login Page</div>} />
                </Routes>
            </MemoryRouter>
        );
        expect(screen.getByText('Login Page')).toBeInTheDocument();
    });

    it('redirects to reservations if authenticated but not admin', () => {
        (useAuth as any).mockReturnValue({ isLoading: false, currentUser: { id: 1, is_admin: false } });
        render(
            <MemoryRouter initialEntries={['/admin']}>
                <Routes>
                    <Route path="/admin" element={<AdminRoute />} />
                    <Route path="/reservations" element={<div>Reservations Page</div>} />
                </Routes>
            </MemoryRouter>
        );
        expect(screen.getByText('Reservations Page')).toBeInTheDocument();
    });

    it('renders outlet if authenticated and admin', () => {
        (useAuth as any).mockReturnValue({ isLoading: false, currentUser: { id: 1, is_admin: true } });
        render(
            <MemoryRouter initialEntries={['/admin']}>
                <Routes>
                    <Route path="/admin" element={<AdminRoute />}>
                        <Route path="" element={<div>Admin Content</div>} />
                    </Route>
                </Routes>
            </MemoryRouter>
        );
        expect(screen.getByText('Admin Content')).toBeInTheDocument();
    });
});
