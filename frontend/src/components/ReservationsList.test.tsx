import { render, screen, fireEvent } from '@testing-library/react';
import { vi, describe, it, expect } from 'vitest';
import { ReservationsList } from './ReservationsList';
import type { Reservation } from '../types';

describe('ReservationsList', () => {
    const mockReservations: Reservation[] = [
        { id: 1, user_id: 1, date: '2025-12-25', time: '20:00', people: 4, status: 'pending' },
        { id: 2, user_id: 1, date: '2025-12-26', time: '19:00', people: 2, status: 'confirmed' },
        { id: 3, user_id: 1, date: '2025-12-27', time: '21:00', people: 3, status: 'cancelled' },
    ];

    it('renders reservations', () => {
        render(<ReservationsList reservations={mockReservations} onCancel={vi.fn()} />);

        expect(screen.getByText(/20:00/)).toBeInTheDocument();
        expect(screen.getByText(/19:00/)).toBeInTheDocument();
        expect(screen.getByText(/21:00/)).toBeInTheDocument();
        expect(screen.getByText(/PENDING/)).toBeInTheDocument();
        expect(screen.getByText(/CONFIRMED/)).toBeInTheDocument();
        expect(screen.getByText(/CANCELLED/)).toBeInTheDocument();
    });

    it('shows cancel button only for pending/confirmed', () => {
        render(<ReservationsList reservations={mockReservations} onCancel={vi.fn()} />);

        const cancelButtons = screen.getAllByRole('button', { name: /cancel/i });
        expect(cancelButtons).toHaveLength(2); // Only for pending and confirmed
    });

    it('calls onCancel when cancel button is clicked', () => {
        const onCancelMock = vi.fn();
        render(<ReservationsList reservations={[mockReservations[0]]} onCancel={onCancelMock} />);

        fireEvent.click(screen.getByRole('button', { name: /cancel/i }));
        expect(onCancelMock).toHaveBeenCalledWith(1);
    });

    it('shows confirm button for admin on pending reservations', () => {
        const onConfirmMock = vi.fn();
        render(
            <ReservationsList
                reservations={[mockReservations[0]]}
                onCancel={vi.fn()}
                isAdmin={true}
                onConfirm={onConfirmMock}
            />
        );

        expect(screen.getByText(/confirm/i)).toBeInTheDocument();
        fireEvent.click(screen.getByText(/confirm/i));
        expect(onConfirmMock).toHaveBeenCalledWith(1);
    });
});
