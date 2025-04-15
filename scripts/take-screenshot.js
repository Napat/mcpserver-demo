const puppeteer = require('puppeteer');
const path = require('path');
const fs = require('fs');

async function wait(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function takeScreenshot() {
  console.log('Starting screenshot process...');
  const browser = await puppeteer.launch({ 
    headless: "new",
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
    defaultViewport: {width: 1280, height: 720}
  });
  const page = await browser.newPage();
  
  try {
    console.log('Navigating to login page...');
    // Navigate to login page
    await page.goto('http://localhost:8001/login', { 
      waitUntil: 'networkidle0',
      timeout: 30000 
    });
    console.log('Page loaded');
    
    // Wait for form elements to be ready
    await wait(1000); // Small delay to ensure form is interactive
    
    // Fill in email
    console.log('Filling in email...');
    await page.waitForSelector('input[type="email"]');
    await page.type('input[type="email"]', 'user@example.com', { delay: 100 });
    
    // Fill in password
    console.log('Filling in password...');
    await page.waitForSelector('input[type="password"]');
    await page.type('input[type="password"]', 'user123', { delay: 100 });
    
    // Wait a moment before clicking
    await wait(500);
    
    // Click login button
    console.log('Clicking login button...');
    await Promise.all([
      page.waitForNavigation({ waitUntil: 'networkidle0', timeout: 30000 }),
      page.click('button[type="submit"]')
    ]);
    
    // Wait for the profile page to load completely
    await wait(2000);
    
    // Ensure screenshots directory exists
    const screenshotDir = path.join(__dirname, '../screenshots');
    if (!fs.existsSync(screenshotDir)) {
      console.log('Creating screenshots directory...');
      fs.mkdirSync(screenshotDir, { recursive: true });
    }
    
    // Take screenshot
    console.log('Taking screenshot...');
    const screenshotPath = path.join(screenshotDir, 'login-result.png');
    await page.screenshot({ 
      path: screenshotPath,
      fullPage: true 
    });
    
    console.log(`Screenshot saved successfully at: ${screenshotPath}`);
  } catch (error) {
    console.error('Error details:', error);
  } finally {
    await browser.close();
    console.log('Browser closed');
  }
}

takeScreenshot();
