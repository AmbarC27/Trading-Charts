"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import './style.css';

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

export default function Stocks() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchStocks = async () => {
      try {
        const response = await fetch('http://localhost:8080/latest');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log(data)
        setStocks(data);
      } catch (error) {
        setError("Failed to fetch data from the server. Please try again.");
        console.error("Fetch error:", error);
      }
    };
    fetchStocks();
  }, []);

  if (error || stocks.length === 0) return <div>{error}</div>;

  return (
    <div>
      <h1>Stocks list</h1>
      <table>
        <thead>
          <tr>
            <th>Datetime</th>
            <th>Ticker</th>
            <th>Open</th>
            <th>High</th>
            <th>Low</th>
            <th>Close</th>
            <th>Volume</th>
            {/* <th>Dividends</th>
            <th>Stock Splits</th> */}
            <th>Diff</th>
          </tr>
        </thead>
        <tbody>
          {stocks.map((stock) => {
            const diff = stock.close - stock.open;
            return (
              <tr 
                key={stock.ticker}
                onClick={() => router.push(`/stocks/${stock.ticker}`)}
              >
                <td>{new Date(stock.datetime).toLocaleString()}</td>
                <td>{stock.ticker}</td>
                <td>{stock.open.toFixed(2)}</td>
                <td>{stock.high.toFixed(2)}</td>
                <td>{stock.low.toFixed(2)}</td>
                <td>{stock.close.toFixed(2)}</td>
                <td>{stock.volume}</td>
                {/* <td>{stock.dividends.toFixed(2)}</td>
                <td>{stock.stock_splits.toFixed(2)}</td> */}
                <td>{diff.toFixed(2)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
