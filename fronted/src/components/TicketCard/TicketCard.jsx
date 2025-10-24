// components/TicketCard/TicketCard.js
import React from 'react';
import './TicketCard.css';

const TicketCard = ({ ticket }) => {
  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString('ru-RU', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatDuration = (minutes) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}ч ${mins}м`;
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU').format(price) + ' ₽';
  };

  const getStopsText = (stops) => {
    if (stops === 0) return 'Без пересадок';
    if (stops === 1) return '1 пересадка';
    return `${stops} пересадки`;
  };

  return (
    <div className="ticket-card">
      <div className="ticket-card__header">
        <div className="ticket-card__price">{formatPrice(ticket.price)}</div>
        <div className="ticket-card__airline">
          <img 
            src={`/airlines/${ticket.airline}.png`} 
            alt={ticket.airline}
            className="ticket-card__airline-logo"
          />
          <span>{ticket.airline}</span>
        </div>
      </div>
      
      <div className="ticket-card__body">
        <div className="ticket-card__segment">
          <div className="ticket-card__route">
            <div className="ticket-card__cities">
              <span className="ticket-card__city">{ticket.from} – {ticket.to}</span>
            </div>
            <div className="ticket-card__times">
              <span className="ticket-card__time">{formatTime(ticket.departure)}</span>
              <span className="ticket-card__time-separator">–</span>
              <span className="ticket-card__time">{formatTime(ticket.arrival)}</span>
            </div>
            <div className="ticket-card__duration">
              В пути: {formatDuration(ticket.duration)}
            </div>
            <div className="ticket-card__stops">
              {getStopsText(ticket.stops)}
            </div>
          </div>
        </div>
      </div>
      
      <button className="ticket-card__button">
        Выбрать
      </button>
    </div>
  );
};

export default TicketCard;