// GENERATED: API response types for frontend

export interface APIError {
  code: string;
  message: string;
  details?: any;
}

export interface ErrorEnvelope {
  success: false;
  status: number;
  error: APIError;
}

export interface SuccessEnvelope<T> {
  success: true;
  status: number;
  data: T;
  meta?: Record<string, any>;
}

export interface Product {
  id: number;
  code: string;
  price: number;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
}

export type ProductListResponse = SuccessEnvelope<Product[]>;
export type ProductResponse = SuccessEnvelope<Product>;
export type APIFailure = ErrorEnvelope;
