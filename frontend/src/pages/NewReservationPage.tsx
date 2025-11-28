import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { reservationsApi } from '../api/reservations';
import { Layout } from '../components/Layout';
import { Input } from '../components/Input';
import { Button } from '../components/Button';
import styles from './Reservations.module.css';

export const NewReservationPage: React.FC = () => {
    const [date, setDate] = useState('');
    const [time, setTime] = useState('');
    const [people, setPeople] = useState(2);
    const [comment, setComment] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        console.log('handleSubmit called', { date, time, people, comment });
        setError('');

        if (people <= 0) {
            setError('Number of people must be greater than 0');
            return;
        }

        setIsLoading(true);

        try {
            await reservationsApi.create({ date, time, people, comment });
            navigate('/reservations');
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to create reservation. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    // Get today's date for min attribute
    const today = new Date().toISOString().split('T')[0];

    return (
        <Layout>
            <div className={styles.container}>
                <div className={styles.card}>
                    <h1 className={styles.title}>New Reservation</h1>
                    <p className={styles.subtitle}>Book a table at Vesuvio</p>

                    <form onSubmit={handleSubmit} noValidate>
                        <div className={styles.row}>
                            <Input
                                id="date"
                                label="Date"
                                type="date"
                                value={date}
                                min={today}
                                onChange={(e) => setDate(e.target.value)}
                                required
                                fullWidth
                            />
                            <Input
                                id="time"
                                label="Time"
                                type="time"
                                value={time}
                                onChange={(e) => setTime(e.target.value)}
                                required
                                fullWidth
                            />
                        </div>

                        <Input
                            id="people"
                            label="Number of People"
                            type="number"
                            min="1"
                            value={people}
                            onChange={(e) => setPeople(parseInt(e.target.value))}
                            required
                            fullWidth
                        />

                        <div className={styles.formGroup}>
                            <label htmlFor="comment" className={styles.label}>Special Requests (Optional)</label>
                            <textarea
                                id="comment"
                                className={styles.textarea}
                                value={comment}
                                onChange={(e) => setComment(e.target.value)}
                                rows={3}
                            />
                        </div>

                        {error && <div className={styles.errorMessage}>{error}</div>}

                        <Button type="submit" fullWidth disabled={isLoading} style={{ marginTop: '24px' }}>
                            {isLoading ? 'Booking...' : 'Confirm Reservation'}
                        </Button>
                    </form>
                </div>
            </div>
        </Layout>
    );
};
