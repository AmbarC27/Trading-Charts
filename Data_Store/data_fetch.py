import yfinance as yf
import pandas as pd
import certifi
import ssl
import urllib.request
from concurrent.futures import ThreadPoolExecutor, as_completed
from sqlalchemy import create_engine
import psycopg2
from dotenv import load_dotenv
import os

def main():
    # Set up SSL context using certifi's certificates
    ssl_context = ssl.create_default_context(cafile=certifi.where())

    # Scrape the S&P 500 list from Wikipedia using urllib with the SSL context
    url = "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"
    response = urllib.request.urlopen(url, context=ssl_context)
    html = response.read()

    # Use pandas to parse the HTML and extract the table
    table = pd.read_html(html)

    # Extract the first table (which contains the ticker symbols)
    sp500_df = table[0]

    # Get the list of ticker symbols
    tickers = sp500_df['Symbol'].tolist()

    def fetch_ticker_data(ticker_symbol,period="1d",interval="1m"):

        # Fetch the ticker data
        ticker = yf.Ticker(ticker_symbol)

        # Get historical market data (trading info)
        trading_info = ticker.history(period=period,interval=interval)
        trading_info['ticker'] = ticker_symbol

        return trading_info

    def fetch_data_parallel(tickers, max_workers=10):
        all_data = pd.DataFrame()

        # Use ThreadPoolExecutor to fetch data in parallel
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            # Submit the tasks to the thread pool
            futures = {executor.submit(fetch_ticker_data, ticker): ticker for ticker in tickers}

            # As the tasks complete, process the results
            for future in as_completed(futures):
                ticker = futures[future]
                try:
                    # Get the result of the completed task
                    trading_info = future.result()
                    trading_info.reset_index(inplace=True)
                    if not trading_info.empty:  # Only append if data is returned
                        all_data = pd.concat([all_data, trading_info],ignore_index=True)
                except Exception as e:
                    print(f"Error processing data for {ticker}: {e}")

        return all_data

    all_data = fetch_data_parallel(tickers, max_workers=10)
    all_data.columns = all_data.columns.str.lower()
    all_data.drop_duplicates(subset=['datetime', 'ticker'], keep='first',inplace=True)

    memory_usage = all_data.memory_usage(deep=True).sum()/(1024**2)
    print('memory_usage:', memory_usage, 'MB'), 

    load_dotenv()

    db_user = os.getenv('DB_USER')
    db_password = os.getenv('DB_PASSWORD')
    db_host = os.getenv('DB_HOST')
    db_port = os.getenv('DB_PORT')
    db_name = os.getenv('DB_NAME')
    db_table = os.getenv('DB_TABLE')

    # Create SQLAlchemy engine
    engine = create_engine(f'postgresql+psycopg2://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}')

    all_data.to_sql(db_table, engine, if_exists='append', index=False)

if __name__ == "__main__":
    main()