import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Filters from './Filters';

describe('Filters Component', () => {
  const mockFilters = {
    stops: [],
    price: { min: 0, max: 100000 },
    airlines: []
  };

  test('renders all filter sections', () => {
    render(<Filters filters={mockFilters} onFiltersChange={jest.fn()} />);
    
    expect(screen.getByText(/количество пересадок/i)).toBeInTheDocument();
    expect(screen.getByText(/авиакомпании/i)).toBeInTheDocument();
    
    expect(screen.getByLabelText(/без пересадок/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/1 пересадка/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/2 пересадки/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/3 пересадки/i)).toBeInTheDocument();
  });

  test('calls onFiltersChange when stop filter is clicked', async () => {
    const mockOnFiltersChange = jest.fn();
    const user = userEvent.setup();
    
    render(<Filters filters={mockFilters} onFiltersChange={mockOnFiltersChange} />);
    
    const noStopsCheckbox = screen.getByLabelText(/без пересадок/i);
    await user.click(noStopsCheckbox);

    expect(mockOnFiltersChange).toHaveBeenCalledWith({
      ...mockFilters,
      stops: [0]
    });
  });

  test('handles multiple stop selections', async () => {
    const mockOnFiltersChange = jest.fn();
    const user = userEvent.setup();
    
    render(<Filters filters={mockFilters} onFiltersChange={mockOnFiltersChange} />);
    
    const noStopsCheckbox = screen.getByLabelText(/без пересадок/i);
    const oneStopCheckbox = screen.getByLabelText(/1 пересадка/i);
    
    await user.click(noStopsCheckbox);
    await user.click(oneStopCheckbox);

    expect(mockOnFiltersChange).toHaveBeenLastCalledWith({
      ...mockFilters,
      stops: [0, 1]
    });
  });

  test('unchecks stop filter when clicked again', async () => {
    const mockOnFiltersChange = jest.fn();
    const user = userEvent.setup();
    const filtersWithStop = { ...mockFilters, stops: [0] };
    
    render(<Filters filters={filtersWithStop} onFiltersChange={mockOnFiltersChange} />);
    
    const noStopsCheckbox = screen.getByLabelText(/без пересадок/i);
    await user.click(noStopsCheckbox);

    expect(mockOnFiltersChange).toHaveBeenCalledWith({
      ...mockFilters,
      stops: []
    });
  });
});