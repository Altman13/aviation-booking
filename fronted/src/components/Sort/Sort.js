import React from 'react';
import './Sort.css';

const Sort = ({ sortBy, onSortChange }) => {
  const sortOptions = [
    { value: 'price', label: 'Самый дешёвый' },
    { value: 'duration', label: 'Самый быстрый' },
    { value: 'optimal', label: 'Оптимальный' }
  ];

  const handleSortChange = (value) => {
    onSortChange(value);
  };

  return (
    <div className="sort">
      {sortOptions.map((option) => (
        <button
          key={option.value}
          type="button"
          className={`sort__button ${sortBy === option.value ? 'sort__button--active' : ''}`}
          onClick={() => handleSortChange(option.value)}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
};

export default Sort;