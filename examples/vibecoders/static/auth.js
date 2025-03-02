document.addEventListener('DOMContentLoaded', function() {
    // DOM Elements
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const logoutButton = document.getElementById('logout-button');
    const userLoggedIn = document.getElementById('user-logged-in');
    const loggedUsername = document.getElementById('logged-username');
    const messageContainer = document.getElementById('message-container');
    const taskContainer = document.getElementById('task-container');
    const taskInfo = document.getElementById('task-info');
    const splashContent = document.getElementById('splash-content');
    const navbarLoggedOut = document.getElementById('navbar-logged-out');
    const navbarLoggedIn = document.getElementById('navbar-logged-in');
    const loginButton = document.getElementById('login-button-nav');
    const registerButton = document.getElementById('register-button-nav');
    
    // Check if user is already logged in (from localStorage)
    checkLoginStatus();

    // Handle login form submission
    loginForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        
        // Send login request
        fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                showMessage(data.error, 'error');
            } else {
                // Store user info in localStorage
                localStorage.setItem('user', JSON.stringify(data.user));
                showMessage('Login successful!', 'success');
                // Redirect to dashboard after successful login
                window.location.href = '/';
            }
        })
        .catch(error => {
            showMessage('An error occurred during login: ' + error, 'error');
        });
    });

    // Handle register form submission
    registerForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const username = document.getElementById('reg-username').value;
        const password = document.getElementById('reg-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        
        if (password !== confirmPassword) {
            showMessage('Passwords do not match', 'error');
            return;
        }
        
        // Send registration request
        fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                showMessage(data.error, 'error');
            } else {
                showMessage('Registration successful! Please log in.', 'success');
                // Redirect to login page after successful registration
                window.location.href = '/login';
            }
        })
        .catch(error => {
            showMessage('An error occurred during registration: ' + error, 'error');
        });
    });

    // Handle logout
    if (logoutButton) {
        logoutButton.addEventListener('click', function() {
            localStorage.removeItem('user');
            updateUIAfterLogout();
            showMessage('Logged out successfully', 'success');
        });
    }

    // Helper functions to display messages

    // Helper function to display messages
    function showMessage(message, type) {
        messageContainer.textContent = message;
        messageContainer.className = type;
        // Clear message after 5 seconds
        setTimeout(() => {
            messageContainer.textContent = '';
            messageContainer.className = '';
        }, 5000);
    }

    // Update UI after successful login
    function updateUIAfterLogin(username) {
        if (userLoggedIn) userLoggedIn.style.display = 'block';
        if (navbarLoggedIn) navbarLoggedIn.style.display = 'block';
        if (navbarLoggedOut) navbarLoggedOut.style.display = 'none';
        if (loggedUsername) loggedUsername.textContent = username;
        
        if (splashContent) {
            splashContent.style.display = 'none';
        }
        if (taskContainer) taskContainer.style.display = 'block';
        if (taskInfo) taskInfo.style.display = 'block';
    }

    // Update UI after logout
    function updateUIAfterLogout() {
        if (navbarLoggedIn) navbarLoggedIn.style.display = 'none';
        if (navbarLoggedOut) navbarLoggedOut.style.display = 'block';
        if (userLoggedIn) userLoggedIn.style.display = 'none';
        if (splashContent) splashContent.style.display = 'block';
        
        // Redirect to home page if we're on a protected page
        const currentPath = window.location.pathname;
        if (currentPath !== '/' && currentPath !== '/login' && currentPath !== '/register') {
            window.location.href = '/';
        }
    }

    // Check if user is already logged in
    function checkLoginStatus() {
        const user = JSON.parse(localStorage.getItem('user'));
        const currentPath = window.location.pathname;
        
        if (user) {
            updateUIAfterLogin(user.username);
            
            // If user is logged in and tries to access login/register pages, redirect to home
            if (currentPath === '/login' || currentPath === '/register') {
                window.location.href = '/';
            }
        } else {
            // If user is not logged in and tries to access protected pages, redirect to login
            // No protected pages yet, but can be added when needed
        }
    }
});
