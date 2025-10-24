import React from 'react';
import TicketCard from '../TicketCard/TicketCard';
import './TicketList.css';

const TicketList = ({ tickets }) => {
  if (!tickets || tickets.length === 0) {
    return (
      <div className="ticket-list ticket-list--empty">
        <p>Билеты не найдены</p>
      </div>
    );
  }

  return (
    <div className="ticket-list">
      {tickets.map(ticket => (
        <TicketCard key={ticket.id} ticket={ticket} />
      ))}
    </div>
  );
};

export default TicketList;