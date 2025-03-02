// Show login modal
export function showLogin() {
  const modal = document.getElementById('auth-modal');
  const loginForm = document.getElementById('login-form');
  const registerForm = document.getElementById('register-form');
  
  if (modal && loginForm && registerForm) {
    modal.classList.remove('hidden');
    loginForm.classList.remove('hidden');
    registerForm.classList.add('hidden');
    
    // Focus on username input
    document.getElementById('login-username').focus();
  }
  
  setupAuthListeners();
}

// Show register modal
export function showRegister() {
  const modal = document.getElementById('auth-modal');
  const loginForm = document.getElementById('login-form');
  const registerForm = document.getElementById('register-form');
  
  if (modal && loginForm && registerForm) {
    modal.classList.remove('hidden');
    registerForm.classList.remove('hidden');
    loginForm.classList.add('hidden');
    
    // Focus on username input
    document.getElementById('register-username').focus();
  }
  
  setupAuthListeners();
}

// Close auth modal
function closeAuthModal() {
  const modal = document.getElementById('auth-modal');
  if (modal) {
    modal.classList.add('hidden');
  }
}

// Setup auth event listeners
function setupAuthListeners() {
  // Close modal when clicking outside
  const modal = document.getElementById('auth-modal');
  const modalBg = document.getElementById('modal-bg');
  const closeBtn = document.getElementById('close-modal');
  
  if (modalBg) {
    modalBg.addEventListener('click', closeAuthModal);
  }
  
  if (closeBtn) {
    closeBtn.addEventListener('click', closeAuthModal);
  }
  
  // Prevent clicks inside modal from closing it
  const modalContent = document.getElementById('modal-content');
  if (modalContent) {
    modalContent.addEventListener('click', (e) => {
      e.stopPropagation();
    });
  }
  
  // Toggle between login and register
  const switchToRegister = document.getElementById('switch-to-register');
  const switchToLogin = document.getElementById('switch-to-login');
  
  if (switchToRegister) {
    switchToRegister.addEventListener('click', (e) => {
      e.preventDefault();
      showRegister();
    });
  }
  
  if (switchToLogin) {
    switchToLogin.addEventListener('click', (e) => {
      e.preventDefault();
      showLogin();
    });
  }
  
  // Handle form submissions
  const loginForm = document.getElementById('login-form');
  const registerForm = document.getElementById('register-form');
  
  if (loginForm) {
    loginForm.addEventListener('submit', handleLogin);
  }
  
  if (registerForm) {
    registerForm.addEventListener('submit', handleRegister);
  }
}

// Handle login submission
async function handleLogin(e) {
  e.preventDefault();
  
  const username = document.getElementById('login-username').value;
  const password = document.getElementById('login-password').value;
  const errorMsg = document.getElementById('login-error');
  
  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Login failed');
    }
    
    // Store token and user data
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    
    // Close modal and refresh page
    closeAuthModal();
    window.location.reload();
  } catch (error) {
    if (errorMsg) {
      errorMsg.textContent = error.message;
      errorMsg.classList.remove('hidden');
    }
  }
}

// Handle register submission
async function handleRegister(e) {
  e.preventDefault();
  
  const username = document.getElementById('register-username').value;
  const password = document.getElementById('register-password').value;
  const bio = document.getElementById('register-bio').value;
  const linkedinURL = document.getElementById('register-linkedin').value;
  const githubURL = document.getElementById('register-github').value;
  const photoURL = document.getElementById('register-photo').value;
  const errorMsg = document.getElementById('register-error');
  
  try {
    const response = await fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username,
        password,
        bio,
        linkedin_url: linkedinURL,
        github_url: githubURL,
        photo_url: photoURL,
      }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Registration failed');
    }
    
    // Store token and user data
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    
    // Close modal and refresh page
    closeAuthModal();
    window.location.reload();
  } catch (error) {
    if (errorMsg) {
      errorMsg.textContent = error.message;
      errorMsg.classList.remove('hidden');
    }
  }
}