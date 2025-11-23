import { use } from "react";
import { latestProductPromise } from "../apis";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "./ui/card";

export default function LatestProduct() {
  // The use() hook reads the promise value directly
  const {ID: id, Code: code, Price: price} = use(latestProductPromise);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Latest product</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-slate-700 space-y-1">
          <div>
            <span className="font-medium">ID:</span> {id}
          </div>
          <div>
            <span className="font-medium">Code:</span> {code}
          </div>
          <div>
            <span className="font-medium">Price:</span> {price}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
