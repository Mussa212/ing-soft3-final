import React, { useEffect, useState } from 'react';
import { reservationsApi } from '../api/reservations';
import { Layout } from '../components/Layout';
import { ReservationsList } from '../components/ReservationsList';
import type { Reservation } from '../types';
import styles from './Reservations.module.css';

export const MyReservationsPage: React.FC = () => {
    const [reservations, setReservations] = useState<Reservation[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filter, setFilter] = useState('all');

    const fetchReservations = async () => {
        setIsLoading(true);
        try {
            const data = await reservationsApi.getMyReservations(filter);
            setReservations(data);
        } catch (err) {
            console.error('Failed to fetch reservations', err);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchReservations();
    }, [filter]);

    const handleCancel = async (id: number) => {
        if (!window.confirm('Are you sure you want to cancel this reservation?')) return;

        try {
            await reservationsApi.cancel(id);
            fetchReservations();
        } catch (err) {
            console.error('Failed to cancel reservation', err);
            alert('Failed to cancel reservation');
        }
    };

    return (
        <Layout>
            <div className={styles.container}>
                <div className={styles.header}>
                    <h1 className={styles.title}>My Reservations</h1>
                    <div className={styles.filter}>
                        <label htmlFor="status-filter">Status:</label>
                        <select
                            id="status-filter"
                            value={filter}
                            onChange={(e) => setFilter(e.target.value)}
                            className={styles.select}
                        >
                            <option value="all">All</option>
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
                    />
                )}
            </div>
        </Layout>
    );
};
