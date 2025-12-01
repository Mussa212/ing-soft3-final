import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MyReservationsPage } from './MyReservationsPage';
import { AuthContext } from '../context/AuthContext';
import { BrowserRouter } from 'react-router-dom';
import { reservationsApi } from '../api/reservations';

// Mock reservationsApi
vi.mock('../api/reservations', () => ({
    reservationsApi: {
        getMyReservations: vi.fn(),
        cancel: vi.fn(),
    },
}));

describe('MyReservationsPage', () => {
    const mockReservations = [
        { id: 1, date: '2025-12-25', time: '20:00', people: 2, status: 'pending' },
        { id: 2, date: '2025-12-26', time: '19:00', people: 4, status: 'confirmed' },
    ];

    const renderComponent = () => {
        render(
            <AuthContext.Provider value={{ currentUser: { id: 1, name: 'User', email: 'u@u.com', is_admin: false }, login: vi.fn(), logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <MyReservationsPage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders reservations list', async () => {
        (reservationsApi.getMyReservations as any).mockResolvedValue(mockReservations);
        renderComponent();

        expect(screen.getByRole('heading', { name: /my reservations/i })).toBeInTheDocument();
        await waitFor(() => {
            expect(screen.getAllByText(/2025/)).toHaveLength(3); // 2 reservations + 1 footer
            expect(screen.getAllByText(/people/i)).toHaveLength(2);
            expect(screen.getAllByRole('article')).toHaveLength(2);
        });
    });

    it('handles cancellation', async () => {
        (reservationsApi.getMyReservations as any).mockResolvedValue(mockReservations);
        (reservationsApi.cancel as any).mockResolvedValue({});

        // Mock window.confirm
        vi.spyOn(window, 'confirm').mockImplementation(() => true);

        renderComponent();

        await waitFor(() => {
            expect(screen.getAllByRole('button', { name: /cancel/i })).toHaveLength(2);
        });

        fireEvent.click(screen.getAllByRole('button', { name: /cancel/i })[0]);

        expect(window.confirm).toHaveBeenCalled();
        await waitFor(() => {
            expect(reservationsApi.cancel).toHaveBeenCalledWith(1);
            expect(reservationsApi.getMyReservations).toHaveBeenCalledTimes(2); // Initial + after cancel
        });
    });
});
