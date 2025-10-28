import React from 'react';
import { render, screen } from '@testing-library/react';
import TicketCard from './TicketCard';

const mockTicket = {
  id: 1,
  from: 'MOW',
  to: 'LED',
  departure: '2024-01-15T10:00:00',
  arrival: '2024-01-15T12:00:00',
  duration: 120,
  price: 4500,
  airline: 'S7',
  stops: 0
};

describe('TicketCard Component', () => {
  test('renders ticket information correctly', () => {
    render(<TicketCard ticket={mockTicket} />);
    
    expect(screen.getByText('4 500 ₽')).toBeInTheDocument();
    expect(screen.getByText('S7')).toBeInTheDocument();
    expect(screen.getByText('MOW – LED')).toBeInTheDocument();
    expect(screen.getByText('10:00')).toBeInTheDocument();
    expect(screen.getByText('12:00')).toBeInTheDocument();
    expect(screen.getByText(/в пути: 2ч 0м/i)).toBeInTheDocument();
    expect(screen.getByText(/без пересадок/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /выбрать/i })).toBeInTheDocument();
  });

  test('displays correct stops text for different stop counts', () => {
    const ticketWithOneStop = { ...mockTicket, stops: 1 };
    const ticketWithTwoStops = { ...mockTicket, stops: 2 };
    
    const { rerender } = render(<TicketCard ticket={ticketWithOneStop} />);
    expect(screen.getByText(/1 пересадка/i)).toBeInTheDocument();
    
    rerender(<TicketCard ticket={ticketWithTwoStops} />);
    expect(screen.getByText(/2 пересадки/i)).toBeInTheDocument();
  });

  test('formats price correctly with thousands separator', () => {
    const expensiveTicket = { ...mockTicket, price: 15000 };
    render(<TicketCard ticket={expensiveTicket} />);
    
    expect(screen.getByText('15 000 ₽')).toBeInTheDocument();
  });
});