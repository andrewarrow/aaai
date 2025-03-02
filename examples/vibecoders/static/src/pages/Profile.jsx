import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';

const Profile = () => {
  const { user, updateProfile, getUserPrompts, createPrompt, updatePrompt, deletePrompt } = useAuth();
  
  // Tab state
  const [activeTab, setActiveTab] = useState('profile');
  
  // Profile form state
  const [profileFormData, setProfileFormData] = useState({
    bio: user?.bio || '',
    linked_in_url: user?.linked_in_url || '',
    github_url: user?.github_url || '',
    photo_url: user?.photo_url || '',
  });
  
  // Prompts state
  const [prompts, setPrompts] = useState([]);
  const [selectedPrompt, setSelectedPrompt] = useState(null);
  const [promptFormData, setPromptFormData] = useState({
    title: '',
    content: '',
    tags: '',
  });

  // UI state
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [promptsLoading, setPromptsLoading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  // Fetch user's prompts when component mounts
  useEffect(() => {
    if (activeTab === 'prompts') {
      fetchPrompts();
    }
  }, [activeTab]);

  const fetchPrompts = async () => {
    setPromptsLoading(true);
    try {
      const result = await getUserPrompts();
      if (result.success) {
        setPrompts(result.prompts);
      } else {
        setError(result.error);
      }
    } catch (err) {
      setError('Failed to fetch prompts');
    } finally {
      setPromptsLoading(false);
    }
  };

  // Profile form handlers
  const handleProfileChange = (e) => {
    const { name, value } = e.target;
    setProfileFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleProfileSubmit = async (e) => {
    e.preventDefault();
    
    setSuccess('');
    setError('');
    setLoading(true);
    
    try {
      const result = await updateProfile(profileFormData);
      
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

  // Prompt form handlers
  const handlePromptChange = (e) => {
    const { name, value } = e.target;
    setPromptFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const resetPromptForm = () => {
    setPromptFormData({
      title: '',
      content: '',
      tags: '',
    });
    setSelectedPrompt(null);
    setIsEditing(false);
  };

  const handlePromptSubmit = async (e) => {
    e.preventDefault();
    
    setSuccess('');
    setError('');
    setLoading(true);
    
    // Convert tags string to array
    const tagsArray = promptFormData.tags
      ? promptFormData.tags.split(',').map(tag => tag.trim())
      : [];
    
    const promptData = {
      title: promptFormData.title,
      content: promptFormData.content,
      tags: tagsArray
    };
    
    try {
      let result;
      
      if (isEditing && selectedPrompt) {
        // Update existing prompt
        result = await updatePrompt(selectedPrompt.id, promptData);
        if (result.success) {
          setSuccess('Prompt updated successfully!');
          // Update the prompt in the local state
          setPrompts(prompts.map(p => 
            p.id === selectedPrompt.id ? result.prompt : p
          ));
        }
      } else {
        // Create new prompt
        result = await createPrompt(promptData);
        if (result.success) {
          setSuccess('Prompt created successfully!');
          // Add the new prompt to the local state
          setPrompts([result.prompt, ...prompts]);
        }
      }
      
      if (!result.success) {
        setError(result.error || 'Failed to save prompt');
      } else {
        resetPromptForm();
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleEditPrompt = (prompt) => {
    setSelectedPrompt(prompt);
    setPromptFormData({
      title: prompt.title,
      content: prompt.content,
      tags: prompt.tags.join(', ')
    });
    setIsEditing(true);
  };

  const handleDeletePrompt = async (promptId) => {
    if (!window.confirm('Are you sure you want to delete this prompt?')) {
      return;
    }
    
    setLoading(true);
    try {
      const result = await deletePrompt(promptId);
      if (result.success) {
        setSuccess('Prompt deleted successfully!');
        setPrompts(prompts.filter(p => p.id !== promptId));
        if (selectedPrompt && selectedPrompt.id === promptId) {
          resetPromptForm();
        }
      } else {
        setError(result.error || 'Failed to delete prompt');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  // Function to format date
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };

  // Generate a vibrant background gradient based on prompt ID
  const getPromptBackground = (id) => {
    const backgrounds = [
      "bg-gradient-to-r from-purple-600 to-indigo-600",
      "bg-gradient-to-r from-pink-500 to-purple-600",
      "bg-gradient-to-r from-teal-400 to-blue-500",
      "bg-gradient-to-r from-orange-500 to-red-600",
      "bg-gradient-to-r from-green-500 to-teal-500"
    ];
    return backgrounds[id % backgrounds.length];
  };

  return (
    <div className="max-w-4xl mx-auto">
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
      
      {/* Tabs */}
      <div className="border-b border-gray-700 mb-6">
        <nav className="flex -mb-px">
          <button
            className={`py-4 px-6 text-center border-b-2 font-medium text-sm ${
              activeTab === 'profile'
                ? 'border-purple-500 text-purple-500'
                : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-400'
            }`}
            onClick={() => setActiveTab('profile')}
          >
            Basic Info
          </button>
          <button
            className={`py-4 px-6 text-center border-b-2 font-medium text-sm ${
              activeTab === 'prompts'
                ? 'border-purple-500 text-purple-500'
                : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-400'
            }`}
            onClick={() => setActiveTab('prompts')}
          >
            My Prompts
          </button>
        </nav>
      </div>
      
      {/* Success and Error Messages */}
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
      
      {/* Tab Content */}
      {activeTab === 'profile' && (
        <div className="bg-gray-800 rounded-lg shadow-lg p-6">
          <h2 className="text-2xl font-bold text-purple-500 mb-6">Edit Profile</h2>
          
          <form onSubmit={handleProfileSubmit}>
            <div className="mb-4">
              <label htmlFor="bio" className="block text-gray-300 mb-2">
                Bio
              </label>
              <textarea
                id="bio"
                name="bio"
                value={profileFormData.bio}
                onChange={handleProfileChange}
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
                value={profileFormData.linked_in_url}
                onChange={handleProfileChange}
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
                value={profileFormData.github_url}
                onChange={handleProfileChange}
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
                value={profileFormData.photo_url}
                onChange={handleProfileChange}
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
      )}
      
      {activeTab === 'prompts' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Prompt Form */}
          <div className="bg-gray-800 rounded-lg shadow-lg p-6">
            <h2 className="text-2xl font-bold text-purple-500 mb-6">
              {isEditing ? 'Edit Prompt' : 'Add New Prompt'}
            </h2>
            
            <form onSubmit={handlePromptSubmit}>
              <div className="mb-4">
                <label htmlFor="title" className="block text-gray-300 mb-2">
                  Title
                </label>
                <input
                  id="title"
                  name="title"
                  value={promptFormData.title}
                  onChange={handlePromptChange}
                  className="input"
                  placeholder="Prompt title"
                  required
                />
              </div>
              
              <div className="mb-4">
                <label htmlFor="content" className="block text-gray-300 mb-2">
                  Content
                </label>
                <textarea
                  id="content"
                  name="content"
                  value={promptFormData.content}
                  onChange={handlePromptChange}
                  className="input"
                  placeholder="Enter your prompt content..."
                  rows="6"
                  required
                />
              </div>
              
              <div className="mb-6">
                <label htmlFor="tags" className="block text-gray-300 mb-2">
                  Tags (comma-separated)
                </label>
                <input
                  id="tags"
                  name="tags"
                  value={promptFormData.tags}
                  onChange={handlePromptChange}
                  className="input"
                  placeholder="react, javascript, tutorial"
                />
              </div>
              
              <div className="flex gap-3">
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={loading}
                >
                  {loading ? 'Saving...' : isEditing ? 'Update Prompt' : 'Create Prompt'}
                </button>
                
                {isEditing && (
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={resetPromptForm}
                  >
                    Cancel
                  </button>
                )}
              </div>
            </form>
          </div>
          
          {/* Prompts List */}
          <div>
            <h2 className="text-2xl font-bold text-purple-500 mb-6">
              My Prompts
            </h2>
            
            {promptsLoading ? (
              <div className="flex justify-center items-center h-40">
                <div className="w-12 h-12 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
              </div>
            ) : prompts.length > 0 ? (
              <div className="space-y-4 max-h-[600px] overflow-y-auto pr-2">
                {prompts.map(prompt => (
                  <div 
                    key={prompt.id} 
                    className={`rounded-lg overflow-hidden shadow-lg ${getPromptBackground(prompt.id)}`}
                  >
                    <div className="p-4">
                      <h3 className="text-lg font-bold text-white">{prompt.title}</h3>
                      <p className="text-gray-100 text-sm line-clamp-2 mb-3">{prompt.content}</p>
                      
                      <div className="flex flex-wrap gap-1 mb-2">
                        {prompt.tags.map(tag => (
                          <span 
                            key={tag} 
                            className="bg-black bg-opacity-30 text-white px-2 py-1 rounded-full text-xs"
                          >
                            #{tag}
                          </span>
                        ))}
                      </div>
                      
                      <div className="flex justify-between items-center">
                        <span className="text-gray-200 text-xs">{formatDate(prompt.created_at)}</span>
                        
                        <div className="flex gap-2">
                          <button 
                            onClick={() => handleEditPrompt(prompt)}
                            className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-xs"
                          >
                            Edit
                          </button>
                          <button 
                            onClick={() => handleDeletePrompt(prompt.id)}
                            className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-xs"
                          >
                            Delete
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="bg-gray-700 p-6 rounded-lg text-center">
                <p className="text-gray-300 mb-4">You don't have any prompts yet.</p>
                <p className="text-gray-400">Create your first prompt using the form!</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default Profile;