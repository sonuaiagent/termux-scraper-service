from flask import Flask, request, jsonify
from datetime import datetime
import pytz
import requests
from bs4 import BeautifulSoup

app = Flask(__name__)

def get_ist_time():
    """Get current time in IST format"""
    ist = pytz.timezone('Asia/Kolkata')
    now = datetime.now(ist)
    return now.strftime('%Y-%m-%d %H:%M:%S IST')

@app.route('/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({
        "status": "ok",
        "service": "Termux Scraper Service",
        "timestamp": get_ist_time()
    })

def scrape_amazon_product(url):
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
        "Accept-Language": "en-US,en;q=0.9",
        "Accept-Encoding": "gzip, deflate, br"
    }
    try:
        response = requests.get(url, headers=headers, timeout=30)
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
        
        # Extract product details
        product_details = {}
        detail_div = soup.find('div', id='detailBullets_feature_div')
        if detail_div:
            lis = detail_div.find_all('li')
            for li in lis:
                spans = li.find_all('span')
                if len(spans) >= 2:
                    key = spans[0].get_text(strip=True).replace(':', '').replace('‏', '').replace('‎', '')
                    value = spans[1].get_text(strip=True)
                    product_details[key] = value
        
        # Extract all possible prices
        all_prices = []
        price_selectors = [
            "span.a-price span.a-offscreen",
            "span#priceblock_ourprice",
            "span#priceblock_dealprice",
            "span#price_inside_buybox",
            "span.a-color-price",
            "span.offer-price",
            "span.a-price-whole",
            "span.a-price-fraction",
            "div#corePrice_feature_div span.a-offscreen",
            "span#priceblock_saleprice",
            "span#priceblock_regularprice",
            "div#averageCustomerReviews span.a-price"
        ]
        
        for selector in price_selectors:
            candidates = soup.select(selector)
            for candidate in candidates:
                text = candidate.get_text(strip=True)
                if text and text.startswith('₹'):
                    try:
                        clean_text = text.replace('₹', '').replace(',', '').split('.')[0]
                        numeric_value = int(clean_text)
                        all_prices.append((numeric_value, text))
                    except ValueError:
                        pass
        
        # Find specific prices
        core_price_tag = soup.select_one("div#corePrice_feature_div span.a-offscreen")
        if core_price_tag:
            final_price = core_price_tag.get_text(strip=True)
        else:
            final_price = "Price Not Found"
        
        mrp_price = "MRP Not Found"
        if all_prices:
            all_prices.sort()
            mrp_price = all_prices[-1][1]  # Largest price
        
        price = final_price
        if mrp == "MRP Not Found":
            mrp = mrp_price
        
        return {
            "title": title,
            "mrp": mrp,
            "discount": discount,
            "price": price,
            "rating": rating,
            "availability": availability,
            "product_details": product_details,
            "url": url
        }
        
    except requests.RequestException as e:
        return {"error": f"Request error: {e}"}
    except Exception as e:
        return {"error": f"An error occurred: {e}"}

@app.route('/scrape', methods=['POST'])
def scrape():
    """Main scraper endpoint"""
    try:
        # Get JSON data from request
        data = request.get_json(force=True) or {}
        
        # Extract command or URL
        command = data.get("command", "").lower().strip()
        url = data.get("url", "").strip()
        chat_id = data.get("chat_id", "unknown")
        username = data.get("username", "user")
        
        # Handle "hi" command - FIXED F-STRING
        if command == "hi" or url == "hi":
            response_msg = (
                f"Hello! I am scraper from Termux 🤖

"
                f"Current time: {get_ist_time()}
"
                f"Chat ID: {chat_id}
"
                f"Username: @{username}"
            )
            
            return jsonify({
                "success": True,
                "message": response_msg,
                "product_info": {
                    "title": "Termux Scraper Service",
                    "price": "Active",
                    "timestamp": get_ist_time()
                },
                "debug": [f"Received command: {command or url}", f"From user: @{username}"]
            })
        
        # Handle Amazon URL scraping
        elif url and "amazon" in url:
            result = scrape_amazon_product(url)
            if "error" in result:
                return jsonify({
                    "success": False,
                    "error": result["error"],
                    "timestamp": get_ist_time()
                })
            else:
                return jsonify({
                    "success": True,
                    "product": result,
                    "timestamp": get_ist_time()
                })
        
        # For any other input, return a friendly message
        else:
            return jsonify({
                "success": False,
                "error": "Send 'hi' to test the scraper service or provide a valid Amazon URL",
                "message": f"Unknown command or URL: {command or url}",
                "debug": [f"Service is running at {get_ist_time()}"]
            })
            
    except Exception as e:
        return jsonify({
            "success": False,
            "error": f"Service error: {str(e)}",
            "timestamp": get_ist_time()
        }), 500

if __name__ == "__main__":
    print(f"🚀 Starting Termux Scraper Service at {get_ist_time()}")
    print("📱 Service will run on http://0.0.0.0:5000")
    print("🔗 Use Cloudflare Tunnel to expose this service")
    print("💡 Send 'hi' command to test the service")
    app.run(host="0.0.0.0", port=5000, debug=False)
