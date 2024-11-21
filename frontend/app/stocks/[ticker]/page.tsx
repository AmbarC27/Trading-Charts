"use client";

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { Line } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend } from 'chart.js';

// Register Chart.js components
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend);

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

export default function StockDetails() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [filteredStocks, setFilteredStocks] = useState<Stock[]>([]);
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const { ticker } = useParams();

  // Fetch all stock data for the ticker on initial load
  useEffect(() => {
    const fetchStocks = async () => {
      try {
        const response = await fetch(`http://localhost:8080/stocks/${ticker}`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data: Stock[] = await response.json();
        setStocks(data);
        setFilteredStocks(data); // Initially, show all data
      } catch (error) {
        setError("Failed to fetch stock data. Please try again.");
        console.error("Fetch error:", error);
      }
    };

    fetchStocks();
  }, [ticker]);

  // Function to handle date range filtering
  const filterStocksInRange = () => {
    if (!startDate || !endDate) {
      setError("Please select both start and end dates.");
      return;
    }

    const filteredData = stocks.filter((stock) => {
      const stockDate = new Date(stock.datetime);
      const start = new Date(startDate);
      const end = new Date(endDate);

      return stockDate >= start && stockDate <= end;
    });

    if (filteredData.length === 0) {
      setError("No data available for the selected date range.");
    } else {
      setFilteredStocks(filteredData);
      setError(null); // Clear error if data is found
    }
  };

  // Chart Data
  const chartData = {
    labels: filteredStocks.map(stock => new Date(stock.datetime).toLocaleString()),
    datasets: [
      {
        label: `Closing Prices (${ticker})`,
        data: filteredStocks.map(stock => stock.close),
        borderColor: 'rgba(75,192,192,1)',
        backgroundColor: 'rgba(75,192,192,0.2)',
        pointRadius: 3,
        borderWidth: 2,
      },
    ],
  };

  // Chart Options
  const chartOptions = {
    responsive: true,
    plugins: {
      legend: { display: true },
      tooltip: { mode: 'index', intersect: false },
    },
    scales: {
      x: { title: { display: true, text: 'Datetime' } },
      y: { title: { display: true, text: 'Closing Price' } },
    },
  };

  return (
    <div>
      <h1>Stock Closing Prices: {ticker}</h1>

      <div style={{ marginBottom: '20px' }}>
        <label>
          Start Date: 
          <input 
            type="datetime-local" 
            value={startDate} 
            onChange={(e) => setStartDate(e.target.value)} 
            style={{ marginLeft: '10px', marginRight: '20px' }}
          />
        </label>
        <label>
          End Date: 
          <input 
            type="datetime-local" 
            value={endDate} 
            onChange={(e) => setEndDate(e.target.value)} 
            style={{ marginLeft: '10px' }}
          />
        </label>
        <button 
          onClick={filterStocksInRange} 
          style={{ 
            marginLeft: '20px', 
            backgroundColor: '#FFF',
            border: 'none', // Optional: Remove border for a cleaner look
            padding: '10px 20px', // Optional: Adjust padding for better appearance
            cursor: 'pointer' // Optional: Change cursor to pointer on hover
          }}
        >
          Search in Following Range
        </button>
      </div>

      {error && <div style={{ color: 'red', marginBottom: '20px' }}>{error}</div>}

      {filteredStocks.length > 0 ? (
        <Line data={chartData} options={chartOptions} />
      ) : (
        <div>No data available for the selected range.</div>
      )}
    </div>
  );
}
