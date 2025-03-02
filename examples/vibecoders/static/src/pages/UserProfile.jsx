import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const UserProfile = () => {
  const { username } = useParams();
  const [user, setUser] = useState(null);
  const [prompts, setPrompts] = useState([]);
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch user data
        const userResponse = await fetch(`/api/users/${username}`);
        
        if (!userResponse.ok) {
          if (userResponse.status === 404) {
            throw new Error('User not found');
          }
          throw new Error('Failed to fetch user data');
        }
        
        const userData = await userResponse.json();
        setUser(userData);
        
        // Fetch user's public prompts
        const promptsResponse = await fetch(`/api/users/${username}/prompts`);
        
        if (!promptsResponse.ok) {
          throw new Error('Failed to fetch prompts');
        }
        
        const promptsData = await promptsResponse.json();
        setPrompts(promptsData);
        
        // Fetch user's public projects
        const projectsResponse = await fetch(`/api/users/${username}/projects`);
        
        if (!projectsResponse.ok) {
          throw new Error('Failed to fetch projects');
        }
        
        const projectsData = await projectsResponse.json();
        setProjects(projectsData);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [username]);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[50vh]">
        <div className="w-12 h-12 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-500 text-white p-4 rounded-md text-center max-w-md mx-auto">
        {error}
      </div>
    );
  }

  if (!user) {
    return null;
  }

  // Function to format date
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
      return "Today";
    } else if (diffDays === 1) {
      return "Yesterday";
    } else if (diffDays < 7) {
      return `${diffDays} days ago`;
    } else if (diffDays < 30) {
      const weeks = Math.floor(diffDays / 7);
      return `${weeks} ${weeks === 1 ? 'week' : 'weeks'} ago`;
    } else {
      return date.toLocaleDateString();
    }
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
      <div className="flex flex-col md:flex-row items-center md:items-start mb-12 gap-8">
        <div className="w-48 h-48 overflow-hidden rounded-full flex-shrink-0">
          <img
            src={user.photo_url || 'https://via.placeholder.com/200'}
            alt={user.username}
            className="w-full h-full object-cover"
          />
        </div>
        
        <div>
          <h1 className="text-3xl font-bold text-purple-500 mb-2">{user.username}</h1>
          {user.fullname && <p className="text-xl text-gray-200 mb-2">{user.fullname}</p>}
          {user.bio && <p className="text-gray-300 mb-4">{user.bio}</p>}
          
          <div className="flex flex-wrap gap-4">
            {user.github_url && (
              <a
                href={user.github_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-purple-400 hover:text-purple-300 flex items-center"
              >
                GitHub
              </a>
            )}
            
            {user.linked_in_url && (
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

      
      {/* Projects Section */}
      {projects.length > 0 && (
        <div className="mb-10">
          <h2 className="text-2xl font-bold text-purple-500 mb-6">Projects</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {projects.map(project => (
              <div 
                key={project.id} 
                className="bg-gray-800 rounded-lg overflow-hidden shadow-lg border border-gray-700"
              >
                {project.image_url1 && (
                  <div className="h-48 overflow-hidden">
                    <img 
                      src={project.image_url1} 
                      alt={project.title}
                      className="w-full h-full object-cover"
                    />
                  </div>
                )}
                
                <div className="p-4">
                  <h3 className="text-lg font-bold text-white mb-2">{project.title}</h3>
                  <p className="text-gray-300 text-sm mb-3 line-clamp-3">{project.description}</p>
                  
                  <div className="flex flex-wrap gap-2 mt-4">
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
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
      {/* Latest Prompts Section */}
      {/* Latest Prompts Section */}
      <div className="mb-10">
        <h2 className="text-2xl font-bold text-purple-500 mb-6">Latest Prompts</h2>
        
        {prompts.length > 0 ? (
          <div className="space-y-6">
            {prompts.map(prompt => (
              <div 
                key={prompt.id} 
                className={`rounded-xl overflow-hidden shadow-lg transform transition-transform hover:scale-102 ${getPromptBackground(prompt.id)}`}
              >
                <div className="p-6">
                  <h3 className="text-xl font-bold text-white mb-3">{prompt.title}</h3>
                  <p className="text-gray-100 mb-4">{prompt.content}</p>
                  
                  <div className="flex flex-wrap justify-between items-center">
                    <div className="flex flex-wrap gap-2 mb-2 md:mb-0">
                      {prompt.tags.map(tag => (
                        <span 
                          key={tag} 
                          className="bg-black bg-opacity-30 text-white px-3 py-1 rounded-full text-sm"
                        >
                          #{tag}
                        </span>
                      ))}
                    </div>
                    <span className="text-gray-200 text-sm">{formatDate(prompt.created_at)}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-400 text-center py-8">No prompts to display</p>
        )}
      </div>
    </div>
  );
};

export default UserProfile;
