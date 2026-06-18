import React, { useState, useEffect } from 'react';
import './App.css';
import SearchForm from './components/SearchForm/SearchForm';
import TicketList from './components/TicketList/TicketList';
import Filters from './components/Filters/Filters';
import Sort from './components/Sort/Sort';

function App() {
  const [tickets, setTickets] = useState([]);
  const [filteredTickets, setFilteredTickets] = useState([]);
  const [filters, setFilters] = useState({
    stops: [],
    price: { min: 0, max: 100000 },
    airlines: []
  });
  const [sortBy, setSortBy] = useState('price');

  // Моковые данные билетов
  useEffect(() => {
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
      },
      {
        id: 3,
        from: 'MOW',
        to: 'LED',
        departure: '2024-01-15T18:00:00',
        arrival: '2024-01-15T19:30:00',
        duration: 90,
        price: 5200,
        airline: 'Pobeda',
        stops: 0
      },
    ];
    setTickets(mockTickets);
    setFilteredTickets(mockTickets);
  }, []);

  return (
    <div className="app">
      <header className="app-header">
        <div className="container">
          <div className="logo">Astalavista</div>
        </div>
      </header>
      
      <main className="app-main">
        <div className="container">
          <div className="app-layout">
            {/* Боковая панель с фильтрами */}
            <aside className="sidebar">
              <Filters 
                filters={filters}
                onFiltersChange={setFilters}
              />
            </aside>
            
            {/* Основной контент */}
            <div className="content">
              <SearchForm />
              <Sort 
                sortBy={sortBy}
                onSortChange={setSortBy}
              />
              <TicketList tickets={filteredTickets} />
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;