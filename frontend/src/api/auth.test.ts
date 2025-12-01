import { describe, it, expect, vi, beforeEach } from 'vitest';
import { authApi } from './auth';
import { httpClient } from './httpClient';

vi.mock('./httpClient', () => ({
    httpClient: {
        post: vi.fn(),
    },
}));

describe('authApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('register calls httpClient.post with correct data', async () => {
        const mockData = { name: 'Test', email: 'test@test.com', password: 'pass' };
        const mockResponse = { data: { id: 1, ...mockData, is_admin: false } };
        (httpClient.post as any).mockResolvedValue(mockResponse);

        const result = await authApi.register(mockData);

        expect(httpClient.post).toHaveBeenCalledWith('/auth/register', mockData);
        expect(result).toEqual(mockResponse.data);
    });

    it('login calls httpClient.post with correct data', async () => {
        const mockData = { email: 'test@test.com', password: 'pass' };
        const mockResponse = { data: { token: 'fake-token' } };
        (httpClient.post as any).mockResolvedValue(mockResponse);

        const result = await authApi.login(mockData);

        expect(httpClient.post).toHaveBeenCalledWith('/auth/login', mockData);
        expect(result).toEqual(mockResponse.data);
    });
});
