import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const UserProfile = () => {
  const { username } = useParams();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await fetch(`/api/users/${username}`);
        
        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('User not found');
          }
          throw new Error('Failed to fetch user data');
        }
        
        const userData = await response.json();
        setUser(userData);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
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

  // Sample prompts data (normally would come from API)
  const samplePrompts = [
    {
      id: 1,
      title: "Optimizing React Performance",
      content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
      date: "3 days ago",
      tags: ["react", "performance", "optimization"]
    },
    {
      id: 2,
      title: "Building a Neural Network from Scratch",
      content: "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
      date: "1 week ago",
      tags: ["python", "machine-learning", "neural-networks"]
    },
    {
      id: 3,
      title: "Advanced TypeScript Patterns",
      content: "Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet.",
      date: "2 weeks ago",
      tags: ["typescript", "patterns", "advanced"]
    }
  ];

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

      {/* My Latest Prompts Section */}
      <div className="mb-10">
        <h2 className="text-2xl font-bold text-purple-500 mb-6">My Latest Prompts</h2>
        
        <div className="space-y-6">
          {samplePrompts.map(prompt => (
            <div 
              key={prompt.id} 
              className={`rounded-xl overflow-hidden shadow-lg transform transition-transform hover:scale-102 cursor-pointer ${getPromptBackground(prompt.id)}`}
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
                  <span className="text-gray-200 text-sm">{prompt.date}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default UserProfile;