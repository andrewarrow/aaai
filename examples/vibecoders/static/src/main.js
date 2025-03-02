import { showLogin, showRegister } from './auth.js';
import { setupProfilePage } from './profile.js';

// Initialize app based on current page
document.addEventListener('DOMContentLoaded', () => {
  // Get current page path
  const path = window.location.pathname;

  // Check for stored token
  const token = localStorage.getItem('token');
  const userJson = localStorage.getItem('user');
  
  // If we have a token and user data, we're logged in
  const isLoggedIn = token && userJson;
  
  // If on root page
  if (path === '/' || path === '/index.html') {
    setupHomePage(isLoggedIn);
  }
  
  // If on profile page and logged in
  if (path === '/profile.html' && isLoggedIn) {
    setupProfilePage(JSON.parse(userJson));
  } else if (path === '/profile.html' && !isLoggedIn) {
    // Redirect to home if trying to access profile while not logged in
    window.location.href = '/';
  }
});

// Setup home page
function setupHomePage(isLoggedIn) {
  const loginBtn = document.getElementById('login-btn');
  const registerBtn = document.getElementById('register-btn');
  const dashboardBtn = document.getElementById('dashboard-btn');
  const logoutBtn = document.getElementById('logout-btn');
  
  if (isLoggedIn) {
    // Show dashboard and logout buttons
    if (loginBtn) loginBtn.classList.add('hidden');
    if (registerBtn) registerBtn.classList.add('hidden');
    if (dashboardBtn) dashboardBtn.classList.remove('hidden');
    if (logoutBtn) {
      logoutBtn.classList.remove('hidden');
      logoutBtn.addEventListener('click', logout);
    }
  } else {
    // Show login and register buttons
    if (loginBtn) {
      loginBtn.classList.remove('hidden');
      loginBtn.addEventListener('click', showLogin);
    }
    if (registerBtn) {
      registerBtn.classList.remove('hidden');
      registerBtn.addEventListener('click', showRegister);
    }
    if (dashboardBtn) dashboardBtn.classList.add('hidden');
    if (logoutBtn) logoutBtn.classList.add('hidden');
  }
}

// Logout function
function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  window.location.reload();
}