import React, { type InputHTMLAttributes } from 'react';
import classNames from 'classnames';
import styles from './Input.module.css';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    error?: string;
    fullWidth?: boolean;
}

export const Input: React.FC<InputProps> = ({
    label,
    error,
    fullWidth = false,
    className,
    id,
    ...props
}) => {
    const inputId = id || props.name;

    return (
        <div className={classNames(styles.container, { [styles.fullWidth]: fullWidth, [styles.hasError]: !!error })}>
            {label && <label htmlFor={inputId} className={styles.label}>{label}</label>}
            <input
                id={inputId}
                className={classNames(styles.input, className)}
                {...props}
            />
            {error && <span className={styles.errorText}>{error}</span>}
        </div>
    );
};
