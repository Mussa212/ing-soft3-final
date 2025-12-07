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
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '60px', padding: '40px 20px' }}>
                {/* Left Image - Artie */}
                <div style={{ flex: 1, maxWidth: '900px' }} className="desktop-only">
                    <img
                        src="https://fablehouse.tv/wp-content/uploads/2022/11/Artie-Bucco-The-Sopranos-f02.jpg"
                        alt="Artie Bucco"
                        style={{ width: '100%', borderRadius: '8px', boxShadow: '0 4px 8px rgba(0,0,0,0.2)' }}
                    />
                </div>

                {/* Main Content */}
                <div style={{
                    textAlign: 'center',
                    maxWidth: '600px',
                    flex: 2
                }}>
                    <h1 style={{ fontSize: '48px', marginBottom: '16px', color: 'var(--color-primary)' }}>
                        Vesuvio Ristorante
                    </h1>
                    <p style={{ fontSize: '20px', color: '#555', marginBottom: '40px' }}>
                        Experience authentic Neapolitan cuisine in a warm, family atmosphere.
                        <br />
                        <em style={{ display: 'block', marginTop: '16px', fontSize: '18px' }}>"This shouldn't reach production unless I allow it!" - Artie Bucco</em>
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

                {/* Right Image - Restaurant */}
                <div style={{ flex: 1, maxWidth: '900px' }} className="desktop-only">
                    <img
                        src="https://pbs.twimg.com/media/FAjXILrWQAQV_zU.jpg"
                        alt="Vesuvio Restaurant"
                        style={{ width: '100%', borderRadius: '8px', boxShadow: '0 4px 8px rgba(0,0,0,0.2)' }}
                    />
                </div>
            </div>
        </Layout>
    );
};
