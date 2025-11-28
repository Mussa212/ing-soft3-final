import React, { type ReactNode } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Button } from './Button';
import styles from './Layout.module.css';

interface LayoutProps {
    children: ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
    const { currentUser, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <div className={styles.layout}>
            <header className={styles.header}>
                <div className={styles.headerContent}>
                    <Link to="/" className={styles.logo}>
                        Vesuvio Ristorante
                    </Link>
                    <nav className={styles.nav}>
                        {currentUser ? (
                            <>
                                {currentUser.is_admin ? (
                                    <Link to="/admin/reservations" className={styles.navLink}>Admin Panel</Link>
                                ) : (
                                    <>
                                        <Link to="/reservations" className={styles.navLink}>My Reservations</Link>
                                        <Link to="/reservations/new" className={styles.navLink}>New Reservation</Link>
                                    </>
                                )}
                                <div className={styles.userMenu}>
                                    <span className={styles.userName}>{currentUser.name}</span>
                                    <Button variant="outline" onClick={handleLogout} style={{ padding: '6px 12px', fontSize: '14px' }}>
                                        Logout
                                    </Button>
                                </div>
                            </>
                        ) : (
                            <>
                                <Link to="/login" className={styles.navLink}>Login</Link>
                                <Link to="/register" className={styles.navLink}>Register</Link>
                            </>
                        )}
                    </nav>
                </div>
            </header>
            <main className={styles.main}>
                {children}
            </main>
            <footer className={styles.footer}>
                <p>&copy; {new Date().getFullYear()} Vesuvio Ristorante. All rights reserved.</p>
            </footer>
        </div>
    );
};
