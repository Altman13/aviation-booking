import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';

describe('App Integration Tests', () => {
  test('renders main application with tickets', async () => {
    render(<App />);
    
    expect(screen.getByText(/Astalavista/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /найти билеты/i })).toBeInTheDocument();
    
    await waitFor(() => {
      expect(screen.getAllByRole('button', { name: /выбрать/i })).toHaveLength(4);
    });
  });

  test('filters tickets by stops', async () => {
    const user = userEvent.setup();
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getAllByRole('button', { name: /выбрать/i })).toHaveLength(4);
    });

    const noStopsCheckbox = screen.getByLabelText(/без пересадок/i);
    await user.click(noStopsCheckbox);

    await waitFor(() => {
      const selectButtons = screen.getAllByRole('button', { name: /выбрать/i });
      expect(selectButtons.length).toBeLessThan(4);
    });
  });

  test('changes sort order', async () => {
    const user = userEvent.setup();
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getAllByRole('button', { name: /выбрать/i })).toHaveLength(4);
    });

    const fastestButton = screen.getByRole('button', { name: /самый быстрый/i });
    await user.click(fastestButton);

    expect(fastestButton).toHaveClass('sort__button--active');
  });
});