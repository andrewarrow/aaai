import React, { useState, useEffect } from 'react';
import UserCard from '../components/UserCard';

const Home = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch('/api/homepage-users');
        if (response.ok) {
          const data = await response.json();
          setUsers(data);
        }
      } catch (error) {
        console.error('Error fetching users:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, []);

  return (
    <div>
      <section className="mb-12 text-center">
        <h1 className="text-4xl font-bold text-purple-500 mb-4">Welcome to VibeCoders</h1>
        <p className="text-xl text-gray-300 max-w-3xl mx-auto">
          Find software engineers that are very good at vibecoding
        </p>
      </section>

      <section className="mb-12">
        <div className="bg-gray-800 rounded-lg p-8 shadow-lg">
          <blockquote className="text-xl text-gray-300 italic mb-4">
            "Vibe coding is the art of leveraging AI tools to their fullest potential in software development, creating a seamless fusion between human creativity and machine intelligence. It represents a paradigm shift where developers orchestrate AI systems rather than writing every line manually."
          </blockquote>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-bold text-purple-500 mb-6">Top VibeCoders</h2>
        
        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-500"></div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {users.map(user => (
              <UserCard key={user.id} user={user} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
};

export default Home;