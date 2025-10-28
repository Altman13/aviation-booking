import React from 'react';
import { render, screen } from '@testing-library/react';
import TicketList from './TicketList';

const mockTickets = [
  {
    id: 1,
    from: 'MOW',
    to: 'LED',
    departure: '2024-01-15T10:00:00',
    arrival: '2024-01-15T12:00:00',
    duration: 120,
    price: 4500,
    airline: 'S7',
    stops: 0
  },
  {
    id: 2,
    from: 'MOW',
    to: 'LED',
    departure: '2024-01-15T14:00:00',
    arrival: '2024-01-15T17:30:00',
    duration: 210,
    price: 3200,
    airline: 'Aeroflot',
    stops: 1
  }
];

describe('TicketList Component', () => {
  test('renders list of tickets', () => {
    render(<TicketList tickets={mockTickets} />);
    
    expect(screen.getAllByRole('button', { name: /выбрать/i })).toHaveLength(2);
    expect(screen.getByText('S7')).toBeInTheDocument();
    expect(screen.getByText('Aeroflot')).toBeInTheDocument();
  });

  test('renders empty state when no tickets', () => {
    render(<TicketList tickets={[]} />);
    
    expect(screen.getByText(/билеты не найдены/i)).toBeInTheDocument();
  });

  test('renders empty state when tickets is null', () => {
    render(<TicketList tickets={null} />);
    
    expect(screen.getByText(/билеты не найдены/i)).toBeInTheDocument();
  });
});