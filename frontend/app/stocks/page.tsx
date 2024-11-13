"use client";

import { useEffect, useState } from 'react';
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

  if (error) return <div>{error}</div>;

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
            <th>Dividends</th>
            <th>Stock Splits</th>
          </tr>
        </thead>
        <tbody>
          {stocks.map((stock) => (
            <tr key={stock.ticker}>
              <td>{stock.datetime}</td>
              <td>{stock.ticker}</td>
              <td>{stock.open}</td>
              <td>{stock.high}</td>
              <td>{stock.low}</td>
              <td>{stock.close}</td>
              <td>{stock.volume}</td>
              <td>{stock.dividends}</td>
              <td>{stock.stock_splits}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
