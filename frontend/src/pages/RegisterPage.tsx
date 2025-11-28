import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../api/auth';
import { Layout } from '../components/Layout';
import { Input } from '../components/Input';
import { Button } from '../components/Button';
import styles from './Auth.module.css';

export const RegisterPage: React.FC = () => {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            await authApi.register({ name, email, password });
            navigate('/login');
        } catch (err: any) {
            setError(err.response?.data?.message || 'Registration failed. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Layout>
            <div className={styles.authContainer}>
                <div className={styles.card}>
                    <h1 className={styles.title}>Create Account</h1>
                    <p className={styles.subtitle}>Join Vesuvio Ristorante</p>

                    <form onSubmit={handleSubmit}>
                        <Input
                            id="name"
                            name="name"
                            label="Full Name"
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            required
                            fullWidth
                        />
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
                            {isLoading ? 'Creating Account...' : 'Register'}
                        </Button>
                    </form>

                    <div className={styles.footer}>
                        Already have an account? <Link to="/login">Sign in</Link>
                    </div>
                </div>
            </div>
        </Layout>
    );
};
