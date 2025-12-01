import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { NewReservationPage } from './NewReservationPage';
import { AuthContext } from '../context/AuthContext';
import { BrowserRouter } from 'react-router-dom';
import { reservationsApi } from '../api/reservations';

// Mock reservationsApi
vi.mock('../api/reservations', () => ({
    reservationsApi: {
        create: vi.fn(),
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

describe('NewReservationPage', () => {
    const renderComponent = () => {
        render(
            <AuthContext.Provider value={{ currentUser: { id: 1, name: 'User', email: 'u@u.com', is_admin: false }, login: vi.fn(), logout: vi.fn(), isLoading: false }}>
                <BrowserRouter>
                    <NewReservationPage />
                </BrowserRouter>
            </AuthContext.Provider>
        );
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders the form', () => {
        renderComponent();
        expect(screen.getByLabelText(/date/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/time/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/number of people/i)).toBeInTheDocument();
    });

    it('validates people count', async () => {
        const user = userEvent.setup();
        renderComponent();

        const peopleInput = screen.getByLabelText(/number of people/i);
        await user.clear(peopleInput);
        await user.type(peopleInput, '0');

        await user.click(screen.getByRole('button', { name: /confirm reservation/i }));

        expect(await screen.findByText(/number of people must be greater than 0/i)).toBeInTheDocument();
        expect(reservationsApi.create).not.toHaveBeenCalled();
    });

    it('submits valid form', async () => {
        const user = userEvent.setup();
        (reservationsApi.create as any).mockResolvedValue({});

        renderComponent();

        const dateInput = screen.getByLabelText(/date/i);
        const timeInput = screen.getByLabelText(/time/i);
        const peopleInput = screen.getByLabelText(/number of people/i);

        // Date input needs strict format yyyy-mm-dd
        fireEvent.change(dateInput, { target: { value: '2025-12-25' } });
        fireEvent.change(timeInput, { target: { value: '20:00' } });
        await user.clear(peopleInput);
        await user.type(peopleInput, '4');

        await user.click(screen.getByRole('button', { name: /confirm reservation/i }));

        await waitFor(() => {
            expect(reservationsApi.create).toHaveBeenCalledWith({
                date: '2025-12-25',
                time: '20:00',
                people: 4,
                comment: ''
            });
            expect(navigateMock).toHaveBeenCalledWith('/reservations');
        });
    });
});
