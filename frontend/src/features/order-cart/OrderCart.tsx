import React from 'react';
import type { Order } from '../../entities/order';
import styles from './OrderCart.module.css'

interface OrderCardProps {
  order: Order;
}

export const OrderCart: React.FC<OrderCardProps> = ({ order }) => {
  return (
    <div className={styles.orderCart}>
      <h2>Order Details</h2>
      <p><strong>Order UID:</strong> {order.order_uid}</p>
      <p><strong>Track Number:</strong> {order.track_number}</p>
      <p><strong>Customer ID:</strong> {order.customer_id}</p>
    </div>
  );
};
