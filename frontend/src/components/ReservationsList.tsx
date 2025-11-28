import React from 'react';
import { Button } from './Button';
import type { Reservation } from '../types';
import styles from '../pages/Reservations.module.css';
import classNames from 'classnames';

interface ReservationsListProps {
    reservations: Reservation[];
    onCancel: (id: number) => void;
    isAdmin?: boolean;
    onConfirm?: (id: number) => void;
}

export const ReservationsList: React.FC<ReservationsListProps> = ({
    reservations,
    onCancel,
    isAdmin = false,
    onConfirm
}) => {
    if (reservations.length === 0) {
        return <div className={styles.empty}>No reservations found.</div>;
    }

    return (
        <div className={styles.list}>
            {reservations.map((reservation) => (
                <div key={reservation.id} className={styles.reservationCard} role="article">
                    <div className={styles.reservationInfo}>
                        <div className={styles.date}>
                            {isAdmin ? (
                                <>
                                    {reservation.time} - {reservation.people} people
                                </>
                            ) : (
                                <>
                                    {new Date(reservation.date).toLocaleDateString()} at {reservation.time}
                                </>
                            )}
                        </div>

                        <div className={styles.details}>
                            {isAdmin ? (
                                <>
                                    <strong>Client:</strong> {reservation.user?.name} ({reservation.user?.email})
                                </>
                            ) : (
                                <>
                                    {reservation.people} people
                                </>
                            )}
                            {reservation.comment && <span className={styles.comment}> • "{reservation.comment}"</span>}
                        </div>

                        <div className={classNames(styles.status, styles[reservation.status])}>
                            {reservation.status.toUpperCase()}
                        </div>
                    </div>

                    <div className={styles.actions}>
                        {isAdmin && reservation.status === 'pending' && onConfirm && (
                            <Button
                                onClick={() => onConfirm(reservation.id)}
                                style={{ marginRight: '8px', backgroundColor: 'var(--color-accent)', color: 'white' }}
                            >
                                Confirm
                            </Button>
                        )}

                        {(reservation.status === 'pending' || reservation.status === 'confirmed') && (
                            <Button
                                variant="danger"
                                onClick={() => onCancel(reservation.id)}
                                className={styles.cancelButton}
                            >
                                Cancel
                            </Button>
                        )}
                    </div>
                </div>
            ))}
        </div>
    );
};
