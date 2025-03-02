import React, { createContext, useState, useEffect, useContext } from 'react';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCurrentUser = async () => {
      try {
        const response = await fetch('/api/user');
        if (response.ok) {
          const userData = await response.json();
          setUser(userData);
        } else {
          setUser(null);
        }
      } catch (error) {
        console.error('Failed to fetch user:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchCurrentUser();
  }, []);

  const login = async (username, password) => {
    setError(null);
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Login failed');
      }

      // Fetch the complete user data
      const userResponse = await fetch('/api/user');
      if (userResponse.ok) {
        const userData = await userResponse.json();
        setUser(userData);
        return { success: true };
      } else {
        throw new Error('Failed to get user data after login');
      }
    } catch (error) {
      setError(error.message);
      return { success: false, error: error.message };
    }
  };

  const register = async (userData) => {
    setError(null);
    try {
      const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Registration failed');
      }

      return { success: true };
    } catch (error) {
      setError(error.message);
      return { success: false, error: error.message };
    }
  };

  const logout = async () => {
    try {
      const response = await fetch('/api/logout', {
        method: 'DELETE',
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Logout failed');
      }

      setUser(null);
      return { success: true };
    } catch (error) {
      setError(error.message);
      return { success: false, error: error.message };
    }
  };

  const updateProfile = async (userData) => {
    setError(null);
    try {
      const response = await fetch('/api/user', {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Update failed');
      }

      // Refresh user data
      const userResponse = await fetch('/api/user');
      if (userResponse.ok) {
        const updatedUserData = await userResponse.json();
        setUser(updatedUserData);
        return { success: true };
      } else {
        throw new Error('Failed to get updated user data');
      }
    } catch (error) {
      setError(error.message);
      return { success: false, error: error.message };
    }
  };

  // Prompt management functions
  const getUserPrompts = async () => {
    try {
      const response = await fetch('/api/prompts');
      
      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to fetch prompts');
      }
      
      const prompts = await response.json();
      return { success: true, prompts };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const createPrompt = async (promptData) => {
    try {
      const response = await fetch('/api/prompts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(promptData),
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to create prompt');
      }
      
      return { success: true, prompt: data };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const updatePrompt = async (promptId, promptData) => {
    try {
      const response = await fetch(`/api/prompts/${promptId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(promptData),
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to update prompt');
      }
      
      return { success: true, prompt: data };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const deletePrompt = async (promptId) => {
    try {
      const response = await fetch(`/api/prompts/${promptId}`, {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to delete prompt');
      }
      
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };

  // Project management functions
  const getUserProjects = async () => {
    try {
      const response = await fetch('/api/projects');
      
      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to fetch projects');
      }
      
      const projects = await response.json();
      return { success: true, projects };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const createProject = async (projectData) => {
    try {
      const response = await fetch('/api/projects', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(projectData),
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to create project');
      }
      
      return { success: true, project: data };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const updateProject = async (projectId, projectData) => {
    try {
      const response = await fetch(`/api/projects/${projectId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(projectData),
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Failed to update project');
      }
      
      return { success: true, project: data };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };
  
  const deleteProject = async (projectId) => {
    try {
      const response = await fetch(`/api/projects/${projectId}`, {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to delete project');
      }
      
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };

  const value = {
    user,
    loading,
    error,
    login,
    register,
    logout,
    updateProfile,
    getUserPrompts,
    createPrompt,
    updatePrompt,
    deletePrompt,
    getUserProjects,
    createProject,
    updateProject,
    deleteProject,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};