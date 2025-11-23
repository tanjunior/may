import type { Product, APIFailure } from "./apiTypes";
import { CodeInternalError } from "./errorCodes";

class APIClientError extends Error {
  code: string | undefined;
  details: any;
  constructor(code: string | undefined, message: string, details?: any) {
    super(message);
    this.code = code;
    this.details = details;
  }
}

async function handleResponse<T>(res: Response): Promise<T> {
  const body = await res.json().catch(() => null);
  if (!body) throw new APIClientError(CodeInternalError, "invalid response body");
  if (body.success) {
    return body.data as T;
  }
  const err = body as APIFailure;
  throw new APIClientError(err.error?.code, err.error?.message || "api error", err.error?.details);
}

export const fetchLatestProduct = async (): Promise<Product> => {
  const response = await fetch("http://localhost:8080/product/latest");
  if (!response.ok) {
    throw new APIClientError(CodeInternalError, "Failed to fetch data");
  }
  return handleResponse<Product>(response);
};

// You can also use a promise directly without an async function wrapper
export const latestProductPromise = fetch("http://localhost:8080/product/latest").then((res) => handleResponse<Product>(res));

export const fetchProductById = async (id: number): Promise<Product> => {
  const response = await fetch(`http://localhost:8080/product/${id}`);
  if (!response.ok) {
    throw new APIClientError(CodeInternalError, "Failed to fetch data");
  }
  return handleResponse<Product>(response);
};

export const productByIdPromise = (id: number) =>
  fetch(`http://localhost:8080/product/${id}`).then((res) => handleResponse<Product>(res));

export const createProduct = async (productData: { code: string; price: number }): Promise<Product> => {
  const response = await fetch("http://localhost:8080/product", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(productData),
  });
  if (!response.ok) {
    throw new APIClientError(CodeInternalError, "Failed to create product");
  }
  return handleResponse<Product>(response);
};

export const createProductPromise = (productData: { code: string; price: number }) =>
  fetch("http://localhost:8080/product", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(productData),
  }).then((res) => handleResponse<Product>(res));

export const fetchAllProducts = async (): Promise<Product[]> => {
  const response = await fetch("http://localhost:8080/products");
  if (!response.ok) {
    throw new APIClientError(CodeInternalError, "Failed to fetch data");
  }
  return handleResponse<Product[]>(response);
};

export const allProductsPromise = fetch("http://localhost:8080/products").then((res) => handleResponse<Product[]>(res));

export const deleteProductById = async (id: number): Promise<{ message: string }> => {
  const response = await fetch(`http://localhost:8080/product/${id}`, {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new APIClientError(CodeInternalError, "Failed to delete product");
  }
  return handleResponse<{ message: string }>(response);
};

export const deleteProductByIdPromise = (id: number) =>
  fetch(`http://localhost:8080/product/${id}`, { method: "DELETE" }).then((res) => handleResponse<{ message: string }>(res));

export const updateProductById = async (id: number, productData: { code?: string; price?: number }): Promise<Product> => {
  const response = await fetch(`http://localhost:8080/product/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(productData),
  });
  if (!response.ok) {
    throw new APIClientError(CodeInternalError, "Failed to update product");
  }
  return handleResponse<Product>(response);
};

export const updateProductByIdPromise = (id: number, productData: { code?: string; price?: number }) =>
  fetch(`http://localhost:8080/product/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(productData),
  }).then((res) => handleResponse<Product>(res));

// Add more API functions as needed
