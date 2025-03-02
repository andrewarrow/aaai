import React from 'react';
import { Link, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Layout = () => {
  const { user, logout, loading } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    const result = await logout();
    if (result.success) {
      navigate('/');
    }
  };

  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-gray-800 py-4 px-4 sm:px-6">
        <div className="max-w-6xl mx-auto flex justify-between items-center">
          <Link to="/" className="text-2xl font-bold text-purple-500">VibeCoders</Link>
          
          <nav className="flex space-x-4">
            <Link to="/" className="text-gray-300 hover:text-white px-3 py-2 rounded-md">Home</Link>
            {!loading && (
              user ? (
                <>
                  <Link to="/profile" className="text-gray-300 hover:text-white px-3 py-2 rounded-md">Profile</Link>
                  {user.is_admin && (
                    <Link to="/admin" className="text-yellow-400 hover:text-yellow-300 px-3 py-2 rounded-md">Admin</Link>
                  )}
                  <button 
                    onClick={handleLogout} 
                    className="text-gray-300 hover:text-white px-3 py-2 rounded-md"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link to="/login" className="text-gray-300 hover:text-white px-3 py-2 rounded-md">Login</Link>
                  <Link to="/register" className="text-gray-300 hover:text-white px-3 py-2 rounded-md">Register</Link>
                </>
              )
            )}
          </nav>
        </div>
      </header>

      <main className="flex-grow py-8 px-4 sm:px-6">
        <div className="max-w-6xl mx-auto">
          <Outlet />
        </div>
      </main>

      <footer className="bg-gray-800 py-6 px-4 sm:px-6">
        <div className="max-w-6xl mx-auto">
          <p className="text-center text-gray-400">
            &copy; {new Date().getFullYear()} VibeCoders. All rights reserved.
          </p>
          <div className="flex justify-center space-x-12 mt-6">
            <a href="https://aider.chat/" target="_blank" rel="noopener noreferrer" className="flex flex-col items-center text-gray-400 hover:text-white">
              <img src="/aider.png" alt="Aider" className="w-8 h-8 mb-2" />
              <span>Aider</span>
            </a>
            <a href="https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview" target="_blank" rel="noopener noreferrer" className="flex flex-col items-center text-gray-400 hover:text-white">
              <img src="/claude.png" alt="Claude Code" className="w-8 h-8 mb-2" />
              <span>Claude Code</span>
            </a>
            <a href="https://www.cursor.com/" target="_blank" rel="noopener noreferrer" className="flex flex-col items-center text-gray-400 hover:text-white">
              <img src="/cursor.svg" alt="Cursor" className="w-8 h-8 mb-2" />
              <span>Cursor</span>
            </a>
            <a href="https://cline.bot/" target="_blank" rel="noopener noreferrer" className="flex flex-col items-center text-gray-400 hover:text-white">
              <img src="/cline.png" alt="Cline" className="w-8 h-8 mb-2" />
              <span>Cline</span>
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout;