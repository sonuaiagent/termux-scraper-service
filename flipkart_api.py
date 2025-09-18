from fastapi import FastAPI
import undetected_chromedriver as uc
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium_stealth import stealth
from fake_useragent import UserAgent
import asyncio
import time
import random
import re
from datetime import datetime

app = FastAPI()

def log_flipkart_product_info(product_data):
    """Beautiful product information logging for Flipkart"""
    border = "═" * 79
    print(f"\n╔{border}╗")
    print("║                    🛍️ FLIPKART PRODUCT FOUND (UNDETECTED)                   ║")
    print(f"╠{border}╣")
    print("║                                                                             ║")
    
    # Truncate long product names
    name = product_data.get('name', 'N/A')
    if len(name) > 55:
        name = name[:52] + "..."
    print(f"║ 🏷️  Product Name: {name:<55} ║")
    
    print(f"║ 💰 Current Price: {product_data.get('price', 'N/A'):<55} ║")
    print(f"║ ⭐ Rating:        {product_data.get('rating', 'N/A'):<55} ║")
    print(f"║ 🛡️  Anti-Bot:      Undetected ChromeDriver BYPASSED                        ║")
    
    print("║                                                                             ║")
    print(f"║ 📅 Processed: {datetime.now().strftime('%Y-%m-%d %H:%M:%S IST'):<60} ║")
    print(f"╚{border}╝\n")

async def create_undetected_driver():
    """Create undetected ChromeDriver with advanced stealth"""
    print("🛡️ Creating undetected ChromeDriver...")
    
    # Configure Chrome options for maximum stealth
    options = uc.ChromeOptions()
    
    # Basic stealth options
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-blink-features=AutomationControlled")
    options.add_argument("--disable-extensions")
    options.add_argument("--disable-plugins")
    options.add_argument("--disable-images")  # Faster loading
    options.add_argument("--window-size=1920,1080")
    
    # Random user agent
    ua = UserAgent()
    options.add_argument(f"--user-agent={ua.random}")
    
    # Additional anti-detection
    options.add_argument("--disable-web-security")
    options.add_argument("--allow-running-insecure-content")
    options.add_argument("--disable-features=VizDisplayCompositor")
    
    # Create undetected driver
    driver = uc.Chrome(options=options, version_main=None)
    
    # Apply additional stealth
    stealth(driver,
        languages=["en-US", "en", "hi"],
        vendor="Google Inc.",
        platform="Linux x86_64",
        webgl_vendor="Intel Inc.",
        renderer="Intel Iris OpenGL Engine",
        fix_hairline=True,
    )
    
    return driver

@app.post("/scrape")
async def scrape_flipkart_product(data: dict):
    driver = None
    try:
        url = data["url"]
        print(f"🛡️ Undetected Flipkart: Starting advanced anti-bot bypass...")
        print(f"🌐 Target URL: {url}")
        
        # Create undetected driver
        driver = await asyncio.to_thread(create_undetected_driver)
        print("✅ Undetected ChromeDriver created successfully")
        
        # Navigate with human-like behavior
        print("🚀 Navigating to Flipkart...")
        driver.get(url)
        
        # Wait for page load with random delay
        random_delay = random.uniform(3, 7)
        print(f"⏳ Waiting {random_delay:.1f}s for page load...")
        time.sleep(random_delay)
        
        # Human-like mouse movement
        actions = driver.execute_script("""
            function randomMouseMove() {
                var event = new MouseEvent('mousemove', {
                    'view': window,
                    'bubbles': true,
                    'cancelable': true,
                    'clientX': Math.random() * window.innerWidth,
                    'clientY': Math.random() * window.innerHeight
                });
                document.dispatchEvent(event);
            }
            randomMouseMove();
            return true;
        """)
        
        # Extract product information with multiple selectors
        result = await asyncio.to_thread(extract_flipkart_data, driver)
        
        print(f"✅ Undetected Flipkart: Product extraction completed")
        
        # Log beautiful product information
        log_flipkart_product_info(result)
        
        return {
            "success": True,
            "name": result["name"],
            "price": result["price"],
            "rating": result["rating"]
        }
        
    except Exception as e:
        print(f"❌ Undetected Flipkart: Scraping failed - {str(e)}")
        return {
            "success": False,
            "error": str(e)
        }
    finally:
        if driver:
            try:
                driver.quit()
                print("🔒 ChromeDriver closed safely")
            except:
                pass

