import type { NextPage } from 'next';
import Link from 'next/link';

const Home: NextPage = () => {
  return (
    <div>
      <h1>Welcome to the Stock Dashboard</h1>
      <p>Click below to view all stocks or go to the dashboard.</p>
      <nav>
        <ul>
          <li><Link href="/stocks">View All Stocks</Link></li>
          <li><Link href="/dashboard">Dashboard</Link></li>
        </ul>
      </nav>
    </div>
  );
};

export default Home;