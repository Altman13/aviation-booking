// Вспомогательные функции для тестов
export const createMockTicket = (overrides = {}) => ({
  id: 1,
  from: 'MOW',
  to: 'LED',
  departure: '2024-01-15T10:00:00',
  arrival: '2024-01-15T12:00:00',
  duration: 120,
  price: 4500,
  airline: 'S7',
  stops: 0,
  ...overrides
});

// Мок для localStorage
export const mockLocalStorage = () => {
  const localStorageMock = {
    getItem: jest.fn(),
    setItem: jest.fn(),
    removeItem: jest.fn(),
    clear: jest.fn(),
  };
  global.localStorage = localStorageMock;
  return localStorageMock;
};