def extract_flipkart_data(driver):
    """Extract product data with advanced selector fallbacks"""
    
    # Product Name - Comprehensive selector list
    name_selectors = [
        'span.B_NuCI',                    # Primary title
        'h1.x-oua-w9',                    # Alternative title
        'h1._35KyD6',                     # Backup title
        'span.VU-ZEz',                    # New layout
        'span._35KyD6',                   # Alternative span
        'div._4rR01T',                    # Container title
        '.B_NuCI',                        # Class only
        '[data-testid*="title"]',         # Data attributes
        'h1[class*="title"]',             # Partial class match
        'span[class*="title"]'            # Span title variants
    ]
    
    product_name = "Product name not found"
    for selector in name_selectors:
        try:
            element = driver.find_element(By.CSS_SELECTOR, selector)
            if element and element.text.strip():
                product_name = element.text.strip()
                print(f"✅ Product name found with: {selector}")
                break
        except:
            continue
    
    # Product Price - Comprehensive selector list
    price_selectors = [
        'div._30jeq3._16Jk6d',           # Primary price
        'div._30jeq3',                   # Alternative price
        'div.Nx9bqj.CxhGGd',            # Price container
        'div._1_WHN1',                   # New price layout
        'span._30jeq3._16Jk6d',         # Span price
        'div._25b18c',                   # Updated selector
        '._30jeq3',                      # Class only
        '[data-testid*="price"]',        # Data attributes
        'div[class*="price"]',           # Partial class
        'span[class*="price"]'           # Span price variants
    ]
    
    product_price = "Price not found"
    for selector in price_selectors:
        try:
            element = driver.find_element(By.CSS_SELECTOR, selector)
            if element and element.text.strip():
                text = element.text.strip()
                # Check if it contains price indicators
                if '₹' in text or 'Rs' in text or any(char.isdigit() for char in text):
                    product_price = text
                    print(f"✅ Price found with: {selector}")
                    break
        except:
            continue
    
    # Product Rating - Comprehensive selector list
    rating_selectors = [
        'div._3LWZlK',                   # Primary rating
        'span._1lRcqv',                  # Rating span
        'div.gUuXy-',                    # Rating container
        'div._3aeaXv',                   # Alternative rating
        'span._2_R_DZ',                  # Rating text
        'div._13vcmD',                   # New rating layout
        '._3LWZlK',                      # Class only
        '[data-testid*="rating"]',       # Data attributes
        'div[class*="rating"]',          # Partial class
        'span[class*="star"]'            # Star ratings
    ]
    
    product_rating = "Rating not available"
    for selector in rating_selectors:
        try:
            element = driver.find_element(By.CSS_SELECTOR, selector)
            if element and element.text.strip():
                text = element.text.strip()
                # Validate rating format
                if any(char.isdigit() for char in text) and (len(text) < 20):
                    product_rating = text
                    print(f"✅ Rating found with: {selector}")
                    break
        except:
            continue
    
    # Clean extracted data
    product_name = re.sub(r'\s+', ' ', product_name).strip()
    product_price = re.sub(r'\s+', ' ', product_price).strip()
    product_rating = re.sub(r'\s+', ' ', product_rating).strip()
    
    return {
        "name": product_name,
        "price": product_price,
        "rating": product_rating
    }

@app.get("/")
async def root():
    return {
        "message": "Undetected ChromeDriver Flipkart Scraper",
        "version": "3.0",
        "features": ["undetected_chromedriver", "stealth_mode", "anti_bot_bypass", "advanced_selectors"],
        "protection": "Bypasses Cloudflare, Akamai, and other anti-bot systems"
    }

@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "service": "Undetected Flipkart Scraper",
        "anti_bot": "advanced_bypass_enabled"
    }
