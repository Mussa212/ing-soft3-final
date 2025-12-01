import { describe, it, expect, vi, beforeEach } from 'vitest';
import { reservationsApi } from './reservations';
import { httpClient } from './httpClient';

vi.mock('./httpClient', () => ({
    httpClient: {
        get: vi.fn(),
        post: vi.fn(),
        patch: vi.fn(),
    },
}));

describe('reservationsApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('create calls httpClient.post', async () => {
        const mockData = { date: '2023-01-01', time: '12:00', people: 2, comment: 'test' };
        const mockResponse = { data: { id: 1, ...mockData } };
        (httpClient.post as any).mockResolvedValue(mockResponse);

        const result = await reservationsApi.create(mockData);

        expect(httpClient.post).toHaveBeenCalledWith('/reservations', mockData);
        expect(result).toEqual(mockResponse.data);
    });

    it('getAll calls httpClient.get', async () => {
        const mockResponse = { data: [] };
        (httpClient.get as any).mockResolvedValue(mockResponse);

        const result = await reservationsApi.getAll();

        expect(httpClient.get).toHaveBeenCalledWith('/admin/reservations', { params: {} });
        expect(result).toEqual(mockResponse.data);
    });

    it('getMyReservations calls httpClient.get', async () => {
        const mockResponse = { data: [] };
        (httpClient.get as any).mockResolvedValue(mockResponse);

        const result = await reservationsApi.getMyReservations();

        expect(httpClient.get).toHaveBeenCalledWith('/my/reservations', { params: {} });
        expect(result).toEqual(mockResponse.data);
    });

    it('cancel calls httpClient.patch', async () => {
        const mockResponse = { data: { status: 'cancelled' } };
        (httpClient.patch as any).mockResolvedValue(mockResponse);

        const result = await reservationsApi.cancel(1);

        expect(httpClient.patch).toHaveBeenCalledWith('/reservations/1/cancel');
        expect(result).toEqual(mockResponse.data);
    });

    it('confirm calls httpClient.patch', async () => {
        const mockResponse = { data: { status: 'confirmed' } };
        (httpClient.patch as any).mockResolvedValue(mockResponse);

        const result = await reservationsApi.confirm(1);

        expect(httpClient.patch).toHaveBeenCalledWith('/admin/reservations/1/confirm');
        expect(result).toEqual(mockResponse.data);
    });

    it('adminCancel calls httpClient.patch', async () => {
        const mockResponse = { data: { status: 'cancelled' } };
        (httpClient.patch as any).mockResolvedValue(mockResponse);

        const result = await reservationsApi.adminCancel(1);

        expect(httpClient.patch).toHaveBeenCalledWith('/admin/reservations/1/cancel');
        expect(result).toEqual(mockResponse.data);
    });
});
