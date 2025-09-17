from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
import traceback
import sys

app = FastAPI()

# Add exception handler to catch all errors
@app.exception_handler(Exception)
async def generic_exception_handler(request: Request, exc: Exception):
    error_traceback = traceback.format_exc()
    print(f"❌ FastAPI Exception: {str(exc)}")
    print(f"📋 Traceback: {error_traceback}")
    
    return JSONResponse(
        status_code=500,
        content={
            "success": False, 
            "error": str(exc),
            "traceback": error_traceback
        }
    )

@app.post("/scrape")
async def scrape(request: Request):
    try:
        payload = await request.json()
        url = payload.get("url")
        
        if not url:
            return {"success": False, "error": "Missing 'url' in request"}
        
        print(f"📦 FastAPI received URL: {url}")
        
        # Try to import and call the scraper function
        from amazon import scrape_amazon_product
        print("✅ Successfully imported scrape_amazon_product")
        
        result = scrape_amazon_product(url)
        print(f"✅ Scraping completed: {result}")
        
        return result
        
    except ImportError as e:
        error_msg = f"Failed to import amazon module: {str(e)}"
        print(f"❌ Import Error: {error_msg}")
        return {"success": False, "error": error_msg}
        
    except Exception as e:
        error_msg = f"Scraping failed: {str(e)}"
        print(f"❌ Scraping Error: {error_msg}")
        print(f"📋 Traceback: {traceback.format_exc()}")
        return {"success": False, "error": error_msg}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("amazon_api:app", host="0.0.0.0", port=8081, reload=True)
