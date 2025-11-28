import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { PrivateRoute } from './PrivateRoute';
import { AdminRoute } from './AdminRoute';
import { HomePage } from '../pages/HomePage';
import { LoginPage } from '../pages/LoginPage';
import { RegisterPage } from '../pages/RegisterPage';
import { NewReservationPage } from '../pages/NewReservationPage';
import { MyReservationsPage } from '../pages/MyReservationsPage';
import { AdminReservationsPage } from '../pages/AdminReservationsPage';

export const AppRoutes: React.FC = () => {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />

            <Route element={<PrivateRoute />}>
                <Route path="/reservations" element={<MyReservationsPage />} />
                <Route path="/reservations/new" element={<NewReservationPage />} />
            </Route>

            <Route element={<AdminRoute />}>
                <Route path="/admin/reservations" element={<AdminReservationsPage />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
    );
};
