import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { PrivateRoute } from './PrivateRoute';
import { useAuth } from '../hooks/useAuth';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

vi.mock('../hooks/useAuth');

describe('PrivateRoute', () => {
    it('renders loading state', () => {
        (useAuth as any).mockReturnValue({ isLoading: true });
        render(<PrivateRoute />);
        expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('redirects to login if not authenticated', () => {
        (useAuth as any).mockReturnValue({ isLoading: false, currentUser: null });
        render(
            <MemoryRouter initialEntries={['/private']}>
                <Routes>
                    <Route path="/private" element={<PrivateRoute />} />
                    <Route path="/login" element={<div>Login Page</div>} />
                </Routes>
            </MemoryRouter>
        );
        expect(screen.getByText('Login Page')).toBeInTheDocument();
    });

    it('renders outlet if authenticated', () => {
        (useAuth as any).mockReturnValue({ isLoading: false, currentUser: { id: 1 } });
        render(
            <MemoryRouter initialEntries={['/private']}>
                <Routes>
                    <Route path="/private" element={<PrivateRoute />}>
                        <Route path="" element={<div>Private Content</div>} />
                    </Route>
                </Routes>
            </MemoryRouter>
        );
        expect(screen.getByText('Private Content')).toBeInTheDocument();
    });
});
