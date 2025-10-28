import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Sort from './Sort';

describe('Sort Component', () => {
  test('renders all sort options', () => {
    render(<Sort sortBy="price" onSortChange={jest.fn()} />);
    
    expect(screen.getByRole('button', { name: /самый дешёвый/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /самый быстрый/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /оптимальный/i })).toBeInTheDocument();
  });

  test('highlights active sort option', () => {
    render(<Sort sortBy="duration" onSortChange={jest.fn()} />);
    
    const cheapestButton = screen.getByRole('button', { name: /самый дешёвый/i });
    const fastestButton = screen.getByRole('button', { name: /самый быстрый/i });
    const optimalButton = screen.getByRole('button', { name: /оптимальный/i });

    expect(cheapestButton).not.toHaveClass('sort__button--active');
    expect(fastestButton).toHaveClass('sort__button--active');
    expect(optimalButton).not.toHaveClass('sort__button--active');
  });

  test('calls onSortChange with correct value when button is clicked', async () => {
    const mockOnSortChange = jest.fn();
    const user = userEvent.setup();
    
    render(<Sort sortBy="price" onSortChange={mockOnSortChange} />);
    
    const optimalButton = screen.getByRole('button', { name: /оптимальный/i });
    await user.click(optimalButton);

    expect(mockOnSortChange).toHaveBeenCalledWith('optimal');
  });
});