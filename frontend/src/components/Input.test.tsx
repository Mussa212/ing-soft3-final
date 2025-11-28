import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Input } from './Input';

describe('Input', () => {
    it('renders label and input', () => {
        render(<Input label="Test Label" id="test-input" />);
        expect(screen.getByLabelText(/test label/i)).toBeInTheDocument();
        expect(screen.getByRole('textbox')).toBeInTheDocument();
    });

    it('handles change events', () => {
        const handleChange = vi.fn();
        render(<Input onChange={handleChange} />);
        const input = screen.getByRole('textbox');
        fireEvent.change(input, { target: { value: 'test' } });
        expect(handleChange).toHaveBeenCalledTimes(1);
    });

    it('displays error message', () => {
        render(<Input error="Invalid input" />);
        expect(screen.getByText(/invalid input/i)).toBeInTheDocument();
    });

    it('applies required attribute', () => {
        render(<Input required />);
        expect(screen.getByRole('textbox')).toBeRequired();
    });

    it('uses id for label association', () => {
        render(<Input label="Label" id="my-id" />);
        const input = screen.getByLabelText("Label");
        expect(input).toHaveAttribute('id', 'my-id');
    });
});
