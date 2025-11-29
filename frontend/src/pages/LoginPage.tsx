import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { authApi } from '../api/auth';
import { Layout } from '../components/Layout';
import { Input } from '../components/Input';
import { Button } from '../components/Button';
import styles from './Auth.module.css';

export const LoginPage: React.FC = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            const user = await authApi.login({ email, password });
            login(user);
            navigate(user.is_admin ? '/admin/reservations' : '/reservations');
        } catch (err) {
            setError('Invalid credentials. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Layout>
            <div className={styles.authContainer} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '40px' }}>
                {/* Left Image - Artie */}
                <div style={{ flex: 1, maxWidth: '900px', minWidth: '0' }} className="desktop-only">
                    <img
                        src="https://www.spoton.com/blog/content/images/2024/01/restaurant-tech-for-the-sopranos.png"
                        alt="Artie Bucco"
                        style={{ width: '100%', borderRadius: '8px', boxShadow: '0 4px 8px rgba(0,0,0,0.2)' }}
                    />
                </div>

                <div className={styles.card} style={{ flexShrink: 0, width: '400px' }}>
                    <h1 className={styles.title}>Welcome Back</h1>
                    <p className={styles.subtitle}>Sign in to manage your reservations</p>

                    <form onSubmit={handleSubmit}>
                        <Input
                            id="email"
                            name="email"
                            label="Email"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            fullWidth
                        />
                        <Input
                            id="password"
                            name="password"
                            label="Password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            fullWidth
                        />

                        {error && <div className={styles.errorMessage}>{error}</div>}

                        <Button type="submit" fullWidth disabled={isLoading} style={{ marginTop: '16px' }}>
                            {isLoading ? 'Signing in...' : 'Sign In'}
                        </Button>
                    </form>

                    <div className={styles.footer}>
                        Don't have an account? <Link to="/register">Create one</Link>
                    </div>
                </div>

                {/* Right Image - Restaurant */}
                <div style={{ flex: 1, maxWidth: '900px', minWidth: '0' }} className="desktop-only">
                    <img
                        src="https://www.spoton.com/blog/content/images/2024/01/Nuovo-Vesuvio-1024x576-1.jpg"
                        alt="Vesuvio Restaurant"
                        style={{ width: '100%', borderRadius: '8px', boxShadow: '0 4px 8px rgba(0,0,0,0.2)' }}
                    />
                </div>
            </div>
        </Layout>
    );
};
