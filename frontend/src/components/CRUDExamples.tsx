import React, { useState } from "react";
import {
  createProduct,
  fetchProductById,
  fetchAllProducts,
  updateProductById,
  deleteProductById,
} from "../apis";
import type { Product } from "../apiTypes";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";

export default function CRUDExamples() {
  const [createState, setCreateState] = useState<{
    code: string;
    price: string;
  }>({ code: "", price: "" });
  const [readId, setReadId] = useState<string>("");
  const [readResult, setReadResult] = useState<Product | null>(null);
  const [listResult, setListResult] = useState<Product[] | null>(null);
  const [updateState, setUpdateState] = useState<{
    id: string;
    code: string;
    price: string;
  }>({ id: "", code: "", price: "" });
  const [deleteId, setDeleteId] = useState<string>("");
  const [message, setMessage] = useState<string | null>(null);

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    try {
      const p = await createProduct({
        code: createState.code,
        price: Number(createState.price),
      });
      setMessage(`Created id=${p.id} code=${p.code} price=${p.price}`);
      setCreateState({ code: "", price: "" });
    } catch (err: any) {
      setMessage(err?.message || String(err));
    }
  };

  const onRead = async (e?: React.FormEvent) => {
    e?.preventDefault();
    setMessage(null);
    setReadResult(null);
    try {
      const id = Number(readId);
      if (!id) throw new Error("Enter a numeric id");
      const p = await fetchProductById(id);
      console.log(p)
      setReadResult(p);
    } catch (err: any) {
      setMessage(err?.message || String(err));
    }
  };

  const onList = async () => {
    setMessage(null);
    try {
      const all = await fetchAllProducts();
      setListResult(all);
    } catch (err: any) {
      setMessage(err?.message || String(err));
    }
  };

  const onUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    try {
      const id = Number(updateState.id);
      if (!id) throw new Error("Enter a numeric id to update");
      const payload: any = {};
      if (updateState.code) payload.code = updateState.code;
      if (updateState.price) payload.price = Number(updateState.price);
      const p = await updateProductById(id, payload);
      setMessage(`Updated id=${p.id} code=${p.code} price=${p.price}`);
      setUpdateState({ id: "", code: "", price: "" });
    } catch (err: any) {
      setMessage(err?.message || String(err));
    }
  };

  const onDelete = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    try {
      const id = Number(deleteId);
      if (!id) throw new Error("Enter a numeric id to delete");
      const res = await deleteProductById(id);
      setMessage(res.message || "Deleted");
      setDeleteId("");
    } catch (err: any) {
      setMessage(err?.message || String(err));
    }
  };

  return (
    <div className="w-full">
      <h2 className="text-lg font-semibold text-slate-800">CRUD Examples</h2>
      <div className="mt-4 grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <CardTitle>Create</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="flex gap-2 flex-wrap" onSubmit={onCreate}>
              <Input
                className="w-auto"
                placeholder="code"
                value={createState.code}
                onChange={(e) =>
                  setCreateState({ ...createState, code: e.target.value })
                }
              />
              <Input
                className="w-auto"
                placeholder="price"
                value={createState.price}
                onChange={(e) =>
                  setCreateState({ ...createState, price: e.target.value })
                }
              />
              <Button type="submit">Create</Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Read (by id)</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="flex gap-2 flex-wrap" onSubmit={onRead}>
              <Input
                className="w-auto"
                placeholder="id"
                value={readId}
                onChange={(e) => setReadId(e.target.value)}
              />
              <Button
                type="submit"
                className="bg-transparent hover:bg-slate-100 text-slate-700 border border-slate-200"
              >
                Fetch
              </Button>
              <Button
                type="button"
                className="bg-transparent hover:bg-slate-100 text-slate-700 border border-slate-200"
                onClick={onList}
              >
                List All
              </Button>
            </form>
            {readResult && (
              <pre className="mt-2 p-2 bg-slate-900/5 rounded-md text-xs overflow-auto">
                {JSON.stringify(readResult, null, 2)}
              </pre>
            )}
            {listResult && (
              <pre className="mt-2 p-2 bg-slate-900/5 rounded-md text-xs overflow-auto">
                {JSON.stringify(listResult, null, 2)}
              </pre>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Update</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="flex gap-2 flex-wrap" onSubmit={onUpdate}>
              <Input
                className="w-auto"
                placeholder="id"
                value={updateState.id}
                onChange={(e) =>
                  setUpdateState({ ...updateState, id: e.target.value })
                }
              />
              <Input
                className="w-auto"
                placeholder="new code (optional)"
                value={updateState.code}
                onChange={(e) =>
                  setUpdateState({ ...updateState, code: e.target.value })
                }
              />
              <Input
                className="w-auto"
                placeholder="new price (optional)"
                value={updateState.price}
                onChange={(e) =>
                  setUpdateState({ ...updateState, price: e.target.value })
                }
              />
              <Button
                type="submit"
                className="bg-transparent hover:bg-slate-100 text-slate-700 border border-slate-200"
              >
                Update
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Delete</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="flex gap-2 flex-wrap" onSubmit={onDelete}>
              <Input
                className="w-auto"
                placeholder="id"
                value={deleteId}
                onChange={(e) => setDeleteId(e.target.value)}
              />
              <Button
                type="submit"
                className="bg-red-500 hover:bg-red-600 text-white"
              >
                Delete
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>

      {message && (
        <div className="mt-4 p-3 rounded-md bg-blue-50 border border-blue-100 text-sm text-blue-800">
          {message}
        </div>
      )}
    </div>
  );
}
