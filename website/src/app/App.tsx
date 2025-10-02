// App.tsx
import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { SearchPage } from '../pages/search/SearchPage';
import { OrderPage } from '../pages/order/OrderPage';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<SearchPage />} />
        <Route path="/order/:orderId" element={<OrderPage />} />
      </Routes>
    </BrowserRouter>
  );
};
