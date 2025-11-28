import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AdminReservationsPage } from './AdminReservationsPage';
import { AuthContext } from '../context/AuthContext';
import { BrowserRouter } from 'react-router-dom';
import { reservationsApi } from '../api/reservations';

// Mock reservationsApi
vi.mock('../api/reservations', () => ({
    reservationsApi: {
        getAll: vi.fn(),
        confirm: vi.fn(),
        adminCancel: vi.fn(),
    },
}));

describe('AdminReservationsPage', () => {
    const mockReservations = [
        { id: 1, date: '2025-12-25', time: '20:00', people: 2, status: 'pending', user: { name: 'User 1' } },
        { id: 2, date: '2025-12-26', time: '19:00', people: 4, status: 'confirmed', user: { name: 'User 2' } },
    ];

    const renderComponent = () => {
        render(
            <AuthContext.Provider value={{ currentUser: { id: 1, name: 'Admin', email: 'a@a.com', is_admin: true }, login: vi.fn(), logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <AdminReservationsPage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders admin panel and reservations', async () => {
        (reservationsApi.getAll as any).mockResolvedValue(mockReservations);
        renderComponent();

        expect(screen.getByRole('heading', { name: /admin panel/i })).toBeInTheDocument();
        await waitFor(() => {
            expect(screen.getAllByText(/user 1/i).length).toBeGreaterThan(0);
            expect(screen.getAllByText(/user 2/i).length).toBeGreaterThan(0);
            expect(screen.getAllByRole('article')).toHaveLength(2);
        });
    });

    it('handles confirmation', async () => {
        (reservationsApi.getAll as any).mockResolvedValue(mockReservations);
        (reservationsApi.confirm as any).mockResolvedValue({});

        renderComponent();

        await waitFor(() => {
            expect(screen.getByRole('button', { name: /confirm/i })).toBeInTheDocument();
        });

        fireEvent.click(screen.getByRole('button', { name: /confirm/i }));

        await waitFor(() => {
            expect(reservationsApi.confirm).toHaveBeenCalledWith(1);
            expect(reservationsApi.getAll).toHaveBeenCalledTimes(2);
        });
    });

    it('handles cancellation', async () => {
        (reservationsApi.getAll as any).mockResolvedValue(mockReservations);
        (reservationsApi.adminCancel as any).mockResolvedValue({});

        vi.spyOn(window, 'confirm').mockImplementation(() => true);

        renderComponent();

        await waitFor(() => {
            expect(screen.getAllByRole('button', { name: /cancel/i })).toHaveLength(2);
        });

        fireEvent.click(screen.getAllByRole('button', { name: /cancel/i })[0]);

        expect(window.confirm).toHaveBeenCalled();
        await waitFor(() => {
            expect(reservationsApi.adminCancel).toHaveBeenCalledWith(1);
            expect(reservationsApi.getAll).toHaveBeenCalledTimes(2);
        });
    });
});
