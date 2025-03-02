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

  return (
    <div className="max-w-2xl mx-auto">
      <div className="flex flex-col md:flex-row items-center md:items-start mb-8 gap-8">
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
    </div>
  );
};

export default UserProfile;