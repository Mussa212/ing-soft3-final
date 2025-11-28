import { httpClient } from './httpClient';
import type { LoginResponse, User } from '../types';

export const authApi = {
    register: async (data: Omit<User, 'id' | 'is_admin'> & { password: string }) => {
        const response = await httpClient.post('/auth/register', data);
        return response.data;
    },

    login: async (data: Pick<User, 'email'> & { password: string }) => {
        const response = await httpClient.post<LoginResponse>('/auth/login', data);
        return response.data;
    },
};
