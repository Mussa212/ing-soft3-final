export interface User {
    id: number;
    name: string;
    email: string;
    is_admin: boolean;
}

export interface Reservation {
    id: number;
    user_id: number;
    date: string;
    time: string;
    people: number;
    comment?: string;
    status: 'pending' | 'confirmed' | 'cancelled';
    user?: User; // For admin view
}

export interface CreateReservationDto {
    date: string;
    time: string;
    people: number;
    comment?: string;
}

export interface LoginResponse extends User { }
