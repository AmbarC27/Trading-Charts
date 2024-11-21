"use client";

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
// import './style.css';

interface Stock {
  datetime: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  dividends: number;
  stock_splits: number;
  ticker: string;
}

export default function StockList() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [error, setError] = useState<string | null>(null);
  const { ticker } = useParams();

  useEffect(() => {
    const fetchStocks = async () => {
      try {
        const response = await fetch(`http://localhost:8080/stocks/${ticker}`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data: Stock[] = await response.json();
        setStocks(data);
      } catch (error) {
        setError("Failed to fetch stock data. Please try again.");
        console.error("Fetch error:", error);
      }
    };
    if (ticker) {
      fetchStocks();
    }
  }, [ticker]);

  if (error) return <div>{error}</div>;
  if (stocks.length === 0) return <div>Loading...</div>;

  return (
    <div>
      <h1>Stock History for: {ticker}</h1>
      <table>
        <thead>
          <tr>
            <th>Datetime</th>
            <th>Open</th>
            <th>High</th>
            <th>Low</th>
            <th>Close</th>
            <th>Volume</th>
            <th>Diff</th>
          </tr>
        </thead>
        <tbody>
          {stocks.map((stock, index) => {
            const diff = stock.close - stock.open;
            return (
              <tr key={index}>
                <td>{new Date(stock.datetime).toLocaleString()}</td>
                <td>{stock.open.toFixed(2)}</td>
                <td>{stock.high.toFixed(2)}</td>
                <td>{stock.low.toFixed(2)}</td>
                <td>{stock.close.toFixed(2)}</td>
                <td>{stock.volume}</td>
                <td>{diff.toFixed(2)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
