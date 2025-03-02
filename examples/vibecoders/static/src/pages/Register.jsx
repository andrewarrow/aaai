import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Register = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirm_password: '',
    bio: '',
    linked_in_url: '',
    github_url: '',
    photo_url: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.username || !formData.password || !formData.confirm_password) {
      setError('Username, password, and confirm password are required');
      return;
    }
    
    if (formData.password !== formData.confirm_password) {
      setError('Passwords do not match');
      return;
    }
    
    setError('');
    setLoading(true);
    
    try {
      const result = await register(formData);
      
      if (result.success) {
        navigate('/login');
      } else {
        setError(result.error || 'Failed to register');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold text-center text-purple-500 mb-6">Register</h1>
      
      {error && (
        <div className="bg-red-500 text-white p-3 rounded-md mb-4">
          {error}
        </div>
      )}
      
      <form onSubmit={handleSubmit} className="bg-gray-800 shadow-md rounded-lg p-6">
        <div className="mb-4">
          <label htmlFor="username" className="block text-gray-300 mb-2">
            Username *
          </label>
          <input
            type="text"
            id="username"
            name="username"
            value={formData.username}
            onChange={handleChange}
            className="input"
            placeholder="Enter your username"
            required
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="password" className="block text-gray-300 mb-2">
            Password *
          </label>
          <input
            type="password"
            id="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            className="input"
            placeholder="Enter your password"
            required
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="confirm_password" className="block text-gray-300 mb-2">
            Confirm Password *
          </label>
          <input
            type="password"
            id="confirm_password"
            name="confirm_password"
            value={formData.confirm_password}
            onChange={handleChange}
            className="input"
            placeholder="Confirm your password"
            required
          />
        </div>
        
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
        
        <div className="flex justify-between items-center">
          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading}
          >
            {loading ? 'Registering...' : 'Register'}
          </button>
          
          <Link to="/login" className="text-purple-400 hover:text-purple-300">
            Already have an account?
          </Link>
        </div>
      </form>
    </div>
  );
};

export default Register;