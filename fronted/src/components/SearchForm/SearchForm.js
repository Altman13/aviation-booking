// components/SearchForm/SearchForm.js
import React, { useState } from 'react';
import './SearchForm.css';

const SearchForm = () => {
  const [searchParams, setSearchParams] = useState({
    from: '',
    to: '',
    departureDate: '',
    returnDate: '',
    passengers: 1
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('Search params:', searchParams);
    // Здесь будет логика поиска
  };

  const handleInputChange = (field, value) => {
    setSearchParams(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const swapCities = () => {
    setSearchParams(prev => ({
      ...prev,
      from: prev.to,
      to: prev.from
    }));
  };

  return (
    <form className="search-form" onSubmit={handleSubmit}>
      <div className="search-form__row">
        <div className="search-form__field">
          <label className="search-form__label">Откуда</label>
          <input
            type="text"
            className="search-form__input"
            placeholder="Город вылета"
            value={searchParams.from}
            onChange={(e) => handleInputChange('from', e.target.value)}
          />
        </div>
        
        <button 
          type="button" 
          className="search-form__swap-btn"
          onClick={swapCities}
        >
          ↔
        </button>
        
        <div className="search-form__field">
          <label className="search-form__label">Куда</label>
          <input
            type="text"
            className="search-form__input"
            placeholder="Город прилёта"
            value={searchParams.to}
            onChange={(e) => handleInputChange('to', e.target.value)}
          />
        </div>
        
        <div className="search-form__field">
          <label className="search-form__label">Туда</label>
          <input
            type="date"
            className="search-form__input"
            value={searchParams.departureDate}
            onChange={(e) => handleInputChange('departureDate', e.target.value)}
          />
        </div>
        
        <div className="search-form__field">
          <label className="search-form__label">Обратно</label>
          <input
            type="date"
            className="search-form__input"
            value={searchParams.returnDate}
            onChange={(e) => handleInputChange('returnDate', e.target.value)}
          />
        </div>
        
        <div className="search-form__field">
          <label className="search-form__label">Пассажиры</label>
          <select
            className="search-form__select"
            value={searchParams.passengers}
            onChange={(e) => handleInputChange('passengers', parseInt(e.target.value))}
          >
            {[1,2,3,4,5,6].map(num => (
              <option key={num} value={num}>{num} пассажир{num > 1 ? 'а' : ''}</option>
            ))}
          </select>
        </div>
      </div>
      
      <button type="submit" className="search-form__submit">
        Найти билеты
      </button>
    </form>
  );
};

export default SearchForm;