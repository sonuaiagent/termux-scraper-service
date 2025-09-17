from fastapi import FastAPI
import httpx
import asyncio
from bs4 import BeautifulSoup
import re
import json

app = FastAPI()

@app.post("/scrape")
async def scrape_amazon_product(data: dict):
    try:
        url = data["url"]
        
        # Enhanced headers to avoid bot detection
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'en-US,en;q=0.5',
            'Accept-Encoding': 'gzip, deflate',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Upgrade-Insecure-Requests': '1',
            'Sec-Fetch-Dest': 'document',
            'Sec-Fetch-Mode': 'navigate',
            'Sec-Fetch-Site': 'none',
            'Cache-Control': 'max-age=0'
        }
        
        # Use async HTTP client with timeout
        async with httpx.AsyncClient(timeout=30.0, follow_redirects=True) as client:
            response = await client.get(url, headers=headers)
        
        # Parse HTML using asyncio.to_thread for CPU-bound parsing
        result = await asyncio.to_thread(parse_amazon_html, response.text)
        
        return {
            "success": True,
            "title": result["title"],
            "mrp": result["mrp"],
            "discount": result["discount"],
            "price": result["price"],
            "rating": result["rating"],
            "availability": result["availability"]
        }
        
    except Exception as e:
        return {
            "success": False,
            "error": str(e)
        }

def parse_amazon_html(html):
    """Synchronous HTML parsing function"""
    soup = BeautifulSoup(html, 'html.parser')
    
    # Extract title with multiple selectors
    title_elem = (
        soup.find('span', {'id': 'productTitle'}) or
        soup.find('h1', class_='a-size-large') or
        soup.find('h1', class_='product-title')
    )
    title = title_elem.get_text().strip() if title_elem else "Not found"
    
    # Extract price with multiple selectors
    price_elem = (
        soup.find('span', class_='a-price-whole') or
        soup.find('span', class_='a-offscreen') or
        soup.find('span', id='priceblock_ourprice') or
        soup.find('span', id='priceblock_dealprice') or
        soup.find('span', class_='a-color-price')
    )
    price = price_elem.get_text().strip() if price_elem else "Not found"
    
    # Extract MRP with multiple selectors
    mrp_elem = (
        soup.find('span', class_='a-price a-text-price') or
        soup.find('span', id='listPrice') or
        soup.find('span', class_='a-text-strike')
    )
    if mrp_elem:
        mrp_span = mrp_elem.find('span', class_='a-offscreen')
        mrp = mrp_span.get_text().strip() if mrp_span else mrp_elem.get_text().strip()
    else:
        mrp = "Not found"
    
    # Extract rating with multiple selectors
    rating_elem = (
        soup.find('span', class_='a-icon-alt') or
        soup.find('span', class_='a-icon-star-alt') or
        soup.find('span', {'data-hook': 'rating-out-of-text'})
    )
    rating = rating_elem.get_text().strip() if rating_elem else "Not available"
    
    # Extract availability with multiple selectors
    availability_elem = (
        soup.find('div', {'id': 'availability'}) or
        soup.find('span', class_='a-color-success') or
        soup.find('span', class_='a-color-state') or
        soup.find('div', class_='a-section a-spacing-none')
    )
    if availability_elem:
        availability = availability_elem.get_text().strip()
        # Clean up availability text
        availability = re.sub(r'\s+', ' ', availability).strip()
    else:
        availability = "Not available"
    
    # Calculate discount
    discount = "Not found"
    if price and mrp and price != "Not found" and mrp != "Not found":
        try:
            # Extract numeric values from price strings
            price_num = float(re.sub(r'[^\d.]', '', price.replace(',', '')))
            mrp_num = float(re.sub(r'[^\d.]', '', mrp.replace(',', '')))
            
            if mrp_num > price_num:
                discount_pct = int(((mrp_num - price_num) / mrp_num) * 100)
                discount = f"-{discount_pct}%"
        except (ValueError, ZeroDivisionError):
            pass
    
    # Clean up extracted text
    title = re.sub(r'\s+', ' ', title).strip()
    price = re.sub(r'\s+', ' ', price).strip()
    mrp = re.sub(r'\s+', ' ', mrp).strip()
    rating = re.sub(r'\s+', ' ', rating).strip()
    
    return {
        "title": title,
        "price": price,
        "mrp": mrp,
        "discount": discount,
        "rating": rating,
        "availability": availability
    }

@app.get("/")
async def root():
    return {
        "message": "FastAPI Amazon Scraper with Concurrency",
        "version": "2.1",
        "features": ["async_processing", "concurrent_requests", "enhanced_parsing"]
    }

@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "service": "FastAPI Amazon Scraper",
        "concurrency": "enabled"
    }
