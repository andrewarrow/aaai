import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const UserProfile = () => {
  const { username } = useParams();
  const [user, setUser] = useState(null);
  const [prompts, setPrompts] = useState([]);
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