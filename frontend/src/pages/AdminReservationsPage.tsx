import React, { useEffect, useState } from 'react';
import { reservationsApi } from '../api/reservations';
import { Layout } from '../components/Layout';
import { Input } from '../components/Input';
import { ReservationsList } from '../components/ReservationsList';
import type { Reservation } from '../types';
import styles from './Reservations.module.css';

export const AdminReservationsPage: React.FC = () => {
    const [reservations, setReservations] = useState<Reservation[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('all');
    const [filterDate, setFilterDate] = useState(new Date().toISOString().split('T')[0]);

    const fetchReservations = async () => {
        setIsLoading(true);
        try {
            const data = await reservationsApi.getAll(filterDate, filterStatus);
            setReservations(data);
        } catch (err) {
            console.error('Failed to fetch reservations', err);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchReservations();
    }, [filterDate, filterStatus]);

    const handleConfirm = async (id: number) => {
        try {
            await reservationsApi.confirm(id);
            fetchReservations();
        } catch (err) {
            console.error('Failed to confirm reservation', err);
            alert('Failed to confirm reservation');
        }
    };

    const handleCancel = async (id: number) => {
        if (!window.confirm('Are you sure you want to cancel this reservation?')) return;

        try {
            await reservationsApi.adminCancel(id);
            fetchReservations();
        } catch (err) {
            console.error('Failed to cancel reservation', err);
            alert('Failed to cancel reservation');
        }
    };

    return (
        <Layout>
            <div className={styles.container} style={{ maxWidth: '1000px' }}>
                <div className={styles.header}>
                    <h1 className={styles.title}>Admin Panel</h1>
                    <div className={styles.filter}>
                        <Input
                            id="filter-date"
                            type="date"
                            value={filterDate}
                            onChange={(e) => setFilterDate(e.target.value)}
                            style={{ marginBottom: 0 }}
                        />
                        <select
                            value={filterStatus}
                            onChange={(e) => setFilterStatus(e.target.value)}
                            className={styles.select}
                        >
                            <option value="all">All Statuses</option>
                            <option value="pending">Pending</option>
                            <option value="confirmed">Confirmed</option>
                            <option value="cancelled">Cancelled</option>
                        </select>
                    </div>
                </div>

                {isLoading ? (
                    <div className={styles.loading}>Loading reservations...</div>
                ) : (
                    <ReservationsList
                        reservations={reservations}
                        onCancel={handleCancel}
                        onConfirm={handleConfirm}
                        isAdmin={true}
                    />
                )}
            </div>
        </Layout>
    );
};
