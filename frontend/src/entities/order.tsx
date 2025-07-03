import type { Delivery } from "./delivery";
import type { Item } from "./item";
import type { Payment } from "./payment";

export interface Order {
  order_uid: string; // uuid.UUID в Go — обычно string в TS
  track_number: string;
  entry: string;
  delivery: Delivery;
  payment: Payment;
  items: Item[];
  locale: string;
  internal_signature?: string;
  customer_id: string;
  delivery_service: string;
  shardkey: string;
  sm_id: number;
  date_created: string; // time.Time в Go обычно ISO string
  oof_shard: string;
}
