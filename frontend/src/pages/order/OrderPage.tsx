// pages/OrderDetailsPage.tsx
import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Order } from '../../entities/order';

export const OrderPage: React.FC = () => {
  const { orderId } = useParams<{ orderId: string }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!orderId) return;
    setLoading(true);
    setError(null);

    fetch(`http://localhost:8080/orders/${orderId}`)
      .then((res) => {
        if (!res.ok) throw new Error('Order not found');
        return res.json();
      })
      .then((data: Order) => setOrder(data))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [orderId]);

  if (loading) return <p>Loading order details...</p>;
  if (error) return <p style={{ color: 'red' }}>Error: {error}</p>;
  if (!order) return <p>No order data</p>;

  return (
<div className='order-page'>
      <h1>Order Details</h1>
      <Link to="/">← Back to Search</Link>

      <p><strong>Order UID:</strong> {order.order_uid}</p>
      <p><strong>Track Number:</strong> {order.track_number}</p>
      <p><strong>Entry:</strong> {order.entry}</p>
      <p><strong>Locale:</strong> {order.locale}</p>
      <p><strong>Internal Signature:</strong> {order.internal_signature || '-'}</p>
      <p><strong>Customer ID:</strong> {order.customer_id}</p>
      <p><strong>Delivery Service:</strong> {order.delivery_service}</p>
      <p><strong>Shardkey:</strong> {order.shardkey}</p>
      <p><strong>SM ID:</strong> {order.sm_id}</p>
      <p><strong>Date Created:</strong> {new Date(order.date_created).toLocaleString()}</p>
      <p><strong>OOF Shard:</strong> {order.oof_shard}</p>

      <h2>Delivery Info</h2>
      <p>Name: {order.delivery.name}</p>
      <p>Phone: {order.delivery.phone}</p>
      <p>Zip: {order.delivery.zip}</p>
      <p>City: {order.delivery.city}</p>
      <p>Address: {order.delivery.address}</p>
      <p>Region: {order.delivery.region}</p>
      <p>Email: {order.delivery.email}</p>

      <h2>Payment Info</h2>
      <p>Transaction: {order.payment.transaction}</p>
      <p>Request ID: {order.payment.request_id || '-'}</p>
      <p>Currency: {order.payment.currency}</p>
      <p>Provider: {order.payment.provider}</p>
      <p>Amount: {order.payment.amount}</p>
      <p>Payment Date: {new Date(order.payment.payment_dt * 1000).toLocaleString()}</p>
      <p>Bank: {order.payment.bank}</p>
      <p>Delivery Cost: {order.payment.delivery_cost}</p>
      <p>Goods Total: {order.payment.goods_total}</p>
      <p>Custom Fee: {order.payment.custom_fee || 0}</p>

      <h2>Items</h2>
      <ul>
        {order.items.map(item => (
          <li key={item.chrt_id}>
            <strong>{item.name}</strong> (RID: {item.rid}) — {item.price} {order.payment.currency}<br />
            Size: {item.size}, Brand: {item.brand}, Sale: {item.sale || 0}, Status: {item.status}<br />
            Track Number: {item.track_number}, Total Price: {item.total_price}, NM ID: {item.nm_id}
          </li>
        ))}
      </ul>
    </div>
  );
};
