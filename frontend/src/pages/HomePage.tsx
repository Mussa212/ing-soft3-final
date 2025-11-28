import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Layout } from '../components/Layout';
import { Button } from '../components/Button';
import { Link } from 'react-router-dom';

export const HomePage: React.FC = () => {
    const { currentUser, isLoading } = useAuth();

    if (isLoading) return <div>Loading...</div>;

    if (currentUser) {
        return <Navigate to={currentUser.is_admin ? "/admin/reservations" : "/reservations"} replace />;
    }

    return (
        <Layout>
            <div style={{
                textAlign: 'center',
                marginTop: '60px',
                padding: '0 20px'
            }}>
                <h1 style={{ fontSize: '48px', marginBottom: '16px', color: 'var(--color-primary)' }}>
                    Vesuvio Ristorante
                </h1>
                <p style={{ fontSize: '20px', color: '#555', marginBottom: '40px', maxWidth: '600px', margin: '0 auto 40px' }}>
                    Experience authentic Neapolitan cuisine in a warm, family atmosphere.
                    "We lead the world in computerized data collection!" - Artie Bucco
                </p>

                <div style={{ display: 'flex', gap: '16px', justifyContent: 'center' }}>
                    <Link to="/login">
                        <Button style={{ minWidth: '150px' }}>Login</Button>
                    </Link>
                    <Link to="/register">
                        <Button variant="secondary" style={{ minWidth: '150px' }}>Register</Button>
                    </Link>
                </div>
            </div>
        </Layout>
    );
};
