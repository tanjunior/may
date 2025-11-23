// services/api.js
export const fetchLatestProduct = async () => {
  const response = await fetch('http://localhost:8080/product/latest');
  if (!response.ok) {
    throw new Error('Failed to fetch data');
  }
  return response.json();
};

// You can also use a promise directly without an async function wrapper
export const latestProductPromise = fetch('http://localhost:8080/product/latest').then((res) => res.json());
