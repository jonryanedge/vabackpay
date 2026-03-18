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
