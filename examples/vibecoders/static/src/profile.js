// Setup profile page
export function setupProfilePage(user) {
  if (!user) return;
  
  // Fill in form fields with current user data
  const bioField = document.getElementById('profile-bio');
  const linkedinField = document.getElementById('profile-linkedin');
  const githubField = document.getElementById('profile-github');
  const photoField = document.getElementById('profile-photo');
  const usernameDisplay = document.getElementById('profile-username');
  const profilePhoto = document.getElementById('display-photo');
  
  if (usernameDisplay) {
    usernameDisplay.textContent = user.username;
  }
  
  if (bioField) {
    bioField.value = user.bio || '';
  }
  
  if (linkedinField) {
    linkedinField.value = user.linkedin_url || '';
  }
  
  if (githubField) {
    githubField.value = user.github_url || '';
  }
  
  if (photoField) {
    photoField.value = user.photo_url || '';
  }
  
  if (profilePhoto && user.photo_url) {
    profilePhoto.src = user.photo_url;
    profilePhoto.classList.remove('hidden');
  }
  
  // Setup form submission
  const profileForm = document.getElementById('profile-form');
  if (profileForm) {
    profileForm.addEventListener('submit', (e) => handleProfileUpdate(e, user.id));
  }
}

// Handle profile update
async function handleProfileUpdate(e, userId) {
  e.preventDefault();
  
  const bio = document.getElementById('profile-bio').value;
  const linkedinURL = document.getElementById('profile-linkedin').value;
  const githubURL = document.getElementById('profile-github').value;
  const photoURL = document.getElementById('profile-photo').value;
  const errorMsg = document.getElementById('profile-error');
  const successMsg = document.getElementById('profile-success');
  
  // Hide messages
  if (errorMsg) errorMsg.classList.add('hidden');
  if (successMsg) successMsg.classList.add('hidden');
  
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('Not authenticated');
    }
    
    const response = await fetch(`/api/users/${userId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        bio,
        linkedin_url: linkedinURL,
        github_url: githubURL,
        photo_url: photoURL,
      }),
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Update failed');
    }
    
    // Update stored user data
    localStorage.setItem('user', JSON.stringify(data));
    
    // Show success message
    if (successMsg) {
      successMsg.textContent = 'Profile updated successfully!';
      successMsg.classList.remove('hidden');
    }
    
    // Update profile photo if provided
    const profilePhoto = document.getElementById('display-photo');
    if (profilePhoto && photoURL) {
      profilePhoto.src = photoURL;
      profilePhoto.classList.remove('hidden');
    }
  } catch (error) {
    if (errorMsg) {
      errorMsg.textContent = error.message;
      errorMsg.classList.remove('hidden');
    }
  }
}