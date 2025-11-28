import { httpClient } from './httpClient';
import type { CreateReservationDto, Reservation } from '../types';

export const reservationsApi = {
    create: async (data: CreateReservationDto) => {
        const response = await httpClient.post<Reservation>('/reservations', data);
        return response.data;
    },

    getMyReservations: async (status?: string) => {
        const params = status && status !== 'all' ? { status } : {};
        const response = await httpClient.get<Reservation[]>('/my/reservations', { params });
        return response.data;
    },

    cancel: async (id: number) => {
        const response = await httpClient.patch<Reservation>(`/reservations/${id}/cancel`);
        return response.data;
    },

    // Admin endpoints
    getAll: async (date?: string, status?: string) => {
        const params: Record<string, string> = {};
        if (date) params.date = date;
        if (status && status !== 'all') params.status = status;

        const response = await httpClient.get<Reservation[]>('/admin/reservations', { params });
        return response.data;
    },

    confirm: async (id: number) => {
        const response = await httpClient.patch<Reservation>(`/admin/reservations/${id}/confirm`);
        return response.data;
    },

    adminCancel: async (id: number) => {
        const response = await httpClient.patch<Reservation>(`/admin/reservations/${id}/cancel`);
        return response.data;
    },
};
