import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export const AdminRoute: React.FC = () => {
    const { currentUser, isLoading } = useAuth();

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (!currentUser) {
        return <Navigate to="/login" replace />;
    }

    if (!currentUser.is_admin) {
        return <Navigate to="/reservations" replace />;
    }

    return <Outlet />;
};
