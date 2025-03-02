import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';

const Profile = () => {
  const { user, updateProfile } = useAuth();
  
  const [formData, setFormData] = useState({
    bio: user?.bio || '',
    linked_in_url: user?.linked_in_url || '',
    github_url: user?.github_url || '',
    photo_url: user?.photo_url || '',
  });
  
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    setSuccess('');
    setError('');
    setLoading(true);
    
    try {
      const result = await updateProfile(formData);
      
      if (result.success) {
        setSuccess('Profile updated successfully!');
      } else {
        setError(result.error || 'Failed to update profile');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="flex flex-col md:flex-row items-center md:items-start mb-8 gap-8">
        <div className="w-48 h-48 overflow-hidden rounded-full flex-shrink-0">
          <img
            src={user?.photo_url || 'https://via.placeholder.com/200'}
            alt={user?.username}
            className="w-full h-full object-cover"
          />
        </div>
        
        <div>
          <h1 className="text-3xl font-bold text-purple-500 mb-2">{user?.username}</h1>
          {user?.bio && <p className="text-gray-300 mb-4">{user.bio}</p>}
          
          <div className="flex flex-wrap gap-4">
            {user?.github_url && (
              <a
                href={user.github_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-purple-400 hover:text-purple-300 flex items-center"
              >
                GitHub
              </a>
            )}
            
            {user?.linked_in_url && (
              <a
                href={user.linked_in_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-purple-400 hover:text-purple-300 flex items-center"
              >
                LinkedIn
              </a>
            )}
          </div>
        </div>
      </div>
      
      <div className="bg-gray-800 rounded-lg shadow-lg p-6">
        <h2 className="text-2xl font-bold text-purple-500 mb-6">Edit Profile</h2>
        
        {success && (
          <div className="bg-green-500 text-white p-3 rounded-md mb-4">
            {success}
          </div>
        )}
        
        {error && (
          <div className="bg-red-500 text-white p-3 rounded-md mb-4">
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="bio" className="block text-gray-300 mb-2">
              Bio
            </label>
            <textarea
              id="bio"
              name="bio"
              value={formData.bio}
              onChange={handleChange}
              className="input"
              placeholder="Tell us about yourself"
              rows="3"
            />
          </div>
          
          <div className="mb-4">
            <label htmlFor="linked_in_url" className="block text-gray-300 mb-2">
              LinkedIn URL
            </label>
            <input
              type="url"
              id="linked_in_url"
              name="linked_in_url"
              value={formData.linked_in_url}
              onChange={handleChange}
              className="input"
              placeholder="https://www.linkedin.com/in/yourprofile"
            />
          </div>
          
          <div className="mb-4">
            <label htmlFor="github_url" className="block text-gray-300 mb-2">
              GitHub URL
            </label>
            <input
              type="url"
              id="github_url"
              name="github_url"
              value={formData.github_url}
              onChange={handleChange}
              className="input"
              placeholder="https://github.com/yourusername"
            />
          </div>
          
          <div className="mb-6">
            <label htmlFor="photo_url" className="block text-gray-300 mb-2">
              Photo URL
            </label>
            <input
              type="url"
              id="photo_url"
              name="photo_url"
              value={formData.photo_url}
              onChange={handleChange}
              className="input"
              placeholder="https://example.com/your-photo.jpg"
            />
          </div>
          
          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading}
          >
            {loading ? 'Updating...' : 'Update Profile'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default Profile;