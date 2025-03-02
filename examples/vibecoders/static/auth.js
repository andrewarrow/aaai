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
    
    // Check if user is already logged in (via server session API)
    checkLoginStatus();

    // Handle login form submission
    if (loginForm) {
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
                }),
                credentials: 'include' // Include cookies in the request
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    showMessage(data.error, 'error');
                } else {
                    showMessage('Login successful!', 'success');
                    // Redirect to dashboard after successful login
                    window.location.href = '/';
                }
            })
            .catch(error => {
                showMessage('An error occurred during login: ' + error, 'error');
            });
        });
    }

    // Handle register form submission
    if (registerForm) {
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
    }
    
    // Logout is now handled via form POST to /logout

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

    // Check if user is already logged in using session API
    function checkLoginStatus() {
        fetch('/api/user', {
            method: 'GET',
            credentials: 'include' // Include cookies in the request
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Not authenticated');
            }
            return response.json();
        })
        .then(data => {
            if (data.authenticated) {
                updateUIAfterLogin(data.username);
                
                // If user is logged in and tries to access login/register pages, redirect to home
                const currentPath = window.location.pathname;
                if (currentPath === '/login' || currentPath === '/register') {
                    window.location.href = '/';
                }
            } else {
                // User is not authenticated
                updateUIAfterLogout();
            }
        })
        .catch(error => {
            // User is not authenticated, update UI
            console.log('Not authenticated:', error);
            updateUIAfterLogout();
        });
    }
});
