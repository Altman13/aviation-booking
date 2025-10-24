// components/Filters/Filters.js
import React from 'react';
import './Filters.css';

const Filters = ({ filters, onFiltersChange }) => {
  const handleStopsChange = (stopCount) => {
    const newStops = filters.stops.includes(stopCount)
      ? filters.stops.filter(s => s !== stopCount)
      : [...filters.stops, stopCount];
    
    onFiltersChange({
      ...filters,
      stops: newStops
    });
  };

  return (
    <div className="filters">
      <div className="filters__section">
        <h3 className="filters__title">Количество пересадок</h3>
        <div className="filters__options">
          {[0, 1, 2, 3].map(stopCount => (
            <label key={stopCount} className="filters__option">
              <input
                type="checkbox"
                checked={filters.stops.includes(stopCount)}
                onChange={() => handleStopsChange(stopCount)}
                className="filters__checkbox"
              />
              <span className="filters__label">
                {stopCount === 0 ? 'Без пересадок' : 
                 stopCount === 1 ? '1 пересадка' : 
                 `${stopCount} пересадки`}
              </span>
            </label>
          ))}
        </div>
      </div>
      
      <div className="filters__section">
        <h3 className="filters__title">Авиакомпании</h3>
        <div className="filters__options">
          {['S7', 'Aeroflot', 'Pobeda', 'Utair'].map(airline => (
            <label key={airline} className="filters__option">
              <input
                type="checkbox"
                className="filters__checkbox"
              />
              <span className="filters__label">{airline}</span>
            </label>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Filters;