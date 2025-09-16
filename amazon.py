#!/usr/bin/env python3
import requests
from bs4 import BeautifulSoup
import json
import sys

def scrape_amazon_product(url):
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
        "Accept-Language": "en-US,en;q=0.9",
        "Accept-Encoding": "gzip, deflate, br"
    }
    
    try:
        # Debug output to stderr (won't interfere with JSON)
        print("📦 Making request to Amazon...", file=sys.stderr)
        response = requests.get(url, headers=headers, timeout=30)
        print(f"📊 Response status: {response.status_code}", file=sys.stderr)
        print(f"📐 Response length: {len(response.text)} characters", file=sys.stderr)
        
        response.raise_for_status()
        soup = BeautifulSoup(response.text, "html.parser")

        # Extract product title
        title_tag = soup.select_one("span#productTitle")
        title = title_tag.get_text(strip=True) if title_tag else "Title Not Found"

        # Extract MRP (original price)
        mrp_selectors = [
            "span.a-price.a-text-price span.a-offscreen",
            "span#priceblock_mrp",
            "span.a-text-price",
            "span.a-price.a-text-price"
        ]
        mrp = "MRP Not Found"
        for selector in mrp_selectors:
            mrp_tag = soup.select_one(selector)
            if mrp_tag:
                mrp_text = mrp_tag.get_text(strip=True)
                if mrp_text:
                    mrp = mrp_text
                    break

        # Extract discount percentage
        discount_tag = soup.select_one("span.savingsPercentage") or soup.select_one("span.a-color-price")
        discount = discount_tag.get_text(strip=True) if discount_tag else "Discount Not Found"

        # Extract rating
        rating_tag = soup.select_one("span.a-icon-alt") or soup.select_one("span#acrPopover")
        rating = rating_tag.get_text(strip=True) if rating_tag else "Rating Not Found"

        # Extract availability
        availability_tag = soup.select_one("div#availability span") or soup.select_one("span#availability")
        availability = availability_tag.get_text(strip=True) if availability_tag else "Availability Not Found"

        # Extract final price
        core_price_tag = soup.select_one("div#corePrice_feature_div span.a-offscreen")
        if core_price_tag:
            final_price = core_price_tag.get_text(strip=True)
        else:
            final_price = "Price Not Found"

        # Debug output to stderr
        print("✅ Amazon scraping completed successfully!", file=sys.stderr)
        print(f"📱 Title: {title}", file=sys.stderr)
        print(f"💰 Price: {final_price}", file=sys.stderr)

        return {
            "success": True,
            "title": title,
            "mrp": mrp,
            "discount": discount,
            "price": final_price,
            "rating": rating,
            "availability": availability
        }
    except Exception as e:
        print(f"❌ Amazon scraping error: {e}", file=sys.stderr)
        return {
            "success": False,
            "error": f"An error occurred: {e}"
        }

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(json.dumps({"success": False, "error": "Usage: python3 amazon.py <amazon_url>"}))
        sys.exit(1)
    
    url = sys.argv[1]
    result = scrape_amazon_product(url)
    
    # Output ONLY clean JSON to stdout (main.go will parse this)
    print(json.dumps(result))
