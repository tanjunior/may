import { use } from "react";
import { latestProductPromise } from "../apis";

export default function LatestProduct() {
  // The use() hook reads the promise value directly
  const product = use(latestProductPromise);
  const { ID: id, Code: code, Price: price } = product.product;

  return (
    <div>
      <h1>id: {id}</h1>
      <p>code: {code}</p>
      <p>price: {price}</p>
    </div>
  );
}
