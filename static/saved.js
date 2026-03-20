const COOKIE_NAME = 'vabackpay_saved';
const MAX_ITEMS = 10;

function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}

function setCookie(name, value, days = 365) {
    const expires = new Date(Date.now() + days * 864e5).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`;
}

function loadSavedItems() {
    const cookie = getCookie(COOKIE_NAME);
    if (!cookie) return [];
    try {
        return JSON.parse(decodeURIComponent(cookie));
    } catch (e) {
        return [];
    }
}

function saveAllItems(items) {
    setCookie(COOKIE_NAME, JSON.stringify(items));
}

function saveItem(item) {
    const items = loadSavedItems();
    const newItem = {
        id: generateId(),
        savedAt: new Date().toISOString(),
        ...item
    };
    items.unshift(newItem);
    if (items.length > MAX_ITEMS) {
        items.pop();
    }
    saveAllItems(items);
    return newItem;
}

function updateItem(id, updates) {
    const items = loadSavedItems();
    const index = items.findIndex(item => item.id === id);
    if (index !== -1) {
        items[index] = { ...items[index], ...updates };
        saveAllItems(items);
        return items[index];
    }
    return null;
}

function deleteItem(id) {
    const items = loadSavedItems();
    const filtered = items.filter(item => item.id !== id);
    saveAllItems(filtered);
}

function clearAllItems() {
    setCookie(COOKIE_NAME, '');
}

function getSavedItem(id) {
    const items = loadSavedItems();
    return items.find(item => item.id === id) || null;
}

function getSavedBackpayItems() {
    return loadSavedItems().filter(item => item.type === 'backpay');
}

function getSavedDaysSinceItems() {
    return loadSavedItems().filter(item => item.type === 'dayssince');
}

function formatCurrency(amount) {
    return '$' + amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatDate(month, day, year) {
    const months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    return `${months[month - 1]} ${day}, ${year}`;
}

function calculateDaysSince(month, day, year) {
    const selectedDate = new Date(year, month - 1, day);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    selectedDate.setHours(0, 0, 0, 0);
    const diffTime = today - selectedDate;
    return Math.floor(diffTime / (1000 * 60 * 60 * 24));
}

const rates = {
    2026: {10: 180.42, 20: 356.66, 30: 552.47, 40: 795.84, 50: 1132.90, 60: 1435.02, 70: 1808.45, 80: 2102.15, 90: 2362.30, 100: 3938.58},
    2025: {10: 175.51, 20: 346.95, 30: 537.42, 40: 774.16, 50: 1102.04, 60: 1395.93, 70: 1759.19, 80: 2044.89, 90: 2297.96, 100: 3831.30},
    2024: {10: 171.25, 20: 337.88, 30: 523.43, 40: 754.08, 50: 1074.19, 60: 1361.50, 70: 1715.71, 80: 1994.83, 90: 2241.39, 100: 3737.85},
    2023: {10: 165.92, 20: 327.42, 30: 507.40, 40: 731.00, 50: 1040.84, 60: 1319.91, 70: 1662.75, 80: 1933.06, 90: 2172.39, 100: 3621.95},
    2022: {10: 152.64, 20: 301.74, 30: 467.35, 40: 673.28, 50: 958.38, 60: 1215.28, 70: 1531.30, 80: 1780.67, 90: 2000.97, 100: 3332.06},
    2021: {10: 144.14, 20: 284.93, 30: 441.35, 40: 635.77, 50: 905.04, 60: 1146.39, 70: 1444.71, 80: 1679.35, 90: 1887.18, 100: 3146.42},
    2020: {10: 142.29, 20: 281.27, 30: 435.69, 40: 627.61, 50: 893.43, 60: 1131.68, 70: 1426.17, 80: 1657.80, 90: 1862.96, 100: 3106.04},
    2019: {10: 140.05, 20: 276.84, 30: 428.83, 40: 617.73, 50: 879.36, 60: 1113.86, 70: 1403.71, 80: 1631.69, 90: 1833.62, 100: 3057.13},
    2018: {10: 136.24, 20: 269.30, 30: 417.15, 40: 600.90, 50: 855.41, 60: 1083.52, 70: 1365.48, 80: 1587.25, 90: 1783.68, 100: 2973.86},
    2017: {10: 133.57, 20: 264.22, 30: 408.97, 40: 589.12, 50: 838.64, 60: 1062.97, 70: 1338.70, 80: 1556.04, 90: 1748.71, 100: 2915.55},
    2016: {10: 133.17, 20: 263.23, 30: 407.75, 40: 587.36, 50: 836.13, 60: 1059.09, 70: 1334.71, 80: 1551.48, 90: 1743.48, 100: 2906.83},
    2015: {10: 133.17, 20: 263.23, 30: 407.75, 40: 587.36, 50: 836.13, 60: 1059.09, 70: 1334.71, 80: 1551.48, 90: 1743.48, 100: 2858.75},
    2014: {10: 130.94, 20: 258.83, 30: 400.93, 40: 577.54, 50: 822.15, 60: 1041.39, 70: 1312.40, 80: 1525.55, 90: 1714.34, 100: 2816.00},
    2013: {10: 128.60, 20: 254.24, 30: 393.87, 40: 567.53, 50: 808.04, 60: 1023.40, 70: 1289.60, 80: 1498.89, 90: 1684.40, 100: 2769.00},
    2012: {10: 123.84, 20: 244.95, 30: 379.28, 40: 546.33, 50: 777.63, 60: 985.22, 70: 1241.87, 80: 1442.87, 90: 1621.21, 100: 2673.75},
    2011: {10: 123.84, 20: 244.95, 30: 379.28, 40: 546.33, 50: 777.63, 60: 985.22, 70: 1241.87, 80: 1442.87, 90: 1621.21, 100: 2673.75},
    2010: {10: 117.00, 20: 231.33, 30: 358.13, 40: 515.79, 50: 734.18, 60: 930.08, 70: 1172.32, 80: 1362.45, 90: 1531.20, 100: 2527.25},
    2009: {10: 110.61, 20: 218.85, 30: 338.76, 40: 487.95, 50: 694.41, 60: 879.73, 70: 1108.77, 80: 1288.22, 90: 1447.72, 100: 2389.00},
    2008: {10: 117.00, 20: 231.33, 30: 358.13, 40: 515.79, 50: 734.18, 60: 930.08, 70: 1172.32, 80: 1362.45, 90: 1531.20, 100: 2527.00},
    2007: {10: 115.00, 20: 227.40, 30: 352.00, 40: 507.00, 50: 721.00, 60: 913.00, 70: 1150.00, 80: 1336.00, 90: 1501.00, 100: 2471.00},
    2006: {10: 110.00, 20: 217.48, 30: 336.61, 40: 484.94, 50: 690.17, 60: 874.36, 70: 1101.87, 80: 1280.22, 90: 1438.46, 100: 2389.00},
    2005: {10: 107.00, 20: 211.54, 30: 327.38, 40: 471.60, 50: 671.00, 60: 849.86, 70: 1070.86, 80: 1243.86, 90: 1397.54, 100: 2316.00},
    2004: {10: 108.00, 20: 210.00, 30: 324.00, 40: 466.00, 50: 663.00, 60: 839.00, 70: 1056.00, 80: 1227.00, 90: 1380.00, 100: 2299.00},
    2003: {10: 106.00, 20: 205.00, 30: 316.00, 40: 454.00, 50: 646.00, 60: 817.00, 70: 1029.00, 80: 1195.00, 90: 1344.00, 100: 2239.00},
    2002: {10: 104.00, 20: 201.00, 30: 310.00, 40: 445.00, 50: 633.00, 60: 801.00, 70: 1008.00, 80: 1171.00, 90: 1317.00, 100: 2193.00},
    2001: {10: 103.00, 20: 199.00, 30: 306.00, 40: 439.00, 50: 625.00, 60: 790.00, 70: 995.00, 80: 1155.00, 90: 1299.00, 100: 2163.00},
    2000: {10: 101.00, 20: 194.00, 30: 298.00, 40: 427.00, 50: 609.00, 60: 769.00, 70: 969.00, 80: 1125.00, 90: 1266.00, 100: 2107.00},
};

function getDependentAddOn(rating, dependents) {
    if (dependents === 0 || rating < 30) {
        return 0;
    }
    if (rating >= 70) {
        return dependents * 153.00;
    }
    return dependents * 109.00;
}

function calculateBackpayTotal(startMonth, startYear, endMonth, endYear, rating, dependents) {
    let totalBackpay = 0;

    for (let year = startYear; year <= endYear; year++) {
        const yearRates = rates[year];
        if (!yearRates) continue;

        const baseRate = yearRates[rating];
        const dependentAddOn = getDependentAddOn(rating, dependents);
        const monthlyRate = baseRate + dependentAddOn;

        let monthsInYear = 12;
        if (year === endYear) {
            monthsInYear = endMonth;
        }
        if (year === startYear) {
            monthsInYear = monthsInYear - startMonth + 1;
        }

        totalBackpay += monthlyRate * monthsInYear;
    }

    return totalBackpay;
}

function getCurrentBackpayAmount(item) {
    const today = new Date();
    const currentMonth = today.getMonth() + 1;
    const currentYear = today.getFullYear();
    
    const endMonth = item.endMonth || currentMonth;
    const endYear = item.endYear || currentYear;
    
    return calculateBackpayTotal(item.startMonth, item.startYear, endMonth, endYear, item.rating, item.dependents);
}
