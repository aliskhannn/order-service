// pages/SearchPage.tsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { SearchBar } from '../../features/search/SearchBar';
import type { Order } from '../../entities/order';
import { OrderCart } from '../../features/order-cart';
import styles from './SearchPage.module.css';

export const SearchPage: React.FC = () => {
  const [order, setOrder] = useState<Order | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();

  const handleSearch = async (query: string) => {
    setLoading(true);
    setError(null);
    setOrder(null);

    try {
      const response = await fetch(`http://localhost:8080/orders/${query}`);
      if (!response.ok) {
        throw new Error(`Order not found: ${response.statusText}`);
      }
      const data: Order = await response.json();
      setOrder(data);
    } catch (err: any) {
      setError(err.message || 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const handleCardClick = () => {
    if (order) {
      navigate(`/order/${order.order_uid}`);
    }
  };

  return (
    <div className={styles.searchPage}>
      <h1>Order Search</h1>
      <SearchBar onSearch={handleSearch} />

      {loading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>Error: {error}</p>}

      {order && (
        <div
          onClick={handleCardClick}
          style={{ minWidth: '350px', marginTop: '50px' }}
        >
          <OrderCart order={order} />
        </div>
      )}
    </div>
  );
};
