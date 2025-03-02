import React from 'react';
import { useNavigate } from 'react-router-dom';

const UserCard = ({ user }) => {
  const navigate = useNavigate();
  
  const handleCardClick = (e) => {
    // Don't navigate if clicking on a link
    if (e.target.tagName.toLowerCase() === 'a' || 
        e.target.parentElement.tagName.toLowerCase() === 'a') {
      return;
    }
    navigate(`/users/${user.username}`);
  };

  return (
    <div 
      className="bg-gray-800 rounded-lg overflow-hidden shadow-lg transform transition-transform hover:scale-105 cursor-pointer"
      onClick={handleCardClick}
    >
      <div className="p-6 flex flex-col items-center">
        <img 
          src={user.photo_url} 
          alt={user.username} 
          className="w-48 h-48 object-cover rounded-full mb-4"
        />
        <h3 className="text-xl font-bold text-white mb-2">{user.username}</h3>
        
        {user.bio && (
          <p className="text-gray-300 text-center mb-4">{user.bio}</p>
        )}
        
        <div className="flex space-x-4 mt-2">
          {user.github_url && (
            <a 
              href={user.github_url} 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white"
            >
              GitHub
            </a>
          )}
          
          {user.linked_in_url && (
            <a 
              href={user.linked_in_url} 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white"
            >
              LinkedIn
            </a>
          )}
        </div>
      </div>
    </div>
  );
};

export default UserCard;