import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SearchForm from './SearchForm';

describe('SearchForm Component', () => {
  test('renders search form with all fields', () => {
    render(<SearchForm />);
    
    expect(screen.getByLabelText(/откуда/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/куда/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/туда/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/обратно/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/пассажиры/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /найти билеты/i })).toBeInTheDocument();
  });

  test('allows user to fill form fields', async () => {
    const user = userEvent.setup();
    render(<SearchForm />);
    
    const fromInput = screen.getByPlaceholderText(/город вылета/i);
    const toInput = screen.getByPlaceholderText(/город прилёта/i);
    const passengersSelect = screen.getByLabelText(/пассажиры/i);

    await user.type(fromInput, 'Moscow');
    await user.type(toInput, 'London');
    await user.selectOptions(passengersSelect, '2');

    expect(fromInput).toHaveValue('Moscow');
    expect(toInput).toHaveValue('London');
    expect(passengersSelect).toHaveValue('2');
  });

  test('submits form with correct data', async () => {
    const user = userEvent.setup();
    const consoleSpy = jest.spyOn(console, 'log');
    
    render(<SearchForm />);
    
    const fromInput = screen.getByPlaceholderText(/город вылета/i);
    const submitButton = screen.getByRole('button', { name: /найти билеты/i });

    await user.type(fromInput, 'Moscow');
    await user.click(submitButton);

    expect(consoleSpy).toHaveBeenCalledWith('Search params:', expect.any(Object));
    
    consoleSpy.mockRestore();
  });

  test('swaps cities when swap button is clicked', async () => {
    const user = userEvent.setup();
    render(<SearchForm />);
    
    const fromInput = screen.getByPlaceholderText(/город вылета/i);
    const toInput = screen.getByPlaceholderText(/город прилёта/i);
    const swapButton = screen.getByRole('button', { name: '↔' });

    await user.type(fromInput, 'Moscow');
    await user.type(toInput, 'London');
    await user.click(swapButton);

    expect(fromInput).toHaveValue('London');
    expect(toInput).toHaveValue('Moscow');
  });
});