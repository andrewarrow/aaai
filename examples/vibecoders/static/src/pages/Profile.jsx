import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';

const Profile = () => {
  const { 
    user, updateProfile, 
    getUserPrompts, createPrompt, updatePrompt, deletePrompt,
    getUserProjects, createProject, updateProject, deleteProject 
  } = useAuth();
  
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
  
  // Projects state
  const [projects, setProjects] = useState([]);
  const [selectedProject, setSelectedProject] = useState(null);
  const [projectFormData, setProjectFormData] = useState({
    title: '',
    description: '',
    github_url: '',
    website_url: '',
    image_url1: '',
    image_url2: '',
    image_url3: '',
  });

  // UI state
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [promptsLoading, setPromptsLoading] = useState(false);
  const [projectsLoading, setProjectsLoading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  // Fetch user's data when tab changes
  useEffect(() => {
    if (activeTab === 'prompts') {
      fetchPrompts();
    } else if (activeTab === 'projects') {
      fetchProjects();
    }
  }, [activeTab]);
  
  const fetchProjects = async () => {
    setProjectsLoading(true);
    try {
      const result = await getUserProjects();
      if (result.success) {
        setProjects(result.projects);
      } else {
        setError(result.error);
      }
    } catch (err) {
      setError('Failed to fetch projects');
    } finally {
      setProjectsLoading(false);
    }
  };

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
  
  // Project form handlers
  const handleProjectChange = (e) => {
    const { name, value } = e.target;
    setProjectFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const resetProjectForm = () => {
    setProjectFormData({
      title: '',
      description: '',
      github_url: '',
      website_url: '',
      image_url1: '',
      image_url2: '',
      image_url3: '',
    });
    setSelectedProject(null);
    setIsEditing(false);
  };

  const handleProjectSubmit = async (e) => {
    e.preventDefault();
    
    setSuccess('');
    setError('');
    setLoading(true);
    
    try {
      let result;
      
      if (isEditing && selectedProject) {
        // Update existing project
        result = await updateProject(selectedProject.id, projectFormData);
        if (result.success) {
          setSuccess('Project updated successfully!');
          // Update the project in the local state
          setProjects(projects.map(p => 
            p.id === selectedProject.id ? result.project : p
          ));
        }
      } else {
        // Create new project
        result = await createProject(projectFormData);
        if (result.success) {
          setSuccess('Project created successfully!');
          // Add the new project to the local state
          setProjects([result.project, ...projects]);
        }
      }
      
      if (!result.success) {
        setError(result.error || 'Failed to save project');
      } else {
        resetProjectForm();
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleEditProject = (project) => {
    setSelectedProject(project);
    setProjectFormData({
      title: project.title,
      description: project.description,
      github_url: project.github_url || '',
      website_url: project.website_url || '',
      image_url1: project.image_url1 || '',
      image_url2: project.image_url2 || '',
      image_url3: project.image_url3 || '',
    });
    setIsEditing(true);
  };

  const handleDeleteProject = async (projectId) => {
    if (!window.confirm('Are you sure you want to delete this project?')) {
      return;
    }
    
    setLoading(true);
    try {
      const result = await deleteProject(projectId);
      if (result.success) {
        setSuccess('Project deleted successfully!');
        setProjects(projects.filter(p => p.id !== projectId));
        if (selectedProject && selectedProject.id === projectId) {
          resetProjectForm();
        }
      } else {
        setError(result.error || 'Failed to delete project');
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

  const getPromptBackground = (id) => {
    const backgrounds = [
      "bg-purple-800",
      "bg-indigo-800",
      "bg-blue-800",
      "bg-teal-800",
      "bg-emerald-800"
    ];
    return backgrounds[id % backgrounds.length];
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex flex-col md:flex-row items-center md:items-start mb-8 gap-8 hover:bg-gray-800/30 p-4 rounded-lg transition-colors cursor-pointer" onClick={() => window.location.href = `/users/${user?.username}`}>
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
                onClick={(e) => e.stopPropagation()} // Prevent triggering parent onClick
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
                onClick={(e) => e.stopPropagation()} // Prevent triggering parent onClick
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
          <button
            className={`py-4 px-6 text-center border-b-2 font-medium text-sm ${
              activeTab === 'projects'
                ? 'border-purple-500 text-purple-500'
                : 'border-transparent text-gray-400 hover:text-gray-300 hover:border-gray-400'
            }`}
            onClick={() => setActiveTab('projects')}
          >
            Projects
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
      
      {activeTab === 'projects' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Project Form */}
          <div className="bg-gray-800 rounded-lg shadow-lg p-6">
            <h2 className="text-2xl font-bold text-purple-500 mb-6">
              {isEditing ? 'Edit Project' : 'Add New Project'}
            </h2>
            
            <form onSubmit={handleProjectSubmit}>
              <div className="mb-4">
                <label htmlFor="project-title" className="block text-gray-300 mb-2">
                  Title
                </label>
                <input
                  id="project-title"
                  name="title"
                  value={projectFormData.title}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="Project title"
                  required
                />
              </div>
              
              <div className="mb-4">
                <label htmlFor="description" className="block text-gray-300 mb-2">
                  Description
                </label>
                <textarea
                  id="description"
                  name="description"
                  value={projectFormData.description}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="Describe your project..."
                  rows="5"
                  required
                />
              </div>
              
              <div className="mb-4">
                <label htmlFor="github_url" className="block text-gray-300 mb-2">
                  GitHub Repository URL
                </label>
                <input
                  type="url"
                  id="github_url"
                  name="github_url"
                  value={projectFormData.github_url}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="https://github.com/username/repo"
                />
              </div>
              
              <div className="mb-4">
                <label htmlFor="website_url" className="block text-gray-300 mb-2">
                  Website URL
                </label>
                <input
                  type="url"
                  id="website_url"
                  name="website_url"
                  value={projectFormData.website_url}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="https://your-project.com"
                />
              </div>
              
              <div className="mb-4">
                <label htmlFor="image_url1" className="block text-gray-300 mb-2">
                  Image URL 1
                </label>
                <input
                  type="url"
                  id="image_url1"
                  name="image_url1"
                  value={projectFormData.image_url1}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="https://example.com/image1.jpg"
                />
              </div>
              
              <div className="mb-4">
                <label htmlFor="image_url2" className="block text-gray-300 mb-2">
                  Image URL 2 (Optional)
                </label>
                <input
                  type="url"
                  id="image_url2"
                  name="image_url2"
                  value={projectFormData.image_url2}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="https://example.com/image2.jpg"
                />
              </div>
              
              <div className="mb-6">
                <label htmlFor="image_url3" className="block text-gray-300 mb-2">
                  Image URL 3 (Optional)
                </label>
                <input
                  type="url"
                  id="image_url3"
                  name="image_url3"
                  value={projectFormData.image_url3}
                  onChange={handleProjectChange}
                  className="input"
                  placeholder="https://example.com/image3.jpg"
                />
              </div>
              
              <div className="flex gap-3">
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={loading}
                >
                  {loading ? 'Saving...' : isEditing ? 'Update Project' : 'Create Project'}
                </button>
                
                {isEditing && (
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={resetProjectForm}
                  >
                    Cancel
                  </button>
                )}
              </div>
            </form>
          </div>
          
          {/* Projects List */}
          <div>
            <h2 className="text-2xl font-bold text-purple-500 mb-6">
              My Projects
            </h2>
            
            {projectsLoading ? (
              <div className="flex justify-center items-center h-40">
                <div className="w-12 h-12 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
              </div>
            ) : projects.length > 0 ? (
              <div className="space-y-6 max-h-[700px] overflow-y-auto pr-2">
                {projects.map(project => (
                  <div 
                    key={project.id} 
                    className="bg-gray-800 rounded-lg overflow-hidden shadow-lg border border-gray-700"
                  >
                    {project.image_url1 && (
                      <div className="h-40 overflow-hidden">
                        <img 
                          src={project.image_url1} 
                          alt={project.title}
                          className="w-full h-full object-cover"
                        />
                      </div>
                    )}
                    
                    <div className="p-4">
                      <h3 className="text-lg font-bold text-white mb-2">{project.title}</h3>
                      <p className="text-gray-300 text-sm mb-3 line-clamp-2">{project.description}</p>
                      
                      <div className="flex flex-wrap gap-2 mb-3">
                        {project.github_url && (
                          <a 
                            href={project.github_url} 
                            target="_blank" 
                            rel="noopener noreferrer" 
                            className="text-blue-400 hover:text-blue-300 text-sm flex items-center gap-1"
                          >
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                            </svg>
                            GitHub
                          </a>
                        )}
                        
                        {project.website_url && (
                          <a 
                            href={project.website_url} 
                            target="_blank" 
                            rel="noopener noreferrer" 
                            className="text-blue-400 hover:text-blue-300 text-sm flex items-center gap-1"
                          >
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                            </svg>
                            Website
                          </a>
                        )}
                      </div>
                      
                      {/* Image thumbnails */}
                      {(project.image_url2 || project.image_url3) && (
                        <div className="flex gap-2 mb-3">
                          {project.image_url2 && (
                            <div className="w-16 h-16 rounded overflow-hidden">
                              <img 
                                src={project.image_url2} 
                                alt={`${project.title} thumbnail 2`}
                                className="w-full h-full object-cover"
                              />
                            </div>
                          )}
                          
                          {project.image_url3 && (
                            <div className="w-16 h-16 rounded overflow-hidden">
                              <img 
                                src={project.image_url3} 
                                alt={`${project.title} thumbnail 3`}
                                className="w-full h-full object-cover"
                              />
                            </div>
                          )}
                        </div>
                      )}
                      
                      <div className="flex justify-between items-center">
                        <span className="text-gray-400 text-xs">{formatDate(project.created_at)}</span>
                        
                        <div className="flex gap-2">
                          <button 
                            onClick={() => handleEditProject(project)}
                            className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-xs"
                          >
                            Edit
                          </button>
                          <button 
                            onClick={() => handleDeleteProject(project.id)}
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
                <p className="text-gray-300 mb-4">You don't have any projects yet.</p>
                <p className="text-gray-400">Create your first project using the form!</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default Profile;
