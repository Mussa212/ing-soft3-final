import React, { type ButtonHTMLAttributes } from 'react';
import classNames from 'classnames';
import styles from './Button.module.css';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'secondary' | 'outline' | 'danger';
    fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
    children,
    variant = 'primary',
    fullWidth = false,
    className,
    ...props
}) => {
    return (
        <button
            className={classNames(
                styles.button,
                styles[variant],
                { [styles.fullWidth]: fullWidth },
                className
            )}
            {...props}
        >
            {children}
        </button>
    );
};
