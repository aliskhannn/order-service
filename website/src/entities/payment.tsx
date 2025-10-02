export interface Payment {
  transaction: string;
  request_id?: string;
  currency: string;
  provider: string;
  amount: number;
  payment_dt: number;
  bank: string;
  delivery_cost: number;
  goods_total: number;
  custom_fee?: number;
}