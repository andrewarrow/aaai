document.addEventListener('DOMContentLoaded', function() {
    // DOM Elements
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const loginFormContainer = document.getElementById('login-form-container');
    const registerFormContainer = document.getElementById('register-form-container');
    const registerToggle = document.getElementById('register-toggle');
    const loginToggle = document.getElementById('login-toggle');
    const logoutButton = document.getElementById('logout-button');
    const userLoggedIn = document.getElementById('user-logged-in');
    const loggedUsername = document.getElementById('logged-username');
    const messageContainer = document.getElementById('message-container');
    const taskContainer = document.getElementById('task-container');
    const taskInfo = document.getElementById('task-info');
    const welcomeMessage = document.getElementById('welcome-message');

    // Check if user is already logged in (from localStorage)
    checkLoginStatus();

    // Toggle between login and register forms
    registerToggle.addEventListener('click', function() {
        loginFormContainer.style.display = 'none';
        registerFormContainer.style.display = 'block';
    });

    loginToggle.addEventListener('click', function() {
        registerFormContainer.style.display = 'none';
        loginFormContainer.style.display = 'block';
    });

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
                // Redirect to dashboard instead of updating UI
                window.location.href = '/dashboard';
                // No need to call updateUIAfterLogin here as we're redirecting
                // updateUIAfterLogin(data.user.username);
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
                registerFormContainer.style.display = 'none';
                loginFormContainer.style.display = 'block';
            }
        })
        .catch(error => {
            showMessage('An error occurred during registration: ' + error, 'error');
        });
    });

    // Handle logout
    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('user');
        showMessage('Logged out successfully', 'success');
        // Redirect to homepage after logout
        window.location.href = '/';
        // updateUIAfterLogout(); // No need as we're redirecting
    });

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
        loginFormContainer.style.display = 'none';
        registerFormContainer.style.display = 'none';
        userLoggedIn.style.display = 'block';
        loggedUsername.textContent = username;
        welcomeMessage.style.display = 'none';
        taskContainer.style.display = 'block';
        taskInfo.style.display = 'block';
    }

    // Update UI after logout
    function updateUIAfterLogout() {
        userLoggedIn.style.display = 'none';
        loginFormContainer.style.display = 'block';
        welcomeMessage.style.display = 'block';
        taskContainer.style.display = 'none';
        taskInfo.style.display = 'none';
    }

    // Check if user is already logged in
    function checkLoginStatus() {
        const user = JSON.parse(localStorage.getItem('user'));
        
        // Check if we're on the index page and logged in
        if (user && window.location.pathname === '/') {
            // Redirect to dashboard
            window.location.href = '/dashboard';
            return;
        }
        
        // Check if we're on the dashboard page and not logged in
        if (!user && window.location.pathname === '/dashboard') {
            // Redirect to homepage
            window.location.href = '/';
            return;
        }
        if (user && window.location.pathname === '/dashboard') {
            updateUIAfterLogin(user.username);
        }
    }
});
